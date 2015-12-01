package mqtt

import (
	"time"
)

type BrokerTimerTask struct {
	context *Context
}

func NewBrokerTimerTask(ctx *Context) *BrokerTimerTask {
	return &BrokerTimerTask{context: ctx}
}

func (task *BrokerTimerTask) DoTask() {
	curTime := time.Now().Unix()

	for _, sub := range task.context.IdToSub {
		if sub.KeepAlive == 0 {
			continue
		}

		if uint16(curTime-sub.LastHbTime) > uint16(2*sub.KeepAlive/2) {
			sub.Close()
		}
	}
}
