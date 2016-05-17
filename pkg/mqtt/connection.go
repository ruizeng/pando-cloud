package mqtt

import (
	"encoding/hex"
	"errors"
	"github.com/PandoCloud/pando-cloud/pkg/rpcs"
	"github.com/PandoCloud/pando-cloud/pkg/server"
	"net"
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
	Mgr             *Manager
	DeviceId        uint64
	Conn            net.Conn
	SendChan        chan Message
	MessageId       uint16
	MessageWaitChan map[uint16]chan error
	KeepAlive       uint16
	LastHbTime      int64
	Token           []byte
}

func NewConnection(conn net.Conn, mgr *Manager) *Connection {
	sendchan := make(chan Message, SendChanLen)
	c := &Connection{
		Conn:            conn,
		SendChan:        sendchan,
		Mgr:             mgr,
		KeepAlive:       defaultKeepAlive,
		MessageWaitChan: make(map[uint16]chan error),
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

// Publish will publish a message , and return a chan to wait for completion.
func (c *Connection) Publish(msg Message, timeout time.Duration) error {
	server.Log.Debugf("publishing message : %v, timeout %v", msg, timeout)

	message := msg.(*Publish)
	message.MessageId = c.MessageId
	c.MessageId++
	c.Submit(message)

	ch := make(chan error)

	// we don't wait for confirm.
	if timeout == 0 {
		return nil
	}

	c.MessageWaitChan[message.MessageId] = ch
	// wait for timeout and
	go func() {
		timer := time.NewTimer(timeout)
		<-timer.C
		waitCh, exist := c.MessageWaitChan[message.MessageId]
		if exist {
			waitCh <- errors.New("timeout pushlishing message.")
			delete(c.MessageWaitChan, message.MessageId)
			close(waitCh)
		}
	}()

	err := <-ch
	return err
}

func (c *Connection) confirmPublish(messageid uint16) {
	waitCh, exist := c.MessageWaitChan[messageid]
	if exist {
		waitCh <- nil
		delete(c.MessageWaitChan, messageid)
		close(waitCh)
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
	deviceid := c.DeviceId
	server.Log.Infof("closing connection of device %v", deviceid)
	if c.Conn != nil {
		c.Conn.Close()
		c.Conn = nil
		c.Mgr.Provider.OnDeviceOffline(deviceid)
	}
	if c.SendChan != nil {
		close(c.SendChan)
		c.SendChan = nil
	}
}

func (c *Connection) RcvMsgFromClient() {
	conn := c.Conn
	host := conn.RemoteAddr().String()
	server.Log.Infof("recieve new connection from %s", host)
	for {
		msg, err := DecodeOneMessage(conn)
		if err != nil {

			server.Log.Errorf("read error: %s", err)
			c.Close()
			return
		}

		server.Log.Infof("%s, come msg===\n%v\n=====", host, msg)
		c.LastHbTime = time.Now().Unix()
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
				server.Log.Warn("invalid clientid length: %d", len(msg.ClientId))
				ret = RetCodeIdentifierRejected
				c.Close()
				return
			}

			deviceid, err := ClientIdToDeviceId(msg.ClientId)
			if err != nil {
				server.Log.Warn("invalid Identify: %d", ret)
				c.Close()
				return
			}
			c.DeviceId = deviceid

			token, err := hex.DecodeString(msg.Password)
			if err != nil {
				server.Log.Warn("token format error : %v", err)
				ret = RetCodeNotAuthorized
				c.Close()
				return
			}
			err = c.ValidateToken(token)
			if err != nil {
				server.Log.Warn("validate token error : %v", err)
				ret = RetCodeNotAuthorized
				c.Close()
				return
			}

			if ret != RetCodeAccepted {
				server.Log.Warn("invalid CON: %d", ret)
				c.Close()
				return
			}

			args := rpcs.ArgsGetOnline{
				Id:                c.DeviceId,
				ClientIP:          host,
				AccessRPCHost:     server.GetRPCHost(),
				HeartbeatInterval: uint32(c.KeepAlive),
			}

			c.Mgr.AddConn(c.DeviceId, c)
			connack := &ConnAck{
				ReturnCode: ret,
			}

			c.Submit(connack)
			c.KeepAlive = msg.KeepAliveTimer

			err = c.Mgr.Provider.OnDeviceOnline(args)
			if err != nil {
				server.Log.Warn("device online error : %v", err)
				c.Close()
				return
			}

			server.Log.Infof("device %d, connected to server now, host: %s", c.DeviceId, host)

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

			err := c.Mgr.Provider.OnDeviceHeartBeat(c.DeviceId)
			if err != nil {
				server.Log.Warnf("%s, heartbeat set error %s, close now...", host, err)
				c.Close()
				return
			}

		case *PubAck:
			server.Log.Infof("%s, comes publish ack", host)
			c.confirmPublish(msg.MessageId)
			err := c.Mgr.Provider.OnDeviceHeartBeat(c.DeviceId)
			if err != nil {
				server.Log.Warnf("%s, heartbeat set error %s, close now...", host, err)
				c.Close()
				return
			}

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
			c.confirmPublish(msg.MessageId)
			err := c.Mgr.Provider.OnDeviceHeartBeat(c.DeviceId)
			if err != nil {
				server.Log.Warnf("%s, heartbeat set error %s, close now...", host, err)
				c.Close()
				return
			}

		case *PingReq:
			server.Log.Infof("%s, ping req comes", host)
			pingrsp := &PingResp{}
			err := c.Mgr.Provider.OnDeviceHeartBeat(c.DeviceId)
			if err != nil {
				server.Log.Warnf("%s, heartbeat set error %s, close now...", host, err)
				c.Close()
				return
			}
			c.Submit(pingrsp)

		case *Subscribe:
			server.Log.Infof("%s, subscribe topic: %v", host, msg.Topics)

		case *Unsubscribe:
			server.Log.Infof("%s, unsubscribe topic: %v", host, msg.Topics)

		case *Disconnect:
			server.Log.Infof("%s, disconnect now, exit...", host)
			c.Close()
			return

		default:
			server.Log.Errorf("unknown msg type %T", msg)
			c.Close()
			return
		}
	}

}

func (c *Connection) SendMsgToClient() {
	host := c.Conn.RemoteAddr()
	for {
		msg, ok := <-c.SendChan
		if !ok {
			server.Log.Errorf("%s is end now", host)
			return
		}

		server.Log.Debugf("send msg to %s=======\n%v\n=========", host, msg)
		err := msg.Encode(c.Conn)
		if err != nil {
			server.Log.Errorf("send msg err: %s=====\n%v\n=====", err, msg)
			continue
		}
	}
}
