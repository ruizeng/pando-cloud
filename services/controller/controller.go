package main

import (
	"github.com/PandoCloud/pando-cloud/pkg/mongo"
	"github.com/PandoCloud/pando-cloud/pkg/rpcs"
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
	return nil
}

func (c *Controller) GetStatus(args rpcs.ArgsGetStatus, reply *rpcs.ReplyGetStatus) error {
	return nil
}

func (c *Controller) OnEvent(args rpcs.ArgsOnEvent, reply *rpcs.ReplyOnEvent) error {
	return nil
}

func (c *Controller) SendCommand(args rpcs.ArgsSendCommand, reply *rpcs.ReplySendCommand) error {
	return nil
}
