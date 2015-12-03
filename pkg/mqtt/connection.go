package mqtt

import (
	"encoding/hex"
	"github.com/PandoCloud/pando-cloud/pkg/rpcs"
	"github.com/PandoCloud/pando-cloud/pkg/server"
	"io"
	"net"
	"strings"
	"time"
)

const (
	SendChanLen      = 16
	defaultKeepAlive = 30
)

type ResponseType struct {
	SendTime    uint8
	PublishType uint8
	DataType    string
}

type Connection struct {
	DeviceId   uint64
	Conn       net.Conn
	SendChan   chan Message
	Mgr        *Manager
	MessageId  uint16
	KeepAlive  uint16
	LastHbTime int64
	Token      []byte
}

func NewConnection(conn net.Conn, mgr *Manager) *Connection {
	sendchan := make(chan Message, SendChanLen)
	c := &Connection{
		Conn:      conn,
		SendChan:  sendchan,
		Mgr:       mgr,
		KeepAlive: defaultKeepAlive,
	}

	go c.SendMsgToClient()
	go c.RcvMsgFromClient()

	return c
}

func (c *Connection) Submit(msg Message) {
	if c.Conn != nil {
		c.SendChan <- msg
	}
}

func (c *Connection) ValidateToken(token []byte) error {

	err := c.Mgr.Provider.ValidateDeviceToken(c.DeviceId, token)
	if err != nil {
		return err
	}

	c.Token = token

	return nil
}

func (c *Connection) Close() {
	c.Mgr.DelConn(c.DeviceId)
	if c.Conn != nil {
		c.Conn.Close()
		c.Conn = nil
	}
	if c.SendChan != nil {
		close(c.SendChan)
		c.SendChan = nil
	}
	c.Mgr.Provider.OnDeviceOffline(c.DeviceId)
}

func (c *Connection) RcvMsgFromClient() {
	conn := c.Conn
	host := conn.RemoteAddr().String()
	server.Log.Infof("recieve new connection from %s", host)
	for {
		msg, err := DecodeOneMessage(conn)
		if err != nil {
			if err == io.EOF {
				server.Log.Debug("the end of io")
				c.Close()
				return
			}

			if strings.HasSuffix(err.Error(), "use of closed network connection") {
				server.Log.Debug("use of closed network connection")
				c.Close()
				return
			}

			server.Log.Error("read error: %s", err)
			c.Close()
			return
		}

		server.Log.Infof("%s, come msg===\n%v\n=====", host, msg)
		switch msg := msg.(type) {
		case *Connect:
			ret := RetCodeAccepted
			if msg.ProtocolVersion == 3 && msg.ProtocolName != "MQIsdp" {
				ret = RetCodeUnacceptableProtocolVersion
			} else if msg.ProtocolVersion == 4 && msg.ProtocolName != "MQTT" {
				ret = RetCodeUnacceptableProtocolVersion
			} else if msg.ProtocolVersion > 4 {
				ret = RetCodeUnacceptableProtocolVersion
			}

			if len(msg.ClientId) < 1 || len(msg.ClientId) > 23 {
				ret = RetCodeIdentifierRejected
			}

			deviceid, err := ClientIdToDeviceId(msg.ClientId)
			if err != nil {
				server.Log.Warn("invalid Identify: %d", ret)
				goto CLOSE
			}
			c.DeviceId = deviceid

			token, _ := hex.DecodeString(msg.Password)
			err = c.ValidateToken(token)
			if err != nil {
				server.Log.Warn("validate token error : %v", err)
				ret = RetCodeNotAuthorized
			}

			if ret != RetCodeAccepted {
				server.Log.Warn("invalid CON: %d", ret)
				goto CLOSE
			}

			c.Mgr.AddConn(c.DeviceId, c)
			connack := &ConnAck{
				ReturnCode: ret,
			}

			c.Submit(connack)
			c.KeepAlive = msg.KeepAliveTimer
			server.Log.Infof("%s, connected to server now", host)
			args := rpcs.ArgsGetOnline{
				Id:                c.DeviceId,
				ClientIP:          host,
				AccessRPCHost:     server.GetRPCHost(),
				HeartbeatInterval: uint32(c.KeepAlive),
			}
			c.Mgr.Provider.OnDeviceOnline(args)

		case *Publish:
			server.Log.Infof("%s, publish topic: %s", host, msg.TopicName)

			c.Mgr.PublishMessage2Server(c.DeviceId, msg)
			if msg.QosLevel.IsAtLeastOnce() {
				server.Log.Infof("publish ack send now")
				publishack := &PubAck{MessageId: msg.MessageId}
				c.Submit(publishack)
			} else if msg.QosLevel.IsExactlyOnce() {
				server.Log.Infof("publish Rec send now")
				publishRec := &PubRec{MessageId: msg.MessageId}
				c.Submit(publishRec)
			}

		case *PubAck:
			server.Log.Infof("%s, comes publish ack", host)
			// TODO  - notify sender

		case *PubRec:
			server.Log.Infof("%s, comes publish rec", host)
			publishRel := &PubRel{MessageId: msg.MessageId}
			c.Submit(publishRel)

		case *PubRel:
			server.Log.Infof("%s, comes publish rel", host)
			publishCom := &PubComp{MessageId: msg.MessageId}
			c.Submit(publishCom)

		case *PubComp:
			server.Log.Infof("%s, comes publish comp", host)
			// TODO  - notify sender

		case *PingReq:
			server.Log.Infof("%s, ping req comes", host)
			c.LastHbTime = time.Now().Unix()
			pingrsp := &PingResp{}
			c.Submit(pingrsp)
			c.Mgr.Provider.OnDeviceHeartBeat(c.DeviceId)

		case *Subscribe:
			server.Log.Infof("%s, subscribe topic: %v", host, msg.Topics)

		case *Unsubscribe:
			server.Log.Infof("%s, unsubscribe topic: %v", host, msg.Topics)

		case *Disconnect:
			server.Log.Infof("%s, disconnect now, exit...", host)
			c.Close()
			return

		default:
			server.Log.Error("unknown msg type %T", msg)
			c.Close()
			return
		}
	}

CLOSE:
	c.Close()

}

func (c *Connection) SendMsgToClient() {
	host := c.Conn.RemoteAddr()
	for {
		msg, ok := <-c.SendChan
		if !ok {
			server.Log.Error("%s is end now", host)
			return
		}

		switch msg := msg.(type) {
		case *Publish:
			server.Log.Infof("publish msg, check for resend")
			if msg.QosLevel.IsAtLeastOnce() || msg.QosLevel.IsExactlyOnce() {
				msg.MessageId = c.MessageId
				c.MessageId++
			}
		case *PubRel:
			server.Log.Infof("pubrel msg, check for resend")
			msg.MessageId = c.MessageId
			c.MessageId++
		default:
			server.Log.Infof("normal msg")
		}

		server.Log.Debug("send msg to %s=======\n%v\n=========", host, msg)
		err := msg.Encode(c.Conn)
		if err != nil {
			server.Log.Error("send msg err: %s=====\n%v\n=====", err, msg)
			continue
		}
	}
}
