package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	"strings"
	"time"
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
	Url string

	// basic info
	ProductKey string
	DeviceCode string
	Version    string

	// private things
	id      int64
	secrect string
	token   []byte
	access  string
}

func NewDevice(url string, productkey string, code string, version string) *Device {
	return &Device{
		Url:        url,
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
	regUrl := fmt.Sprintf("%v%v", d.Url, "/v1/devices/registration")
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
		Protocol:     "mqtt",
	}
	regUrl := fmt.Sprintf("%v%v", d.Url, "/v1/devices/authentication")
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

func (d *Device) messageHandler(client *MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
	topicPieces := strings.Split(msg.Topic())
	clientid := topicPieces[0]
	msgtype := topicPieces[1]
}

func (d *Device) DoAccess() error {

	//create a ClientOptions struct setting the broker address, clientid, turn
	//off trace output and set the default message handler
	opts := MQTT.NewClientOptions().AddBroker("tcp://" + d.access)
	clientid := fmt.Sprintf("%x", d.id)
	opts.SetClientID(clientid)
	opts.SetUsername(clientid) // clientid as username
	opts.SetPassword(hex.EncodeToString(d.token))
	opts.SetKeepAlive(30 * time.Second)
	opts.SetDefaultPublishHandler(d.messageHandler)

	//create and start a client using the above ClientOptions
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	// we just pause here to wait for messages
	<-make(chan int)

	defer c.Disconnect(250)

	return nil
}
