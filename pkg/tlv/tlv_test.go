package tlv

import (
	"bytes"
	"reflect"
	"testing"
)

func TestTlvLen(t *testing.T) {
	float64Tlv, _ := MakeTLV(float64(0.12))
	if float64Tlv.Length() != 10 {
		t.Errorf("float64 len is not right\n")
	}

	int64Tlv, _ := MakeTLV(int64(100))
	if int64Tlv.Length() != 10 {
		t.Errorf("int64 len is not right\n")
	}

	uint64Tlv, _ := MakeTLV(uint64(100))
	if uint64Tlv.Length() != 10 {
		t.Errorf("uint64 len is not right\n")
	}

	float32Tlv, _ := MakeTLV(float32(0.12))
	if float32Tlv.Length() != 6 {
		t.Errorf("float32 len is not right\n")
	}

	int32Tlv, _ := MakeTLV(int32(100))
	if int32Tlv.Length() != 6 {
		t.Errorf("int32 len is not right\n")
	}

	uint32Tlv, _ := MakeTLV(uint32(100))
	if uint32Tlv.Length() != 6 {
		t.Errorf("uint32 len is not right\n")
	}

	int16Tlv, _ := MakeTLV(int16(100))
	if int16Tlv.Length() != 4 {
		t.Errorf("int16 len is not right\n")
	}

	uint16Tlv, _ := MakeTLV(uint16(100))
	if uint16Tlv.Length() != 4 {
		t.Errorf("uint16 len is not right\n")
	}

	int8Tlv, _ := MakeTLV(int8(100))
	if int8Tlv.Length() != 3 {
		t.Errorf("int8 len is not right\n")
	}

	uint8Tlv, _ := MakeTLV(uint8(100))
	if uint8Tlv.Length() != 3 {
		t.Errorf("uint8 len is not right\n")
	}

	byteValue := []byte{'1', '0', '0'}
	byteTLV, _ := MakeTLV(byteValue)
	if byteTLV.Length() != len(byteValue)+4 {
		t.Errorf("byte len is not right\n")
	}

	str := "100"
	strTLV, _ := MakeTLV(str)
	if strTLV.Length() != len(str)+4 {
		t.Errorf("string len is not right\n")
	}
}

func TestUintAndByte(t *testing.T) {
	value := uint16(100)
	byteValue := Uint16ToByte(value)
	newValue := ByteToUint16(byteValue)

	if value != newValue {
		t.Errorf("origin: %d, now: %d\n", value, newValue)
	}
}

func TestTlvs(t *testing.T) {
	str := "itachili"
	params := []interface{}{float64(0.1), int64(100), uint64(200), uint32(300), int32(16), float32(3.2), int16(20), uint16(30), int8(1), uint8(2), []byte{'1', '2', '3'}, str}

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
	str := "itachili"
	params := []interface{}{float64(0.1), int64(100), uint64(200), uint32(300), int32(16), float32(3.2), int16(20), uint16(30), int8(1), uint8(2), []byte{'1', '2', '3'}, str}
	tlv, err := MakeTLV(params[0])
	if err != nil {
		t.Error(err)
	}

	bin := tlv.ToBinary()
	buf := bytes.NewReader(bin)
	newTlv := &TLV{}
	newTlv.FromBinary(buf)

	if !reflect.DeepEqual(tlv, newTlv) {
		t.Errorf("the origin:\n%x\n, now:\n%x\n", tlv, newTlv)
	}
}
