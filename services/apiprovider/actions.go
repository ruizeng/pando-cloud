package main

import (
	"github.com/PandoCloud/pando-cloud/pkg/models"
	"github.com/PandoCloud/pando-cloud/pkg/rpcs"
	"github.com/PandoCloud/pando-cloud/pkg/server"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
)

const (
	ErrOK             = 0
	ErrSystemFault    = 10001
	ErrDeviceNotFound = 10002
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

	result := DeviceStatusResponse{}

	onlineargs := rpcs.ArgsGetDeviceStatus{
		Id: uint64(device.ID),
	}
	onlinereply := rpcs.ReplyGetDeviceStatus{}
	err = server.RPCCallByName("devicemanger", "DeviceManager.GetDeviceStatus", onlineargs, &onlinereply)
	if err != nil {
		// if device is not online, just return
		server.Log.Errorf("get devie status error: %v", err)
		r.JSON(http.StatusOK, result)
		return
	}

	// device is online, try read status
	// todo
}

func GetDeviceLatestStatus() {

}

func SetDeviceStatus() {

}
