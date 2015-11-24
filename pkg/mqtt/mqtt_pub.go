package mqtt

import (
	"net"
)

type MqttPub struct {
	conn net.Conn
}

func NewMqttPub(host string) (*MqttPub, error) {
	c, err := net.Dial("tcp", host)
	if err != nil {
		return nil, err
	}

	return &MqttPub{conn:c}, nil
}

func (pub *MqttPub) PublishMsg(msg BytesPayload) error {
	return nil
}
