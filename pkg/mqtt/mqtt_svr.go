package mqtt

import (
	"fmt"
	"net"
	"sync"
)

type MqttSvrHandler struct {
	CxtMutex  sync.RWMutex
	IdToSub   map[string]*Subscription
	TopicToId map[string][]string
	Wildcards []Wild
}

func NewMqttSvrHandler() *MqttSvrHandler {
	return &MqttSvrHandler{}
}

func (handler *MqttSvrHandler) Handle(conn net.Conn) {
	host := conn.RemoteAddr().String()
	NewSubscription(host, conn, handler)
}

func (handler *MqttSvrHandler) AddConn(identify string, sub *Subscription) {
	handler.CxtMutex.Lock()
	oldSub, exist := handler.IdToSub[identify]
	if exist {
		oldSub.Close()
	}

	handler.IdToSub[identify] = sub
	handler.CxtMutex.Unlock()
}

func (handler *MqttSvrHandler) DelConn(identify string) {
	handler.CxtMutex.Lock()
	_, exist := handler.IdToSub[identify]

	if exist {
		delete(handler.IdToSub, identify)
	}
	handler.CxtMutex.Unlock()
}

func (handler *MqttSvrHandler) AddWild(topic string, identify string) {
	if !isWildValid(topic) {
		return
	}

	handler.CxtMutex.Lock()
	for _, w := range handler.Wildcards {
		if w.isExist(topic, identify) {
			return
		}
	}

	wild := NewWild(topic, identify)
	handler.Wildcards = append(handler.Wildcards, *wild)
	handler.CxtMutex.Unlock()
}

func (handler *MqttSvrHandler) DelWild(topic, identify string) {
	newWilds := []Wild{}
	flag := false
	handler.CxtMutex.Lock()
	for _, w := range handler.Wildcards {
		if !w.isExist(topic, identify) {
			newWilds = append(newWilds, w)
		} else {
			flag = true
		}
	}

	if flag {
		handler.Wildcards = newWilds
	}
	handler.CxtMutex.Unlock()
}

func (handler *MqttSvrHandler) AddSub(topic string, identify string) {
	if len(identify) == 0 {
		return
	}

	handler.CxtMutex.Lock()
	ids, exist := handler.TopicToId[topic]
	if !exist {
		handler.TopicToId[topic] = []string{}
	}

	for i := 0; i < len(ids); i++ {
		if ids[i] == identify {
			return
		}
	}

	handler.TopicToId[topic] = append(handler.TopicToId[topic], identify)
	handler.CxtMutex.Unlock()

	return
}

func (handler *MqttSvrHandler) DelSub(topic string, identify string) {
	if len(identify) == 0 {
		return
	}

	handler.CxtMutex.Lock()
	ids, exist := handler.TopicToId[topic]
	if !exist {
		fmt.Printf("topic[%s] not exist\n", topic)
		return
	}

	newids := []string{}
	for i := 0; i < len(ids); i++ {
		if ids[i] == identify {
			continue
		}
		newids = append(newids, ids[i])
	}

	if len(newids) == 0 {
		delete(handler.TopicToId, topic)
	} else {
		handler.TopicToId[topic] = newids
	}

	handler.CxtMutex.Unlock()
	return
}

func (handler *MqttSvrHandler) PublishTopic(topic string, msg *Publish) {
	handler.CxtMutex.Lock()
	ids, exist := handler.TopicToId[topic]
	if !exist {
		return
	}

	newIds := []string{}
	for _, id := range ids {
		sub, exist := handler.IdToSub[id]
		if !exist {
			continue
		}
		fmt.Printf("submit msg to sub:%s\n", sub.Host)
		sub.Submit(msg)
		newIds = append(newIds, id)
	}
	handler.TopicToId[topic] = newIds

	newWilds := []Wild{}
	for _, w := range handler.Wildcards {
		if !w.matches(topic) {
			continue
		}
		sub, exist := handler.IdToSub[w.identify]
		if !exist {
			continue
		}
		fmt.Printf("submit msg to wilds:%s\n", sub.Host)
		sub.Submit(msg)
		newWilds = append(newWilds, w)
	}
	handler.Wildcards = newWilds
	handler.CxtMutex.Unlock()
}
