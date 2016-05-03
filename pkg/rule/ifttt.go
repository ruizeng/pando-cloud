// support ifttt action between two devices.
package rule

import (
	"github.com/PandoCloud/pando-cloud/pkg/models"
	"github.com/PandoCloud/pando-cloud/pkg/productconfig"
	"github.com/PandoCloud/pando-cloud/pkg/server"
)

type Ifttt struct{}

func NewIfttt() *Ifttt {
	return &Ifttt{}
}

func (ift *Ifttt) Check(deviceid uint64, eventid uint16) error {
	actions := &[]models.Rule{}
	query := &models.Rule{
		RuleType: "ifttt",
		DeviceID: int64(deviceid),
	}
	err := server.RPCCallByName("registry", "Registry.QueryRules", query, actions)
	if err != nil {
		server.Log.Warnf("load ifttt rules error : %v", err)
		return err
	}

	if len(*actions) > 0 {
		device := &models.Device{}
		err := server.RPCCallByName("registry", "Registry.FindDeviceById", int64(deviceid), device)
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

		name := ""
		for _, ev := range c.Events {
			if ev.No == int(eventid) {
				name = ev.Name
			}
		}

		for _, action := range *actions {
			if action.Trigger == name {
				err := performRuleAction(action.Target, action.Action)
				if err != nil {
					server.Log.Warnf("ifttt action failed: %v", err)
				}
			}
		}
	}

	return nil
}
