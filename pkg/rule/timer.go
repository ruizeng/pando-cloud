// suport cron like schedule tasks.
package rule

import (
	"fmt"
	"github.com/PandoCloud/pando-cloud/pkg/models"
	"github.com/PandoCloud/pando-cloud/pkg/server"
	"github.com/robfig/cron"
	"time"
)

type Timer struct {
	c *cron.Cron
}

func NewTimer() *Timer {
	t := &Timer{}

	return t
}

func (t *Timer) createTimerFunc(target string, action string) func() {
	return func() {
		err := performRuleAction(target, action)
		if err != nil {
			server.Log.Warnf("timer action failed: %v", err)
		}
	}
}

func (t *Timer) refresh() {
	if t.c != nil {
		t.c.Stop()
	}
	t.c = cron.New()
	timers := &[]models.Rule{}
	query := &models.Rule{
		RuleType: "timer",
	}
	err := server.RPCCallByName("registry", "Registry.QueryRules", query, timers)
	if err != nil {
		server.Log.Warnf("refresh timer rules error : %v", err)
		return
	}

	sec := fmt.Sprintf("%d ", (time.Now().Second()+30)%60)

	for _, one := range *timers {
		t.c.AddFunc(sec+one.Trigger, t.createTimerFunc(one.Target, one.Action))
	}

	t.c.Start()
}

func (t *Timer) Run() {
	go func() {
		for {
			t.refresh()
			time.Sleep(time.Minute)
		}
	}()
}
