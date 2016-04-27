package main

import (
	"encoding/json"
	"fmt"
	"github.com/PandoCloud/pando-cloud/pkg/models"
	"github.com/PandoCloud/pando-cloud/pkg/productconfig"
	"github.com/PandoCloud/pando-cloud/pkg/rpcs"
	"github.com/PandoCloud/pando-cloud/pkg/server"
	"strings"
)

func performRuleAction(target string, action string) error {
	server.Log.Infof("trigger rule action: %v, %v", target, action)

	parts := strings.Split(target, "/")
	if len(parts) != 3 {
		return fmt.Errorf("error target format: %v", target)
	}

	identifier := parts[1]
	device := &models.Device{}
	err := server.RPCCallByName("registry", "Registry.FindDeviceByIdentifier", identifier, device)
	if err != nil {
		return err
	}

	product := &models.Product{}
	err = server.RPCCallByName("registry", "Registry.FindProduct", device.ProductID, product)
	if err != nil {
		return err
	}

	config, err := productconfig.New(product.ProductConfig)
	if err != nil {
		return err
	}

	sendType := parts[2]
	switch sendType {
	case "command":
		var args interface{}
		err := json.Unmarshal([]byte(action), &args)
		if err != nil {
			server.Log.Errorf("marshal action error: %v", err)
			return err
		}

		server.Log.Debugf(": %v", args)

		m, ok := args.(map[string]interface{})
		if !ok {
			server.Log.Errorf("decode action error:%v", err)
			return fmt.Errorf("decode action error:%v", err)
		}

		command, err := config.MapToCommand(m)
		if err != nil {
			server.Log.Errorf("action format error: %v", err)
			return err
		}

		cmdargs := rpcs.ArgsSendCommand{
			DeviceId:  uint64(device.ID),
			SubDevice: uint16(command.Head.SubDeviceid),
			No:        uint16(command.Head.No),
			WaitTime:  uint32(3000),
			Params:    command.Params,
		}
		cmdreply := rpcs.ReplySendCommand{}
		err = server.RPCCallByName("controller", "Controller.SendCommand", cmdargs, &cmdreply)
		if err != nil {
			server.Log.Errorf("send device command error: %v", err)
			return err
		}
	case "status":

	default:
		server.Log.Errorf("wrong action %v", action)
	}

	return nil
}
