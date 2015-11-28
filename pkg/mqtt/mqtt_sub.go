package mqtt

import (
	"net"
)

type MsgHandler interface {
	Handle(msg BytesPayload)
}

type MqttSub struct {
	conn    net.Conn
	handler MsgHandler
}

func NewMqttSub(host string, h MsgHandler) (*MqttSub, error) {
	c, err := net.Dial("tcp", host)
	if err != nil {
		return nil, err
	}

	return &MqttSub{conn: c, handler: h}, nil
}
