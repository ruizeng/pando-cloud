package productconfig

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PandoCloud/pando-cloud/pkg/protocol"
	"github.com/PandoCloud/pando-cloud/pkg/tlv"
)

type CommandOrEventParam struct {
	ValueType int32 `json:"value_type"`
	Name      string
}

type ProductCommandOrEvent struct {
	No       int
	Part     int
	Name     string
	Priority int
	Params   []CommandOrEventParam
}

type StatusParam struct {
	ValueType int32 `json:"value_type"`
	Name      string
}

type ProductObject struct {
	Id     int
	No     int
	Label  string
	Part   int
	Status []StatusParam
}

// product config parses the JSON product config string.
type ProductConfig struct {
	Objects  []ProductObject
	Commands []ProductCommandOrEvent
	Events   []ProductCommandOrEvent
}

func New(config string) (*ProductConfig, error) {
	v := &ProductConfig{}
	err := json.Unmarshal([]byte(config), v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (config *ProductConfig) ValidateStatus(label string, params []interface{}) (*ProductObject, []interface{}, error) {
	// search for status name
	var paramInfo []StatusParam
	var status *ProductObject
	found := false
	for _, obj := range config.Objects {
		if obj.Label == label {
			paramInfo = obj.Status
			status = &obj
			found = true
			break
		}
	}
	if found == false {
		return nil, []interface{}{}, errors.New("object not found.")
	}
	if len(paramInfo) != len(params) {
		return nil, []interface{}{}, errors.New("wrong status parameters.")
	}
	realParams := make([]interface{}, len(params))
	for idx, para := range paramInfo {
		realParams[idx] = tlv.CastTLV(params[idx], para.ValueType)
	}
	return status, realParams, nil
}

func (config *ProductConfig) ValidateCommandOrEvent(name string, params []interface{}, typ string) (*ProductCommandOrEvent, []interface{}, error) {
	var target []ProductCommandOrEvent
	if typ == "command" {
		target = config.Commands
	} else if typ == "event" {
		target = config.Events
	} else {
		return nil, []interface{}{}, errors.New("wrong target type.")
	}

	// search for name
	var paramInfo []CommandOrEventParam
	var coe *ProductCommandOrEvent
	found := false
	for _, one := range target {
		if one.Name == name {
			paramInfo = one.Params
			coe = &one
			found = true
			break
		}
	}
	if found == false {
		return nil, []interface{}{}, errors.New("command or event not found.")
	}
	if len(paramInfo) != len(params) {
		return nil, []interface{}{}, errors.New("wrong parameters.")
	}
	realParams := make([]interface{}, len(params))
	for idx, para := range paramInfo {
		realParams[idx] = tlv.CastTLV(params[idx], para.ValueType)
	}
	return coe, realParams, nil
}

func (config *ProductConfig) StatusToMap(status []protocol.SubData) (map[string][]interface{}, error) {
	result := make(map[string][]interface{})

	for _, sub := range status {
		val, err := tlv.ReadTLVs(sub.Params)
		if err != nil {
			return nil, err
		}
		label := ""
		for _, obj := range config.Objects {
			if obj.No == int(sub.Head.PropertyNum) {
				label = obj.Label
			}
		}
		result[label] = val
	}

	return result, nil
}

func (config *ProductConfig) EventToMap(event *protocol.Event) (map[string][]interface{}, error) {
	result := make(map[string][]interface{})

	name := ""
	for _, ev := range config.Events {
		if ev.No == int(event.Head.No) {
			name = ev.Name
		}
	}
	val, err := tlv.ReadTLVs(event.Params)
	if err != nil {
		return nil, err
	}

	result[name] = val

	return result, nil
}

func (config *ProductConfig) MapToStatus(data map[string]interface{}) ([]protocol.SubData, error) {
	result := []protocol.SubData{}

	for label, one := range data {
		params, ok := one.([]interface{})
		if !ok {
			return nil, fmt.Errorf("status format error: %v", one)
		}
		obj, realParams, err := config.ValidateStatus(label, params)
		if err != nil {
			return nil, err
		}

		tlvs, err := tlv.MakeTLVs(realParams)
		if err != nil {
			return nil, err
		}

		result = append(result, protocol.SubData{
			Head: protocol.SubDataHead{
				SubDeviceid: uint16(obj.Part),
				PropertyNum: uint16(obj.No),
				ParamsCount: uint16(len(realParams)),
			},
			Params: tlvs,
		})
	}

	return result, nil
}

func (config *ProductConfig) MapToCommand(cmd map[string]interface{}) (*protocol.Command, error) {
	result := &protocol.Command{}

	for name, one := range cmd {
		params, ok := one.([]interface{})
		if !ok {
			return nil, fmt.Errorf("command format error: %v", one)
		}

		c, realParams, err := config.ValidateCommandOrEvent(name, params, "command")
		if err != nil {
			return nil, err
		}

		tlvs, err := tlv.MakeTLVs(realParams)
		if err != nil {
			return nil, err
		}

		result.Head.No = uint16(c.No)
		result.Head.Priority = uint16(c.Priority)
		result.Head.SubDeviceid = uint16(c.Part)
		result.Head.ParamsCount = uint16(len(realParams))
		result.Params = tlvs

	}

	return result, nil
}
