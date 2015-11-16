package tlv

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func TestTlvs(t *testing.T) {
	params := []interface{}{int32(16), float32(3.2), []byte{'1', '2', '3'}}

	tlvs, err := MakeTLVs(params)
	if err != nil {
		t.Error(err)
	}

	newParams, err := ReadTLVs(tlvs)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(params, newParams) {
		t.Errorf("the origin:\n%x\n, now:\n%x\n", params, newParams)
	}
}

func TestTlvBinary(t *testing.T) {
	params := []interface{}{int32(15)}
	tlv, err := MakeTLV(params[0])
	if err != nil {
		t.Error(err)
	}

	bin := tlv.ToBinary()
	buf := bytes.NewReader(bin)
	newTlv := &TLV{}
	newTlv.FromBinary(buf)
	fmt.Printf("%x", bin)

	if !reflect.DeepEqual(tlv, newTlv) {
		t.Errorf("the origin:\n%x\n, now:\n%x\n", tlv, newTlv)
	}
}
