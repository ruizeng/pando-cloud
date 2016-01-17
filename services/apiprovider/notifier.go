package main

import (
	"encoding/json"
	"github.com/PandoCloud/pando-cloud/pkg/models"
	"github.com/PandoCloud/pando-cloud/pkg/productconfig"
	"github.com/PandoCloud/pando-cloud/pkg/protocol"
	"github.com/PandoCloud/pando-cloud/pkg/queue"
	"github.com/PandoCloud/pando-cloud/pkg/rpcs"
	"github.com/PandoCloud/pando-cloud/pkg/server"
	"github.com/PandoCloud/pando-cloud/pkg/utils"
	"time"
)

const (
	topicEvents = "events"
	topicStatus = "status"
)

// report structure
type ReportPack struct {
	Tag        string                   `json:"tag"`
	Identifier string                   `json:"identifier"`
	TimeStamp  int64                    `json:"timestamp"`
	Data       map[string][]interface{} `json:"data"`
}

var notifier *Notifier

type Notifier struct {
	eventsQueue *queue.Queue
	statusQueue *queue.Queue
	apps        []*models.Application
}

func NewNotifier(rabbithost string) (*Notifier, error) {
	eq, err := queue.New(rabbithost, topicEvents)
	if err != nil {
		return nil, err
	}

	sq, err := queue.New(rabbithost, topicStatus)
	if err != nil {
		return nil, err
	}

	return &Notifier{
		eventsQueue: eq,
		statusQueue: sq,
	}, nil
}

// TODO
func (n *Notifier) reportStatus(event rpcs.ArgsOnStatus) error {
	return nil
}

// TODO
func (n *Notifier) processStatus() error {
	return nil
}

func (n *Notifier) updateApplications() error {
	for {

		err := server.RPCCallByName("registry", "Registry.GetApplications", 0, &n.apps)
		if err != nil {
			server.Log.Errorf("get applications error : %v", err)
		}

		time.Sleep(time.Minute)
	}
}

func (n *Notifier) reportEvent(event rpcs.ArgsOnEvent) error {
	server.Log.Debugf("reporting event %v", event)

	device := &models.Device{}
	err := server.RPCCallByName("registry", "Registry.FindDeviceById", int64(event.DeviceId), device)
	if err != nil {
		server.Log.Errorf("find device error : %v", err)
		return err
	}

	product := &models.Product{}
	err = server.RPCCallByName("registry", "Registry.FindProduct", device.ProductID, product)
	if err != nil {
		server.Log.Errorf("find product error : %v", err)
		return err
	}

	c, err := productconfig.New(product.ProductConfig)
	if err != nil {
		server.Log.Errorf("product config error : %v", err)
		return err
	}

	ev := &protocol.Event{}
	ev.Head.No = event.No
	ev.Head.SubDeviceid = event.SubDevice
	ev.Params = event.Params

	m, err := c.EventToMap(ev)
	if err != nil {
		server.Log.Errorf("gen event json error : %v", err)
		return err
	}

	res := ReportPack{
		Tag:        "event",
		Identifier: device.DeviceIdentifier,
		Data:       m,
	}

	jsonRes, err := json.Marshal(res)
	if err != nil {
		server.Log.Errorf("json marshal error : %v", err)
		return err
	}

	reqHead := map[string]string{}
	reqHead["Content-Type"] = "application/json"

	for _, app := range n.apps {
		if nil == checkAppDomain(app.AppDomain, device.DeviceIdentifier) {
			reqHead["App-Token"] = app.AppToken
			_, err := utils.SendHttpRequest(app.ReportUrl, string(jsonRes), "POST", reqHead)
			if err != nil {
				server.Log.Errorf("http post json error : %v", err)
			}
			server.Log.Debugf("http post json succ : %v", string(jsonRes))
		}
	}

	return nil
}

func (n *Notifier) processEvents() error {
	for {
		event := rpcs.ArgsOnEvent{}
		err := n.eventsQueue.Receive(&event)
		if err != nil {
			server.Log.Errorf("error when receiving from queue : %v", err)
			return err
		}
		go n.reportEvent(event)
	}

	return nil
}

func (n *Notifier) Run() error {
	go n.updateApplications()
	go n.processEvents()
	go n.processStatus()

	return nil
}

func RunNotifier() error {
	if notifier == nil {
		notifier, err := NewNotifier(*confRabbitHost)
		if err != nil {
			server.Log.Error(err)
		}
		err = notifier.Run()
		if err != nil {
			server.Log.Error(err)
		}
	}
	return nil
}
