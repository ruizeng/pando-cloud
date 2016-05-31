package main

import (
	"errors"
	"github.com/PandoCloud/pando-cloud/pkg/mqtt"
	"github.com/PandoCloud/pando-cloud/pkg/protocol"
	"github.com/PandoCloud/pando-cloud/pkg/rpcs"
	"github.com/PandoCloud/pando-cloud/pkg/server"
	"time"
)

const (
	defaultTimeoutSecond = 5

	commandGetCurrentStatus = uint16(65528)
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
	return a.MqttBroker.SendMessageToDevice(args.DeviceId, "s", msg, defaultTimeoutSecond*time.Second)
}

func (a *Access) GetStatus(args rpcs.ArgsGetStatus, reply *rpcs.ReplyGetStatus) error {
	server.Log.Infof("Access Get Status: %v", args)
	// first send a get status command
	cmdArgs := rpcs.ArgsSendCommand{
		DeviceId:  args.Id,
		SubDevice: 65535,
		No:        commandGetCurrentStatus,
		Priority:  99,
		WaitTime:  0,
	}
	cmdReply := rpcs.ReplySendCommand{}
	err := a.SendCommand(cmdArgs, &cmdReply)
	if err != nil {
		return err
	}

	// then wait for status report
	StatusChan[args.Id] = make(chan *protocol.Data)
	after := time.After(defaultTimeoutSecond * time.Second)
	server.Log.Debug("now waiting 5 seconds for status report...")
	select {
	case <-after:
		// timeout
		close(StatusChan[args.Id])
		delete(StatusChan, args.Id)
		return errors.New("get status timeout.")
	case data := <-StatusChan[args.Id]:
		// go it
		close(StatusChan[args.Id])
		delete(StatusChan, args.Id)
		reply.Status = data.SubData
		return nil
	}
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
	cmd.Head.ParamsCount = uint16(len(args.Params))
	cmd.Params = args.Params
	msg, err := cmd.Marshal()
	if err != nil {
		return err
	}
	return a.MqttBroker.SendMessageToDevice(args.DeviceId, "c", msg, time.Duration(args.WaitTime)*time.Second)
}
