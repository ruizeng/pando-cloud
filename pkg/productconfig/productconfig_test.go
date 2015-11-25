package productconfig

import (
	"encoding/json"
	"github.com/PandoCloud/pando-cloud/pkg/protocol"
	"github.com/PandoCloud/pando-cloud/pkg/tlv"
	"testing"
)

func TestParseProductConfig(t *testing.T) {
	config :=
		`
    {
      "objects": [{
        "id": 2,
        "no": 1,
        "label": "switch",
        "part": 1,
        "status": [{
          "value_type": 7,
          "name": "onoff"
        }]
      }],
      "commands": [{
        "no": 1,
        "part": 1,
        "name": "switch",
        "priority": 0,
        "params": [{
          "value_type": 7,
          "name": "status"
        }]
      }],
      "events": []
    }
  `

	status :=
		`
    {
      "switch": [1]
    }
    `

	c, err := New(config)
	if err != nil {
		t.Fatal(err)
	}

	var v interface{}
	err = json.Unmarshal([]byte(status), &v)
	if err != nil {
		t.Fatal(err)
	}

	for label, onedata := range v.(map[string]interface{}) {
		params := onedata.([]interface{})
		obj, realParams, err := c.ValidateStatus(label, params)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(obj)
		t.Log(realParams)
	}

	one, err := tlv.MakeTLV(uint8(1))
	if err != nil {
		if err != nil {
			t.Error(err)
		}
	}

	params := []tlv.TLV{*one}
	teststatus := []protocol.SubData{protocol.SubData{
		Head: protocol.SubDataHead{
			SubDeviceid: uint16(1),
			PropertyNum: uint16(1),
			ParamsCount: uint16(1),
		},
		Params: params,
	}}

	res, err := c.StatusToMap(teststatus)
	if err != nil {
		t.Error(err)
	}

	t.Log(res)

	m := make(map[string][]interface{})
	m["switch"] = []interface{}{float64(1)}
	_, err = c.MapToStatus(m)
	if err != nil {
		t.Error(err)
	}

}
