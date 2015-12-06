package protocol

import (
	"github.com/PandoCloud/pando-cloud/pkg/tlv"
	"reflect"
	"testing"
	"time"
)

func TestCommand(t *testing.T) {
	param := []interface{}{uint32(1), float32(3.2), []byte{'1', '2'}}
	params, err := tlv.MakeTLVs(param)
	if err != nil {
		t.Fatal(err)
	}

	payloadHead := CommandEventHead{
		Flag:        0,
		Timestamp:   uint64(time.Now().Unix()) * 1000,
		SubDeviceid: uint16(2),
		No:          uint16(12),
		Priority:    uint16(1),
		ParamsCount: uint16(len(param)),
	}
	payload := &Command{
		Head:   payloadHead,
		Params: params,
	}

	buf, err := payload.Marshal()
	if err != nil {
		if err != nil {
			t.Error(err)
		}
	}

	payload2 := &Command{}

	err = payload2.UnMarshal(buf)
	if err != nil {
		if err != nil {
			t.Error(err)
		}
	}

	if !reflect.DeepEqual(payload, payload2) {
		t.Errorf("test command payload failed, want %v, got %v", payload, payload2)
	}
}

func TestEvent(t *testing.T) {
	param := []interface{}{uint32(1), float32(3.2), []byte{'1', '2'}}
	params, err := tlv.MakeTLVs(param)
	if err != nil {
		t.Fatal(err)
	}

	payloadHead := CommandEventHead{
		Flag:        0,
		Timestamp:   uint64(time.Now().Unix()) * 1000,
		SubDeviceid: uint16(2),
		No:          uint16(12),
		Priority:    uint16(1),
		ParamsCount: uint16(len(param)),
	}
	payload := &Event{
		Head:   payloadHead,
		Params: params,
	}

	buf, err := payload.Marshal()
	if err != nil {
		if err != nil {
			t.Error(err)
		}
	}

	payload2 := &Event{}

	err = payload2.UnMarshal(buf)
	if err != nil {
		if err != nil {
			t.Error(err)
		}
	}

	if !reflect.DeepEqual(payload, payload2) {
		t.Errorf("test event payload failed, want %v, got %v", payload, payload2)
	}
}

func TestData(t *testing.T) {
	payloadHead := DataHead{
		Flag:      0,
		Timestamp: uint64(time.Now().Unix() * 1000),
	}
	param1 := []interface{}{uint32(3), float32(1.2), int64(10)}
	params1, err := tlv.MakeTLVs(param1)
	if err != nil {
		t.Fatal(err)
	}
	sub1 := SubData{
		Head: SubDataHead{
			SubDeviceid: uint16(1),
			PropertyNum: uint16(1),
			ParamsCount: uint16(len(params1)),
		},
		Params: params1,
	}
	param2 := []interface{}{uint32(4), int64(11)}
	params2, err := tlv.MakeTLVs(param2)
	if err != nil {
		t.Fatal(err)
	}
	sub2 := SubData{
		Head: SubDataHead{
			SubDeviceid: uint16(1),
			PropertyNum: uint16(2),
			ParamsCount: uint16(len(params2)),
		},
		Params: params2,
	}

	payload := &Data{
		Head:    payloadHead,
		SubData: []SubData{},
	}
	payload.SubData = append(payload.SubData, sub1)
	payload.SubData = append(payload.SubData, sub2)

	buf, err := payload.Marshal()
	if err != nil {
		if err != nil {
			t.Error(err)
		}
	}

	payload2 := &Data{}

	err = payload2.UnMarshal(buf)
	if err != nil {
		if err != nil {
			t.Error(err)
		}
	}

	if !reflect.DeepEqual(payload, payload2) {
		t.Errorf("test data payload failed, want %v, got %v", payload, payload2)
	}
}
