package mqtt

import (
	"encoding/hex"
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

type Subscription struct {
	Identify   string
	DeviceId   uint64
	Host       string
	Conn       net.Conn
	SendChan   chan Message
	Cs         *Context
	SeqState   map[uint16]ResponseType
	SeqType    map[uint16]string
	KeepAlive  uint16
	LastHbTime int64
	Token      []byte
}

func NewSubscription(host string, conn net.Conn, cs *Context) {
	sendchan := make(chan Message, SendChanLen)
	sub := Subscription{
		Host:      host,
		Conn:      conn,
		SendChan:  sendchan,
		Cs:        cs,
		SeqState:  make(map[uint16]ResponseType),
		SeqType:   make(map[uint16]string),
		KeepAlive: defaultKeepAlive,
	}

	go sub.SendMsgToSub()
	go sub.RcvMsgFromClient()
}

func (sub *Subscription) Submit(msg Message) {
	if sub.Conn != nil {
		sub.SendChan <- msg
	}
}

func (sub *Subscription) ValidateToken(token []byte) error {

	err := sub.Cs.Provider.ValidateDeviceToken(sub.DeviceId, token)
	if err != nil {
		return err
	}

	sub.Token = token

	return nil
}

func (sub *Subscription) Close() {
	sub.Cs.DelSub(sub.Identify+"/c", sub.Identify)
	sub.Cs.DelConn(sub.Identify)
	if sub.Conn != nil {
		sub.Conn.Close()
		sub.Conn = nil
	}
	if sub.SendChan != nil {
		close(sub.SendChan)
		sub.SendChan = nil
	}
	sub.Cs.Provider.OnDeviceOffline(sub.DeviceId)
}

func (sub *Subscription) ReSend(msg Message, seq uint16, state uint8) {
	rspType, exist := sub.SeqState[seq]
	if !exist {
		sub.SeqState[seq] = ResponseType{
			SendTime:    0,
			PublishType: state,
			DataType:    sub.SeqType[seq],
		}
	} else {
		if rspType.SendTime == 3 {
			server.Log.Warn("send msg 3 time, abort send================\n")
			return
		}
		rspType.SendTime++
		sub.SeqState[seq] = rspType
	}

	time.Sleep(1 * time.Second)

	var needResend bool
	server.Log.Infof("now resend to %s====\n%v=======", sub.Host, msg)
	switch msg := msg.(type) {
	case *Publish:
		server.Log.Infof("try resend publish to %s", sub.Host)
		msgSeq := msg.MessageId
		oldTspType, exist := sub.SeqState[msgSeq]
		if exist && oldTspType.PublishType == uint8(MsgPublish) {
			needResend = true
		} else {
			needResend = false
		}
	case *PubRel:
		server.Log.Infof("try resend pubrel to %s", sub.Host)
		msgSeq := msg.MessageId
		oldTspType, exist := sub.SeqState[msgSeq]
		if exist && oldTspType.PublishType == uint8(MsgPubRel) {
			needResend = true
		} else {
			needResend = false
		}
	default:
		server.Log.Error("error msg type to resend")
	}

	if needResend {
		sub.SendChan <- msg
	} else {
		server.Log.Infof("no need to resend")
	}
}

func (sub *Subscription) RcvMsgFromClient() {
	host := sub.Host
	conn := sub.Conn
	server.Log.Infof("recieve new connection from %s", host)
	for {
		msg, err := DecodeOneMessage(conn)
		if err != nil {
			if err == io.EOF {
				server.Log.Debug("the end of io")
				sub.Close()
				return
			}

			if strings.HasSuffix(err.Error(), "use of closed network connection") {
				server.Log.Debug("use of closed network connection")
				sub.Close()
				return
			}

			server.Log.Error("read error: %s", err)
			sub.Close()
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

			sub.Identify = msg.ClientId

			token, _ := hex.DecodeString(msg.Password)
			err := sub.ValidateToken(token)
			if err != nil {
				server.Log.Warn("validate token error : %v", err)
				ret = RetCodeNotAuthorized
			}

			if ret != RetCodeAccepted {
				server.Log.Warn("invalid CON: %d", ret)
				goto CLOSE
			}

			deviceid, err := ClientIdToDeviceId(sub.Identify)
			if err != nil {
				server.Log.Warn("invalid Identify: %d", ret)
				goto CLOSE
			}
			sub.DeviceId = deviceid
			sub.Cs.AddConn(msg.ClientId, sub)
			connack := &ConnAck{
				ReturnCode: ret,
			}

			sub.Submit(connack)
			sub.KeepAlive = msg.KeepAliveTimer
			server.Log.Infof("%s, connected to server now", host)
			sub.Cs.Provider.OnDeviceOnline(sub.DeviceId, sub.KeepAlive)

		case *Publish:
			server.Log.Infof("%s, publish topic: %s", host, msg.TopicName)

			sub.Cs.PublishTopic2Server(msg.TopicName, sub.Identify, msg)
			if msg.QosLevel.IsAtLeastOnce() {
				server.Log.Infof("publish ack send now")
				publishack := &PubAck{MessageId: msg.MessageId}
				sub.Submit(publishack)
				sub.SeqType[msg.MessageId] = msg.TopicName
			} else if msg.QosLevel.IsExactlyOnce() {
				server.Log.Infof("publish Rec send now")
				publishRec := &PubRec{MessageId: msg.MessageId}
				sub.Submit(publishRec)
				sub.SeqType[msg.MessageId] = msg.TopicName
			}

		case *PubAck:
			server.Log.Infof("%s, comes publish ack", host)
			delete(sub.SeqState, msg.MessageId)
			delete(sub.SeqType, msg.MessageId)

		case *PubRec:
			server.Log.Infof("%s, comes publish rec", host)
			publishRel := &PubRel{MessageId: msg.MessageId}
			sub.Submit(publishRel)

		case *PubRel:
			server.Log.Infof("%s, comes publish rel", host)
			publishCom := &PubComp{MessageId: msg.MessageId}
			sub.Submit(publishCom)

		case *PubComp:
			server.Log.Infof("%s, comes publish comp", host)
			delete(sub.SeqState, msg.MessageId)
			delete(sub.SeqType, msg.MessageId)

		case *PingReq:
			server.Log.Infof("%s, ping req comes", host)
			sub.LastHbTime = time.Now().Unix()
			pingrsp := &PingResp{}
			sub.Submit(pingrsp)
			sub.Cs.Provider.OnDeviceHeartBeat(sub.DeviceId)

		case *Subscribe:
			server.Log.Infof("%s, subscribe topic: %v", host, msg.Topics)
			suback := &SubAck{
				MessageId: msg.MessageId,
				TopicsQos: make([]TagQosLevel, len(msg.Topics)),
			}
			for i, topic := range msg.Topics {
				if isWild(topic.Topic) {
					sub.Cs.AddWild(topic.Topic, sub.Identify)
				} else {
					sub.Cs.AddSub(topic.Topic, sub.Identify)
					suback.TopicsQos[i] = QosAtMostOnce
				}
			}
			sub.Submit(suback)

		case *Unsubscribe:
			server.Log.Infof("%s, unsubscribe topic: %v", host, msg.Topics)
			for _, topic := range msg.Topics {
				if isWild(topic) {
					sub.Cs.DelWild(topic, sub.Identify)
				} else {
					sub.Cs.DelSub(topic, sub.Identify)
				}
			}
			unsubAck := &UnsubAck{MessageId: msg.MessageId}
			sub.Submit(unsubAck)

		case *Disconnect:
			server.Log.Infof("%s, disconnect now, exit...", host)
			sub.Close()
			return

		default:
			server.Log.Error("unknown msg type %T", msg)
			sub.Close()
			return
		}
	}

CLOSE:
	sub.Close()

}

func (sub *Subscription) SendMsgToSub() {
	for {
		msg, ok := <-sub.SendChan
		if !ok {
			server.Log.Error("%s is end now", sub.Host)
			return
		}

		switch msg := msg.(type) {
		case *Publish:
			server.Log.Infof("publish msg, check for resend")
			if msg.QosLevel.IsAtLeastOnce() || msg.QosLevel.IsExactlyOnce() {
				go sub.ReSend(msg, msg.MessageId, uint8(MsgPublish))
			}
		case *PubRel:
			server.Log.Infof("pubrel msg, check for resend")
			go sub.ReSend(msg, msg.MessageId, uint8(MsgPubRel))
		default:
			server.Log.Infof("normal msg")
		}

		server.Log.Debug("send msg to %s=======\n%v\n=========", sub.Host, msg)
		err := msg.Encode(sub.Conn)
		if err != nil {
			server.Log.Error("send msg err: %s=====\n%v\n=====", err, msg)
			continue
		}
	}
}
