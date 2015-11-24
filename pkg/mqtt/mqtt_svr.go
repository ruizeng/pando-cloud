package mqtt

import (
	"net"
)

type MqttSvr struct {
	listener   net.Listener
}

func NewMqttSvr(host string) (*MqttSvr, error) {
	l, err := net.Listen("tcp", host)	
	if err != nil {
		return nil, err
	}

	return &MqttSvr{listener:l}, nil
}

func (server *MqttSvr) Run() {
	<- make(chan bool)
}
