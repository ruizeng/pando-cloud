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
