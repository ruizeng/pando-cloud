package main

import (
	"encoding/json"
	"github.com/PandoCloud/pando-cloud/pkg/models"
	"github.com/PandoCloud/pando-cloud/pkg/productconfig"
	"github.com/PandoCloud/pando-cloud/pkg/rpcs"
	"github.com/PandoCloud/pando-cloud/pkg/server"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
)

const (
	ErrOK                 = 0
	ErrSystemFault        = 10001
	ErrProductNotFound    = 10002
	ErrDeviceNotFound     = 10003
	ErrDeviceNotOnline    = 10004
	ErrWrongRequestFormat = 10005
	ErrWrongProductConfig = 10006
	ErrWrongQueryFormat   = 10007
	ErrAccessDenied       = 10008
)

const (
	defaultTimeOut = 3 // seconds
)

func renderError(code int, err error) Common {
	result := Common{}
	result.Code = code
	result.Message = err.Error()
	server.Log.Error(err.Error())
	return result
}

func GetDeviceInfoByKey(params martini.Params, req *http.Request, r render.Render) {
	key := req.URL.Query().Get("device_key")
	server.Log.Printf("ACTION GetDeviceInfoByKey, key:: %v", key)
	device := &models.Device{}
	err := server.RPCCallByName("registry", "Registry.ValidateDevice", key, device)
	if err != nil {
		r.JSON(http.StatusOK, renderError(ErrDeviceNotFound, err))
		return
	}

	result := DeviceInfoResponse{
		Data: DeviceInfoData{
			Identifier:  device.DeviceIdentifier,
			Name:        device.DeviceName,
			Description: device.DeviceDescription,
			Version:     device.DeviceVersion,
		},
	}
	r.JSON(http.StatusOK, result)
	return
}

func GetDeviceInfoByIdentifier(urlparams martini.Params, r render.Render) {
	identifier := urlparams["identifier"]
	server.Log.Printf("ACTION GetDeviceInfoByIdentifier, identifier:: %v", identifier)
	device := &models.Device{}
	err := server.RPCCallByName("registry", "Registry.FindDeviceByIdentifier", identifier, device)
	if err != nil {
		r.JSON(http.StatusOK, renderError(ErrDeviceNotFound, err))
		return
	}

	result := DeviceInfoResponse{
		Data: DeviceInfoData{
			Identifier:  device.DeviceIdentifier,
			Name:        device.DeviceName,
			Description: device.DeviceDescription,
			Version:     device.DeviceVersion,
		},
	}
	r.JSON(http.StatusOK, result)
	return
}

func GetDeviceCurrentStatus(device *models.Device, config *productconfig.ProductConfig,
	urlparams martini.Params, r render.Render) {
	server.Log.Printf("ACTION GetDeviceCurrentStatus, identifier:: %v", device.DeviceIdentifier)

	statusargs := rpcs.ArgsGetStatus{
		Id: uint64(device.ID),
	}
	statusreply := rpcs.ReplyGetStatus{}
	err := server.RPCCallByName("controller", "Controller.GetStatus", statusargs, &statusreply)
	if err != nil {
		server.Log.Errorf("get devie status error: %v", err)
		r.JSON(http.StatusOK, renderError(ErrSystemFault, err))
		return
	}

	status, err := config.StatusToMap(statusreply.Status)
	if err != nil {
		r.JSON(http.StatusOK, renderError(ErrWrongRequestFormat, err))
		return
	}
	result := DeviceStatusResponse{
		Data: status,
	}

	r.JSON(http.StatusOK, result)
	return
}

func GetDeviceLatestStatus() {

}

func SetDeviceStatus(device *models.Device, config *productconfig.ProductConfig,
	urlparams martini.Params, req *http.Request, r render.Render) {
	server.Log.Printf("ACTION GetDeviceCurrentStatus, identifier:: %v,request: %v", device.DeviceIdentifier, req.Body)

	var args interface{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&args)
	if err != nil {
		r.JSON(http.StatusOK, renderError(ErrWrongRequestFormat, err))
		return
	}

	m, ok := args.(map[string]interface{})
	if !ok {
		r.JSON(http.StatusOK, renderError(ErrWrongRequestFormat, err))
		return
	}

	status, err := config.MapToStatus(m)
	if err != nil {
		r.JSON(http.StatusOK, renderError(ErrWrongRequestFormat, err))
		return
	}

	statusargs := rpcs.ArgsSetStatus{
		DeviceId: uint64(device.ID),
		Status:   status,
	}
	statusreply := rpcs.ReplySetStatus{}
	err = server.RPCCallByName("controller", "Controller.SetStatus", statusargs, &statusreply)
	if err != nil {
		server.Log.Errorf("set devie status error: %v", err)
		r.JSON(http.StatusOK, renderError(ErrSystemFault, err))
		return
	}

	r.JSON(http.StatusOK, Common{})
	return

}

func SendCommandToDevice(device *models.Device, config *productconfig.ProductConfig,
	urlparams martini.Params, req *http.Request, r render.Render) {
	timeout := req.URL.Query().Get("timeout")

	server.Log.Printf("ACTION SendCommandToDevice, identifier:: %v, request: %v, timeout: %v",
		device.DeviceIdentifier, req.Body, timeout)

	var args interface{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&args)
	if err != nil {
		r.JSON(http.StatusOK, renderError(ErrWrongRequestFormat, err))
		return
	}

	m, ok := args.(map[string]interface{})
	if !ok {
		r.JSON(http.StatusOK, renderError(ErrWrongRequestFormat, err))
		return
	}

	command, err := config.MapToCommand(m)
	if err != nil {
		r.JSON(http.StatusOK, renderError(ErrWrongRequestFormat, err))
		return
	}

	cmdargs := rpcs.ArgsSendCommand{
		DeviceId:  uint64(device.ID),
		SubDevice: uint16(command.Head.SubDeviceid),
		No:        uint16(command.Head.No),
		WaitTime:  uint32(defaultTimeOut),
		Params:    command.Params,
	}
	cmdreply := rpcs.ReplySendCommand{}
	err = server.RPCCallByName("controller", "Controller.SendCommand", cmdargs, &cmdreply)
	if err != nil {
		server.Log.Errorf("send devie command error: %v", err)
		r.JSON(http.StatusOK, renderError(ErrSystemFault, err))
		return
	}

	r.JSON(http.StatusOK, Common{})
	return

}

func AddRule(device *models.Device, req *http.Request, r render.Render) {
	var ruleReq CreateRuleRequest
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&ruleReq)
	if err != nil {
		r.JSON(http.StatusOK, renderError(ErrWrongRequestFormat, err))
		return
	}

	rule := &models.Rule{
		DeviceID: device.ID,
		RuleType: ruleReq.Type,
		Trigger:  ruleReq.Trigger,
		Target:   ruleReq.Target,
		Action:   ruleReq.Action,
	}
	reply := &rpcs.ReplyEmptyResult{}

	err = server.RPCCallByName("registry", "Registry.CreateRule", rule, reply)
	if err != nil {
		server.Log.Errorf("create device rule error: %v", err)
		r.JSON(http.StatusOK, renderError(ErrSystemFault, err))
		return
	}

	r.JSON(http.StatusOK, Common{})
	return

}
