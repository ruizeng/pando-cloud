package main

import (
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	"github.com/PandoCloud/pando-cloud/pkg/protocol"
	"github.com/PandoCloud/pando-cloud/pkg/tlv"
	"log"
	"os"
	"time"
)

const (
	commonCmdGetStatus = uint16(65528)
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

func (d *Device) reportStatus(client *MQTT.Client) {

	payloadHead := protocol.DataHead{
		Flag:      0,
		Timestamp: uint64(time.Now().Unix() * 1000),
	}
	param := []interface{}{uint8(1)}
	params, err := tlv.MakeTLVs(param)
	if err != nil {
		fmt.Println(err)
		return
	}
	sub := protocol.SubData{
		Head: protocol.SubDataHead{
			SubDeviceid: uint16(1),
			PropertyNum: uint16(1),
			ParamsCount: uint16(len(params)),
		},
		Params: params,
	}

	status := protocol.Data{
		Head:    payloadHead,
		SubData: []protocol.SubData{},
	}

	status.SubData = append(status.SubData, sub)

	payload, err := status.Marshal()
	if err != nil {
		fmt.Println(err)
		return
	}

	client.Publish("s", 1, false, payload)

}

func (d *Device) statusHandler(client *MQTT.Client, msg MQTT.Message) {
	status := protocol.Data{}

	err := status.UnMarshal(msg.Payload())
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("device receiving status set : ")

	for _, one := range status.SubData {
		fmt.Println("subdeviceid : ", one.Head.SubDeviceid)
		fmt.Println("no : ", one.Head.PropertyNum)
		fmt.Println("params : ", one.Params)
	}
}

func (d *Device) commandHandler(client *MQTT.Client, msg MQTT.Message) {
	cmd := protocol.Command{}

	err := cmd.UnMarshal(msg.Payload())
	if err != nil {
		fmt.Println(err)
		return
	}

	switch cmd.Head.No {
	case commonCmdGetStatus:
		d.reportStatus(client)
	default:
		fmt.Println("unsuported command : %v", cmd.Head.No)
	}
}

func (d *Device) messageHandler(client *MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %x\n", msg.Payload())
	msgtype := msg.Topic()
	fmt.Println(msgtype)

	switch msgtype {
	case "c":
		d.commandHandler(client, msg)
	case "s":
		d.statusHandler(client, msg)
	default:
		fmt.Println("unsuported message type :", msgtype)
	}
}

func (d *Device) DoAccess() error {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	MQTT.ERROR = logger
	MQTT.CRITICAL = logger
	MQTT.WARN = logger
	MQTT.DEBUG = logger

	//create a ClientOptions struct setting the broker address, clientid, turn
	//off trace output and set the default message handler
	opts := MQTT.NewClientOptions().AddBroker("tls://" + d.access)
	clientid := fmt.Sprintf("%x", d.id)
	opts.SetClientID(clientid)
	opts.SetUsername(clientid) // clientid as username
	opts.SetPassword(hex.EncodeToString(d.token))
	opts.SetKeepAlive(30 * time.Second)
	opts.SetDefaultPublishHandler(d.messageHandler)
	opts.SetTLSConfig(&tls.Config{Certificates: nil, InsecureSkipVerify: true})

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
