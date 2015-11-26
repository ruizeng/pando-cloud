package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
)

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

// common response fields
type Common struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// device register response data field
type DeviceRegisterData struct {
	DeviceId         int64  `json:"device_id"`
	DeviceSecret     string `json:"device_secret"`
	DeviceKey        string `json:"device_key"`
	DeviceIdentifier string `json:"device_identifier"`
}

// device register response
type DeviceRegisterResponse struct {
	Common
	Data DeviceRegisterData `json:"data"`
}

// device auth response data field
type DeviceAuthData struct {
	AccessToken string `json:"access_token"`
	AccessAddr  string `json:"access_addr"`
}

// device auth response
type DeviceAuthResponse struct {
	Common
	Data DeviceAuthData `json:"data"`
}

type Device struct {
	// API URL
	BrokerUrl string

	// basic info
	ProductKey string
	DeviceCode string
	Version    string

	// private thins
	id      int64
	secrect string
	token   []byte
	access  string
}

func NewDevice(broker string, productkey string, code string, version string) *Device {
	return &Device{
		BrokerUrl:  broker,
		ProductKey: productkey,
		DeviceCode: code,
		Version:    version,
	}
}

func (d *Device) DoRegister() error {
	args := DeviceRegisterArgs{
		ProductKey: d.ProductKey,
		DeviceCode: d.DeviceCode,
		Version:    d.Version,
	}
	regUrl := fmt.Sprintf("%v%v", d.BrokerUrl, "/v1/devices/registration")
	request, err := json.Marshal(args)
	if err != nil {
		return err
	}
	jsonresp, err := SendHttpRequest(regUrl, string(request), "POST", nil)
	if err != nil {
		return err
	}
	response := DeviceRegisterResponse{}
	err = json.Unmarshal(jsonresp, &response)
	if err != nil {
		return err
	}
	err = CheckHttpsCode(response)
	if err != nil {
		return err
	}

	d.id = response.Data.DeviceId
	d.secrect = response.Data.DeviceSecret

	return nil
}

func (d *Device) DoLogin() error {
	args := DeviceAuthArgs{
		DeviceId:     d.id,
		DeviceSecret: d.secrect,
		Protocol:     "http",
	}
	regUrl := fmt.Sprintf("%v%v", d.BrokerUrl, "/v1/devices/authentication")
	request, err := json.Marshal(args)
	if err != nil {
		return err
	}
	jsonresp, err := SendHttpRequest(regUrl, string(request), "POST", nil)
	if err != nil {
		return err
	}
	response := DeviceAuthResponse{}
	err = json.Unmarshal(jsonresp, &response)
	if err != nil {
		return err
	}
	err = CheckHttpsCode(response)
	if err != nil {
		return err
	}
	// ecode hex
	htoken, err := hex.DecodeString(response.Data.AccessToken)
	if err != nil {
		return err
	}
	d.token = htoken
	d.access = response.Data.AccessAddr

	return nil
}
