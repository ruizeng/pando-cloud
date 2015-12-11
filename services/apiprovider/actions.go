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
	ErrWrongStatusFormat  = 10005
	ErrWrongProductConfig = 10006
	ErrWrongQueryFormat   = 10007
	ErrAccessDenied       = 10008
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
	err := server.RPCCallByName("registry", "Registry.FindDeviceByKey", key, device)
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

func GetDeviceCurrentStatus(urlparams martini.Params, r render.Render) {
	identifier := urlparams["identifier"]
	server.Log.Printf("ACTION GetDeviceCurrentStatus, identifier:: %v", identifier)

	device := &models.Device{}
	err := server.RPCCallByName("registry", "Registry.FindDeviceByIdentifier", identifier, device)
	if err != nil {
		r.JSON(http.StatusOK, renderError(ErrDeviceNotFound, err))
		return
	}

	onlineargs := rpcs.ArgsGetDeviceOnlineStatus{
		Id: uint64(device.ID),
	}
	onlinereply := rpcs.ReplyGetDeviceOnlineStatus{}
	err = server.RPCCallByName("devicemanager", "DeviceManager.GetDeviceOnlineStatus", onlineargs, &onlinereply)
	if err != nil {
		server.Log.Errorf("get devie online status error: %v", err)
		r.JSON(http.StatusOK, renderError(ErrDeviceNotOnline, err))
		return
	}

	statusargs := rpcs.ArgsGetStatus{
		Id: uint64(device.ID),
	}
	statusreply := rpcs.ReplyGetStatus{}
	err = server.RPCCallByName("controller", "Controller.GetStatus", statusargs, &statusreply)
	if err != nil {
		server.Log.Errorf("get devie status error: %v", err)
		r.JSON(http.StatusOK, renderError(ErrSystemFault, err))
		return
	}

	product := &models.Product{}
	err = server.RPCCallByName("registry", "Registry.FindProduct", device.ProductID, product)
	if err != nil {
		r.JSON(http.StatusOK, renderError(ErrProductNotFound, err))
		return
	}

	c, err := productconfig.New(product.ProductConfig)
	if err != nil {
		r.JSON(http.StatusOK, renderError(ErrWrongProductConfig, err))
		return
	}

	data, err := c.StatusToMap(statusreply.Status)
	if err != nil {
		r.JSON(http.StatusOK, renderError(ErrWrongStatusFormat, err))
		return
	}
	result := DeviceStatusResponse{
		Data: data,
	}

	r.JSON(http.StatusOK, result)
	return
}

func GetDeviceLatestStatus() {

}

func SetDeviceStatus(urlparams martini.Params, req *http.Request, r render.Render) {
	identifier := urlparams["identifier"]
	server.Log.Printf("ACTION GetDeviceCurrentStatus, identifier:: %v", identifier)

	device := &models.Device{}
	err := server.RPCCallByName("registry", "Registry.FindDeviceByIdentifier", identifier, device)
	if err != nil {
		r.JSON(http.StatusOK, renderError(ErrDeviceNotFound, err))
		return
	}

	onlineargs := rpcs.ArgsGetDeviceOnlineStatus{
		Id: uint64(device.ID),
	}
	onlinereply := rpcs.ReplyGetDeviceOnlineStatus{}
	err = server.RPCCallByName("devicemanager", "DeviceManager.GetDeviceOnlineStatus", onlineargs, &onlinereply)
	if err != nil {
		server.Log.Errorf("get devie online status error: %v", err)
		r.JSON(http.StatusOK, renderError(ErrDeviceNotOnline, err))
		return
	}

	var args interface{}
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&args)
	if err != nil {
		r.JSON(http.StatusOK, renderError(ErrWrongStatusFormat, err))
		return
	}

	product := &models.Product{}
	err = server.RPCCallByName("registry", "Registry.FindProduct", device.ProductID, product)
	if err != nil {
		r.JSON(http.StatusOK, renderError(ErrProductNotFound, err))
		return
	}

	c, err := productconfig.New(product.ProductConfig)
	if err != nil {
		r.JSON(http.StatusOK, renderError(ErrWrongProductConfig, err))
		return
	}

	m, ok := args.(map[string]interface{})
	if !ok {
		r.JSON(http.StatusOK, renderError(ErrWrongStatusFormat, err))
		return
	}

	status, err := c.MapToStatus(m)
	if err != nil {
		r.JSON(http.StatusOK, renderError(ErrWrongStatusFormat, err))
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
