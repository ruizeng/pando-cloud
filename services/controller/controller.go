package main

import (
	"github.com/PandoCloud/pando-cloud/pkg/mongo"
	"github.com/PandoCloud/pando-cloud/pkg/rpcs"
	"github.com/PandoCloud/pando-cloud/pkg/server"
)

const (
	mongoSetName = "pandocloud"
)

type Controller struct {
	commandRecorder *mongo.Recorder
	eventRecorder   *mongo.Recorder
	dataRecorder    *mongo.Recorder
}

func NewController(mongohost string) (*Controller, error) {
	cmdr, err := mongo.NewRecorder(mongohost, mongoSetName, "commands")
	if err != nil {
		return nil, err
	}

	ever, err := mongo.NewRecorder(mongohost, mongoSetName, "events")
	if err != nil {
		return nil, err
	}

	datar, err := mongo.NewRecorder(mongohost, mongoSetName, "datas")
	if err != nil {
		return nil, err
	}

	return &Controller{
		commandRecorder: cmdr,
		eventRecorder:   ever,
		dataRecorder:    datar,
	}, nil
}

func (c *Controller) PutData(args rpcs.ArgsPutData, reply *rpcs.ReplyPutData) error {
	return c.dataRecorder.Insert(args)
}

func (c *Controller) SetStatus(args rpcs.ArgsSetStatus, reply *rpcs.ReplySetStatus) error {
	rpchost, err := getAccessRPCHost(args.DeviceId)
	if err != nil {
		return err
	}

	return server.RPCCallByHost(rpchost, "Access.SetStatus", args, reply)
}

func (c *Controller) GetStatus(args rpcs.ArgsGetStatus, reply *rpcs.ReplyGetStatus) error {
	rpchost, err := getAccessRPCHost(args.Id)
	if err != nil {
		return err
	}

	return server.RPCCallByHost(rpchost, "Access.GetStatus", args, reply)
}

func (c *Controller) OnEvent(args rpcs.ArgsOnEvent, reply *rpcs.ReplyOnEvent) error {
	return c.eventRecorder.Insert(args)
}

func (c *Controller) SendCommand(args rpcs.ArgsSendCommand, reply *rpcs.ReplySendCommand) error {
	rpchost, err := getAccessRPCHost(args.DeviceId)
	if err != nil {
		return err
	}

	return server.RPCCallByHost(rpchost, "Access.SendCommand", args, reply)
}

func getAccessRPCHost(deviceid uint64) (string, error) {
	args := rpcs.ArgsGetDeviceOnlineStatus{
		Id: deviceid,
	}
	reply := &rpcs.ReplyGetDeviceOnlineStatus{}
	err := server.RPCCallByName("devicemanager", "DeviceManager.GetDeviceOnlineStatus", args, reply)
	if err != nil {
		return "", err
	}

	return reply.AccessRPCHost, nil
}
