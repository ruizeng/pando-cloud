package protocol

import (
	"github.com/PandoCloud/pando-cloud/pkg/tlv"
)

type CommandEventHead struct {
	Flag        uint8
	Timestamp   uint64
	Token       [16]byte
	SubDeviceid uint16
	No          uint16
	Priority    uint16
	ParamsCount uint16
}

type Command struct {
	Head   CommandEventHead
	Params []tlv.TLV
}

type Event struct {
	Head   CommandEventHead
	Params []tlv.TLV
}

type DataHead struct {
	Flag      uint8
	Timestamp uint64
	Token     [16]byte
}

type Data struct {
	Head    DataHead
	SubData []SubData
}

type SubDataHead struct {
	SubDeviceid uint16
	PropertyNum uint16
	ParamsCount uint16
}

type SubData struct {
	Head   SubDataHead
	Params []tlv.TLV
}
