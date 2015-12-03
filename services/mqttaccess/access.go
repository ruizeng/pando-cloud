package main

import (
	"github.com/PandoCloud/pando-cloud/pkg/mqtt"
	"github.com/PandoCloud/pando-cloud/pkg/protocol"
	"github.com/PandoCloud/pando-cloud/pkg/rpcs"
	"github.com/PandoCloud/pando-cloud/pkg/server"
	"time"
)

type Access struct {
	MqttBroker *mqtt.Broker
}

func NewAccess() (*Access, error) {
	p := NewMQTTProvider()
	return &Access{
		mqtt.NewBroker(p),
	}, nil
}

func (a *Access) SetStatus(args rpcs.ArgsSetStatus, reply *rpcs.ReplySetStatus) error {
	server.Log.Infof("Access Set Status: %v", args)
	data := &protocol.Data{}
	data.Head.Timestamp = uint64(time.Now().Unix())
	token, err := a.MqttBroker.GetToken(args.DeviceId)
	if err != nil {
		return err
	}
	copy(data.Head.Token[:], token[:16])
	data.SubData = args.Status
	msg, err := data.Marshal()
	if err != nil {
		return err
	}
	return a.MqttBroker.SendMessageToDevice(args.DeviceId, "d", msg)
}

func (a *Access) GetStatus(args rpcs.ArgsGetStatus, reply *rpcs.ReplyGetStatus) error {
	server.Log.Infof("Access Get Status: %v", args)
	return nil
}

func (a *Access) SendCommand(args rpcs.ArgsSendCommand, reply *rpcs.ReplySendCommand) error {
	server.Log.Infof("Access Send Command: %v", args)
	cmd := &protocol.Command{}
	cmd.Head.Timestamp = uint64(time.Now().Unix())
	token, err := a.MqttBroker.GetToken(args.DeviceId)
	if err != nil {
		return err
	}
	copy(cmd.Head.Token[:], token[:16])
	cmd.Head.No = args.No
	cmd.Head.Priority = args.Priority
	cmd.Head.SubDeviceid = args.SubDevice
	cmd.Params = args.Params
	msg, err := cmd.Marshal()
	if err != nil {
		return err
	}
	return a.MqttBroker.SendMessageToDevice(args.DeviceId, "c", msg)
}
