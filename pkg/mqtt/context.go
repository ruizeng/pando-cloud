package mqtt

import (
	"errors"
	"net"
	"sync"
	"time"
)

type Context struct {
	Provider  Provider
	CxtMutex  sync.RWMutex
	IdToSub   map[string]*Subscription
	TopicToId map[string][]string
	Wildcards []Wild
}

func NewContext(p Provider) *Context {
	cxt := &Context{
		Provider:  p,
		IdToSub:   make(map[string]*Subscription),
		TopicToId: make(map[string][]string),
	}

	go cxt.CleanWorker()

	return cxt
}

func (cs *Context) NewSub(host string, conn net.Conn) {
	NewSubscription(host, conn, cs)
}

func (cs *Context) AddConn(identify string, sub *Subscription) {
	cs.CxtMutex.Lock()
	oldSub, exist := cs.IdToSub[identify]
	if exist {
		oldSub.Close()
	}

	cs.IdToSub[identify] = sub
	cs.CxtMutex.Unlock()
}

func (cs *Context) DelConn(identify string) {
	cs.CxtMutex.Lock()
	_, exist := cs.IdToSub[identify]

	if exist {
		delete(cs.IdToSub, identify)
	}
	cs.CxtMutex.Unlock()
}

func (cs *Context) PublishTopic2Device(deviceid uint64, msg *Publish) error {
	identify := DeviceIdToClientId(deviceid)
	cs.CxtMutex.RLock()
	sub, exist := cs.IdToSub[identify]
	cs.CxtMutex.RUnlock()
	if !exist {
		return errorf("device not exist: %v[%v]", deviceid, identify)
	}

	sub.Submit(msg)

	return nil
}

func (cs *Context) PublishTopic2Server(topic string, identify string, msg *Publish) error {
	deviceid, err := ClientIdToDeviceId(identify)
	if err != nil {
		return errorf("wrong clientid for event: %s", identify)
	}
	payload := msg.Payload.(BytesPayload)

	cs.Provider.OnDeviceMessage(deviceid, payload)
	return nil
}

func (cs *Context) AddWild(topic string, identify string) error {
	if !isWildValid(topic) {
		return errorf("wrong wildcard expressiong: %v", topic)
	}

	cs.CxtMutex.Lock()

	for _, w := range cs.Wildcards {
		if w.isExist(topic, identify) {
			return nil
		}
	}

	wild := NewWild(topic, identify)
	cs.Wildcards = append(cs.Wildcards, *wild)

	cs.CxtMutex.Unlock()
	return nil
}

func (cs *Context) DelWild(topic, identify string) {
	newWilds := []Wild{}
	flag := false
	cs.CxtMutex.Lock()
	for _, w := range cs.Wildcards {
		if !w.isExist(topic, identify) {
			newWilds = append(newWilds, w)
		} else {
			flag = true
		}
	}

	if flag {
		cs.Wildcards = newWilds
	}
	cs.CxtMutex.Unlock()
}

func (cs *Context) AddSub(topic string, identify string) error {
	if len(identify) == 0 {
		return errors.New("identify is empty")
	}

	cs.CxtMutex.Lock()
	ids, exist := cs.TopicToId[topic]
	if !exist {
		cs.TopicToId[topic] = []string{}
	}

	for i := 0; i < len(ids); i++ {
		if ids[i] == identify {
			return nil
		}
	}

	cs.TopicToId[topic] = append(cs.TopicToId[topic], identify)

	cs.CxtMutex.Unlock()

	return nil
}

func (cs *Context) DelSub(topic string, identify string) error {
	if len(identify) == 0 {
		return errors.New("identify is empty")
	}

	cs.CxtMutex.Lock()
	ids, exist := cs.TopicToId[topic]
	if !exist {
		return errorf("topic[%s] not exist", topic)
	}

	newids := []string{}
	for i := 0; i < len(ids); i++ {
		if ids[i] == identify {
			continue
		}
		newids = append(newids, ids[i])
	}

	if len(newids) == 0 {
		delete(cs.TopicToId, topic)
	} else {
		cs.TopicToId[topic] = newids
	}

	cs.CxtMutex.Unlock()

	return nil
}

func (cs *Context) CleanWorker() {
	for {
		curTime := time.Now().Unix()

		for _, sub := range cs.IdToSub {
			if sub.KeepAlive == 0 {
				continue
			}

			if uint16(curTime-sub.LastHbTime) > uint16(2*sub.KeepAlive/2) {
				sub.Close()
			}
		}

		time.Sleep(60 * time.Second)
	}
}
