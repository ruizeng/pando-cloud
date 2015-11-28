package mqtt

import (
	"fmt"
	"io"
	"net"
	"strings"
)

const (
	SendChanLen      = 16
	defaultKeepAlive = 30
)

type Subscription struct {
	Identify   string
	Host       string
	Conn       net.Conn
	SendChan   chan Message
	mqttHandle *MqttSvrHandler
}

func NewSubscription(host string, conn net.Conn, h *MqttSvrHandler) {
	sub := Subscription{
		Host:       host,
		Conn:       conn,
		SendChan:   make(chan Message, SendChanLen),
		mqttHandle: h,
	}

	go sub.SendMsgToSub()
	go sub.RcvMsgFromClient()
}

func (sub *Subscription) Submit(msg Message) {
	if sub.Conn != nil {
		sub.SendChan <- msg
	}
}

func (sub *Subscription) Close() {

}

func (sub *Subscription) RcvMsgFromClient() {
	host := sub.Host
	conn := sub.Conn
	fmt.Printf("recieve new connection from %s\n", host)
	for {
		msg, err := DecodeOneMessage(conn)
		if err != nil {
			if err == io.EOF {
				fmt.Printf("the end of io\n")
				sub.Close()
				return
			}

			if strings.HasSuffix(err.Error(), "use of closed network connection") {
				fmt.Printf("use of closed network connection\n")
				sub.Close()
				return
			}

			fmt.Printf("read error: %s\n", err)
			sub.Close()
			return
		}

		fmt.Printf("%s, come msg===\n%v\n=====\n", host, msg)
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

			if ret != RetCodeAccepted {
				fmt.Printf("invalid CON: %d\n", ret)
				sub.Close()
				return
			}

			sub.Identify = msg.ClientId
			sub.mqttHandle.AddConn(sub.Identify, sub)

			connack := &ConnAck{
				ReturnCode: ret,
			}

			sub.Submit(connack)
			fmt.Printf("%s, connected to server now\n", host)

		case *Publish:
			fmt.Printf("%s, publish topic: %s\n", host, msg.TopicName)
			sub.mqttHandle.PublishTopic(msg.TopicName, msg)

		case *PubAck:
			fmt.Printf("%s, comes publish ack\n", host)

		case *PubRec:
			fmt.Printf("%s, comes publish rec\n", host)
			publishRel := &PubRel{MessageId: msg.MessageId}
			sub.Submit(publishRel)

		case *PubRel:
			fmt.Printf("%s, comes publish rel\n", host)
			publishCom := &PubComp{MessageId: msg.MessageId}
			sub.Submit(publishCom)

		case *PubComp:
			fmt.Printf("%s, comes publish comp\n", host)

		case *PingReq:
			fmt.Printf("%s, ping req comes\n", host)
			pingrsp := &PingResp{}
			sub.Submit(pingrsp)

		case *Subscribe:
			fmt.Printf("%s, subscribe topic: %v\n", host, msg.Topics)
			suback := &SubAck{
				MessageId: msg.MessageId,
				TopicsQos: make([]TagQosLevel, len(msg.Topics)),
			}
			for i, topic := range msg.Topics {
				if isWild(topic.Topic) {
					sub.mqttHandle.AddWild(topic.Topic, sub.Identify)
				} else {
					sub.mqttHandle.AddSub(topic.Topic, sub.Identify)
					suback.TopicsQos[i] = QosAtMostOnce
				}
			}
			sub.Submit(suback)

		case *Unsubscribe:
			fmt.Printf("%s, unsubscribe topic: %v\n", host, msg.Topics)
			for _, topic := range msg.Topics {
				if isWild(topic) {
					sub.mqttHandle.DelWild(topic, sub.Identify)
				} else {
					sub.mqttHandle.DelSub(topic, sub.Identify)
				}
			}
			unsubAck := &UnsubAck{MessageId: msg.MessageId}
			sub.Submit(unsubAck)

		case *Disconnect:
			fmt.Printf("%s, disconnect now, exit...\n", host)
			sub.Close()
			return

		default:
			fmt.Printf("unknown msg type %T\n", msg)
			sub.Close()
			return
		}
	}
}

func (sub *Subscription) SendMsgToSub() {
	for {
		msg, ok := <-sub.SendChan
		if !ok {
			fmt.Printf("%s is end now\n", sub.Host)
			return
		}

		fmt.Printf("send msg to %s=======\n%v\n=========\n", sub.Host, msg)
		err := msg.Encode(sub.Conn)
		if err != nil {
			fmt.Printf("send msg err: %s=====\n%v\n=====\n", err, msg)
			continue
		}
	}
}
