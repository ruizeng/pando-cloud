package productconfig

import (
	"encoding/json"
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

	c, err := NewProductConfig(config)
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

}
