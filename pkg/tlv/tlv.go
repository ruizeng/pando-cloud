package tlv

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	TLV_FLOAT64 = 1
	TLV_FLOAT32 = 2
	TLV_INT8    = 3
	TLV_INT16   = 4
	TLV_INT32   = 5
	TLV_INT64   = 6
	TLV_UINT8   = 7
	TLV_UINT16  = 8
	TLV_UINT32  = 9
	TLV_UINT64  = 10
	TLV_BYTES   = 11
	TLV_URI     = 12
	TLV_BOOL    = 13
)

type TLV struct {
	Tag   uint16
	Value []byte
}

func Uint16ToByte(value uint16) []byte {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, value)

	return buf.Bytes()
}

func ByteToUint16(buf []byte) uint16 {
	tmpBuf := bytes.NewBuffer(buf)
	var value uint16
	binary.Read(tmpBuf, binary.BigEndian, &value)

	return value
}

func (tlv *TLV) ToBinary() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, &tlv.Tag)
	binary.Write(buf, binary.BigEndian, &tlv.Value)

	return buf.Bytes()
}

func (tlv *TLV) Length() int {
	length := int(0)
	switch tlv.Tag {
	case TLV_FLOAT64:
		length = 8
	case TLV_INT64:
		length = 8
	case TLV_UINT64:
		length = 8
	case TLV_FLOAT32:
		length = 4
	case TLV_INT32:
		length = 4
	case TLV_UINT32:
		length = 4
	case TLV_INT16:
		length = 2
	case TLV_UINT16:
		length = 2
	case TLV_INT8:
		length = 1
	case TLV_UINT8:
		length = 1
	case TLV_BYTES:
		length = int(ByteToUint16(tlv.Value[0:2]))
		length += 2
	case TLV_URI:
		length = int(ByteToUint16(tlv.Value[0:2]))
		length += 2
	default:
		length = 0
	}

	length += 2

	return length
}

func (tlv *TLV) FromBinary(r io.Reader) error {
	binary.Read(r, binary.BigEndian, &tlv.Tag)
	length := uint16(0)
	switch tlv.Tag {
	case TLV_FLOAT64:
		length = 8
		tlv.Value = make([]byte, length)
		binary.Read(r, binary.BigEndian, &tlv.Value)
	case TLV_INT64:
		length = 8
		tlv.Value = make([]byte, length)
		binary.Read(r, binary.BigEndian, &tlv.Value)
	case TLV_UINT64:
		length = 8
		tlv.Value = make([]byte, length)
		binary.Read(r, binary.BigEndian, &tlv.Value)
	case TLV_FLOAT32:
		length = 4
		tlv.Value = make([]byte, length)
		binary.Read(r, binary.BigEndian, &tlv.Value)
	case TLV_INT32:
		length = 4
		tlv.Value = make([]byte, length)
		binary.Read(r, binary.BigEndian, &tlv.Value)
	case TLV_UINT32:
		length = 4
		tlv.Value = make([]byte, length)
		binary.Read(r, binary.BigEndian, &tlv.Value)
	case TLV_INT16:
		length = 2
		tlv.Value = make([]byte, length)
		binary.Read(r, binary.BigEndian, &tlv.Value)
	case TLV_UINT16:
		length = 2
		tlv.Value = make([]byte, length)
		binary.Read(r, binary.BigEndian, &tlv.Value)
	case TLV_INT8:
		length = 1
		tlv.Value = make([]byte, length)
		binary.Read(r, binary.BigEndian, &tlv.Value)
	case TLV_UINT8:
		length = 1
		tlv.Value = make([]byte, length)
		binary.Read(r, binary.BigEndian, &tlv.Value)
	case TLV_BYTES:
		binary.Read(r, binary.BigEndian, &length)
		tlv.Value = make([]byte, length+2)
		copy(tlv.Value[0:2], Uint16ToByte(length))
		binary.Read(r, binary.BigEndian, tlv.Value[2:])
	case TLV_URI:
		binary.Read(r, binary.BigEndian, &length)
		tlv.Value = make([]byte, length+2)
		copy(tlv.Value[0:2], Uint16ToByte(length))
		binary.Read(r, binary.BigEndian, tlv.Value[2:])
	default:
		return errors.New(fmt.Sprintf("unsuport value: %d", tlv.Tag))
	}

	return nil
}

func MakeTLV(a interface{}) (*TLV, error) {
	var tag uint16
	var length uint16
	buf := new(bytes.Buffer)
	switch a.(type) {
	case float64:
		tag = TLV_FLOAT64
		length = 8
		binary.Write(buf, binary.BigEndian, a.(float64))
	case float32:
		tag = TLV_FLOAT32
		length = 4
		binary.Write(buf, binary.BigEndian, a.(float32))
	case int8:
		tag = TLV_INT8
		length = 1
		binary.Write(buf, binary.BigEndian, a.(int8))
	case int16:
		tag = TLV_INT16
		length = 2
		binary.Write(buf, binary.BigEndian, a.(int16))
	case int32:
		tag = TLV_INT32
		length = 4
		binary.Write(buf, binary.BigEndian, a.(int32))
	case int64:
		tag = TLV_INT64
		length = 8
		binary.Write(buf, binary.BigEndian, a.(int64))
	case uint8:
		tag = TLV_UINT8
		length = 1
		binary.Write(buf, binary.BigEndian, a.(uint8))
	case uint16:
		tag = TLV_UINT16
		length = 2
		binary.Write(buf, binary.BigEndian, a.(uint16))
	case uint32:
		tag = TLV_UINT32
		length = 4
		binary.Write(buf, binary.BigEndian, a.(uint32))
	case uint64:
		tag = TLV_UINT64
		length = 8
		binary.Write(buf, binary.BigEndian, a.(uint64))
	case []byte:
		tag = TLV_BYTES
		length = uint16(len(a.([]byte)))
		binary.Write(buf, binary.BigEndian, length)
		binary.Write(buf, binary.BigEndian, a.([]byte))
	case string:
		tag = TLV_URI
		length = uint16(len(a.(string)))
		binary.Write(buf, binary.BigEndian, length)
		binary.Write(buf, binary.BigEndian, []byte(a.(string)))
	default:
		return nil, errors.New(fmt.Sprintf("unsuport value: %v", a))
	}

	tlv := TLV{
		Tag:   tag,
		Value: buf.Bytes(),
	}

	if length == 0 {
		tlv.Value = []byte{}
	}

	return &tlv, nil
}

func ReadTLV(tlv *TLV) (interface{}, error) {
	tag := tlv.Tag
	length := uint16(0)
	value := tlv.Value

	buffer := bytes.NewReader(value)
	switch tag {
	case TLV_FLOAT64:
		retvar := float64(0.0)
		err := binary.Read(buffer, binary.BigEndian, &retvar)
		return retvar, err
	case TLV_FLOAT32:
		retvar := float32(0.0)
		err := binary.Read(buffer, binary.BigEndian, &retvar)
		return retvar, err
	case TLV_INT8:
		retvar := int8(0)
		err := binary.Read(buffer, binary.BigEndian, &retvar)
		return retvar, err
	case TLV_INT16:
		retvar := int16(0)
		err := binary.Read(buffer, binary.BigEndian, &retvar)
		return retvar, err
	case TLV_INT32:
		retvar := int32(0)
		err := binary.Read(buffer, binary.BigEndian, &retvar)
		return retvar, err
	case TLV_INT64:
		retvar := int64(0)
		err := binary.Read(buffer, binary.BigEndian, &retvar)
		return retvar, err
	case TLV_UINT8:
		retvar := uint8(0)
		err := binary.Read(buffer, binary.BigEndian, &retvar)
		return retvar, err
	case TLV_UINT16:
		retvar := uint16(0)
		err := binary.Read(buffer, binary.BigEndian, &retvar)
		return retvar, err
	case TLV_UINT32:
		retvar := uint32(0)
		err := binary.Read(buffer, binary.BigEndian, &retvar)
		return retvar, err
	case TLV_UINT64:
		retvar := uint64(0)
		err := binary.Read(buffer, binary.BigEndian, &retvar)
		return retvar, err
	case TLV_BYTES:
		err := binary.Read(buffer, binary.BigEndian, &length)
		if err != nil {
			return []byte{}, err
		}
		retvar := make([]byte, length)
		err = binary.Read(buffer, binary.BigEndian, &retvar)
		return retvar, err
	case TLV_URI:
		err := binary.Read(buffer, binary.BigEndian, &length)
		if err != nil {
			return string([]byte{}), err
		}
		retvar := make([]byte, length)
		err = binary.Read(buffer, binary.BigEndian, &retvar)
		return string(retvar), err
	default:
		return nil, errors.New("Reading TLV error ,Unkown TLV type: " + string(tag))
	}
}

func MakeTLVs(a []interface{}) ([]TLV, error) {
	tlvs := []TLV{}
	for _, one := range a {
		tlv, err := MakeTLV(one)
		if err != nil {
			return nil, err
		}
		tlvs = append(tlvs, *tlv)
	}
	return tlvs, nil
}

func ReadTLVs(tlvs []TLV) ([]interface{}, error) {
	values := []interface{}{}
	for _, tlv := range tlvs {
		one, err := ReadTLV(&tlv)
		if err != nil {
			return values, err
		}
		values = append(values, one)
	}
	return values, nil
}

func CastTLV(value interface{}, valueType int32) interface{} {
	switch valueType {
	case TLV_FLOAT64:
		return float64(value.(float64))
	case TLV_FLOAT32:
		return float32(value.(float64))
	case TLV_INT8:
		return int8(value.(float64))
	case TLV_INT16:
		return int16(value.(float64))
	case TLV_INT32:
		return int32(value.(float64))
	case TLV_INT64:
		return int64(value.(float64))
	case TLV_UINT8:
		return uint8(value.(float64))
	case TLV_UINT16:
		return uint16(value.(float64))
	case TLV_UINT32:
		return uint32(value.(float64))
	case TLV_UINT64:
		return uint64(value.(float64))
	case TLV_BYTES:
		return []byte(value.(string))
	case TLV_URI:
		return value.(string)
	default:
		return nil
	}
}
