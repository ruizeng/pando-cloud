package main

import (
	"encoding/hex"
	"errors"
	"github.com/PandoCloud/pando-cloud/pkg/models"
	"github.com/PandoCloud/pando-cloud/pkg/rpcs"
	"github.com/PandoCloud/pando-cloud/pkg/server"
	"github.com/PandoCloud/pando-cloud/pkg/token"
	"github.com/martini-contrib/render"
	"math/rand"
	"net/http"
)

const (
	ErrOK                  = 0
	ErrSystemFault         = 10001
	ErrDeviceNotFound      = 10002
	ErrWrongSecret         = 10003
	ErrProtocolNotSuported = 10004
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

func AuthDevice(args DeviceAuthArgs, r render.Render) {
	server.Log.Printf("ACTION AuthDevice, args:: %v", args)
	device := &models.Device{}
	err := server.RPCCallByName("registry", "Registry.FindDeviceById", int64(args.DeviceId), device)
	if err != nil {
		r.JSON(http.StatusOK, renderError(ErrDeviceNotFound, err))
		return
	}

	if device.DeviceSecret != args.DeviceSecret {
		// device secret is wrong.
		r.JSON(http.StatusOK, renderError(ErrWrongSecret, errors.New("wrong device secret.")))
		return
	}

	hepler := token.NewHelper(*confRedisHost)
	token, err := hepler.GenerateToken(uint64(device.ID))
	if err != nil {
		r.JSON(http.StatusOK, renderError(ErrSystemFault, err))
		return
	}

	var hosts []string
	switch args.Protocol {
	case "http":
		hosts, err = server.GetServerHosts(args.Protocol+"access", "httphost")
	case "mqtt":
		hosts, err = server.GetServerHosts(args.Protocol+"access", "tcphost")
	default:
		err = errors.New("unsuported protocol: " + args.Protocol)
	}
	if err != nil {
		r.JSON(http.StatusOK, renderError(ErrProtocolNotSuported, err))
		return
	}

	// just get a random host
	host := hosts[rand.Intn(len(hosts))]

	result := DeviceAuthResponse{}
	result.Data = DeviceAuthData{
		AccessToken: hex.EncodeToString(token),
		AccessAddr:  host,
	}

	server.Log.Infof("auth device success: %v, token: %x, access: %v", device, token, host)

	r.JSON(http.StatusOK, result)
	return
}
