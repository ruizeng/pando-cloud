package productconfig

import (
	"encoding/json"
	"github.com/PandoCloud/pando-cloud/pkg/protocol"
	"github.com/PandoCloud/pando-cloud/pkg/tlv"
	"reflect"
	"testing"
)

func testStatus(c *ProductConfig, t *testing.T) {
	status :=
		`
    {
      "switch": [1]
    }
    `

	var v interface{}
	err := json.Unmarshal([]byte(status), &v)
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

	m := make(map[string]interface{})
	m["switch"] = []interface{}{float64(1)}
	_, err = c.MapToStatus(m)
	if err != nil {
		t.Error(err)
	}
}

func testEvent(c *ProductConfig, t *testing.T) {
	want := `{"alarm":["test"]}`

	testev := &protocol.Event{}
	testev.Head.No = 1
	testev.Head.SubDeviceid = 1
	params, err := tlv.MakeTLVs([]interface{}{"test"})
	if err != nil {
		t.Error(err)
	}
	testev.Params = params

	m, err := c.EventToMap(testev)
	if err != nil {
		t.Error(err)
	}

	result, err := json.Marshal(m)
	if err != nil {
		t.Error(err)
	}

	got := string(result)

	if got != want {
		t.Errorf("event to map error: want: %v, got : %v", want, got)
	}

}

func testCommand(c *ProductConfig, t *testing.T) {
	input := `{"switch":[1,2]}`

	v := make(map[string]interface{})
	err := json.Unmarshal([]byte(input), &v)
	if err != nil {
		t.Fatal(err)
	}
	params, err := tlv.MakeTLVs([]interface{}{uint8(1), uint8(2)})
	want := &protocol.Command{}
	want.Head.No = 1
	want.Head.SubDeviceid = 1
	want.Head.ParamsCount = 2
	want.Params = params

	got, err := c.MapToCommand(v)

	if !reflect.DeepEqual(want, got) {
		t.Errorf("map to command error: want: %v, got %v", want, got)
	}
}

func TestParseProductConfig(t *testing.T) {
	config :=
		`
    {
      "objects": [{
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
          "name": "p1"
        },{
          "value_type": 7,
          "name": "p2"
        }]
      }],
      "events": [{
        "no": 1,
        "part": 1,
        "name": "alarm",
        "priority": 0,
        "params": [{
          "value_type": 12,
          "name": "text"
        }]
      }]
    }
    `

	c, err := New(config)
	if err != nil {
		t.Fatal(err)
	}

	testStatus(c, t)
	testEvent(c, t)
	testCommand(c, t)

}
