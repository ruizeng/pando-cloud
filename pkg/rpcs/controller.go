package rpcs

import (
	"github.com/PandoCloud/pando-cloud/pkg/protocol"
	"github.com/PandoCloud/pando-cloud/pkg/tlv"
)

type ArgsPutData struct {
	DeviceId  uint64
	Sequence  uint64
	Timestamp uint64
	Subdata   protocol.SubData
}
type ReplyPutData ReplyEmptyResult

type ArgsSetStatus struct {
	DeviceId uint64
	Subdata  protocol.SubData
}
type ReplySetStatus ReplyEmptyResult

type ArgsGetStatus ArgsDeviceId
type ReplyGetStatus struct {
	Subdata protocol.SubData
}

type ArgsOnEvent struct {
	DeviceId  uint64
	TimeStamp uint64
	SubDevice uint16
	No        uint16
	Priority  uint16
	Params    []tlv.TLV
}
type ReplyOnEvent ReplyEmptyResult

type ArgsSendCommand struct {
	DeviceId  uint64
	SubDevice uint16
	No        uint16
	Priority  uint16
	WaitTime  uint32
	Params    []tlv.TLV
}
type ReplySendCommand ReplyEmptyResult
