package main

import (
	"github.com/PandoCloud/pando-cloud/pkg/rpcs"
)

const (
	topicEvents = "events"
	topicStatus = "status"
)

var notifier *Notifier

type Notifier struct {
	eventsQueue *queue.Queue
	statusQueue *queue.Queue
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

func (n *Notifier) reportEvent(event rpcs.ArgsOnEvent) error {

}

func (n *Notifier) processEvents() error {
	for {
		event := rpcs.ArgsOnEvent{}
		err := n.eventsQueue.Receive(&event)
		if err != nil {
			return err
		}

	}
}

func (n *Notifier) processStatus() error {
}

func (n *Notifier) Run() error {
}

func RunNotifier() error {

}
