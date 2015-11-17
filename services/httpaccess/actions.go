package main

import (
	"github.com/PandoCloud/pando-cloud/pkg/models"
	"github.com/PandoCloud/pando-cloud/pkg/rpcs"
	"github.com/PandoCloud/pando-cloud/pkg/server"
	"github.com/martini-contrib/render"
	"net/http"
)

const (
	ErrOK          = 0
	ErrSystemFault = 10001
)

func renderError(code int, err error) Common {
	result := Common{}
	result.Code = code
	result.Message = err.Error()
	server.Log.Error(err.Error())
	return result
}

// device register args
type DeviceRegisterArgs struct {
	ProductKey string `json:"product_key"  binding:"required"`
	DeviceCode string `json:"device_code"  binding:"required"`
	Version    string `json:"version"  binding:"required"`
}

// device authentication args
type DeviceAuthArgs struct {
	DeviceId     int64  `json:"device_id" binding:"required"`
	DeviceSecret string `json:"device_secret" binding:"required"`
	Protocol     string `json:"protocol" binding:"required"`
}

func RegisterDevice(args DeviceRegisterArgs, r render.Render) {
	server.Log.Printf("ACTION RegisterDevice, args:: %v ", args)
	rpcargs := &rpcs.ArgsDeviceRegister{
		ProductKey:    args.ProductKey,
		DeviceCode:    args.DeviceCode,
		DeviceVersion: args.Version,
	}
	device := &models.Device{}
	err := server.RPCCallByName("registry", "Registry.RegisterDevice", rpcargs, device)
	if err != nil {
		r.JSON(http.StatusOK, renderError(ErrSystemFault, err))
		return
	}
	server.Log.Infof("register device success: %v", device)

	result := DeviceRegisterResponse{}
	result.Data = DeviceRegisterData{
		DeviceId:         device.ID,
		DeviceSecret:     device.DeviceSecret,
		DeviceKey:        device.DeviceKey,
		DeviceIdentifier: device.DeviceIdentifier,
	}
	r.JSON(http.StatusOK, result)
	return
}
