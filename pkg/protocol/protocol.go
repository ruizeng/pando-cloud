package protocol

import (
	"github.com/PandoCloud/pando-cloud/pkg/tlv"
)

type CommandTypeHead struct {
	Flag        uint8
	Timestamp   uint64
	Token       [16]byte
	SubDeviceid uint16
	CommandId   uint16
	Priority    uint16
	ParamsCount uint16
}

type CommandType struct {
	Head   CommandTypeHead
	Params []tlv.TLV
}

type EventTypeHead struct {
	Flag        uint8
	Timestamp   uint64
	Token       [16]byte
	SubDeviceid uint16
	Eventid     uint16
	Priority    uint16
	ParamsCOunt uint16
}

type EventType struct {
	Head   EventTypeHead
	Params []tlv.TLV
}

type DataTypeHead struct {
	Flag      uint8
	Timestamp uint64
	Token     [16]byte
}

type DataType struct {
	Head    DataTypeHead
	SubData []SubDataType
}

type SubDataTypeHead struct {
	SubDeviceid uint16
	PropertyNum uint16
	ParamsCount uint16
}

type SubDataType struct {
	Head   SubDataTypeHead
	Params []tlv.TLV
}
