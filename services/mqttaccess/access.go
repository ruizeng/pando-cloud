package main

import (
	"github.com/PandoCloud/pando-cloud/pkg/mqtt"
	"github.com/PandoCloud/pando-cloud/pkg/rpcs"
	"github.com/PandoCloud/pando-cloud/pkg/server"
)

type Access struct {
	MqttHandler *mqtt.MqttSvrHandler
}

func NewAccess() (*Access, error) {
	return &Access{
		mqtt.NewMqttSvrHandler(),
	}, nil
}

func (a *Access) SetStatus(args rpcs.ArgsSetStatus, reply *rpcs.ReplySetStatus) error {
	server.Log.Infof("Access Set Status: %v", args)
	return nil
}

func (a *Access) GetStatus(args rpcs.ArgsGetStatus, reply *rpcs.ReplyGetStatus) error {
	server.Log.Infof("Access Get Status: %v", args)
	return nil
}

func (a *Access) SendCommand(args rpcs.ArgsSendCommand, reply *rpcs.ReplySendCommand) error {
	server.Log.Infof("Access Send Command: %v", args)
	return nil
}
