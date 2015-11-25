package rpcs

import (
	"github.com/PandoCloud/pando-cloud/pkg/protocol"
	"github.com/PandoCloud/pando-cloud/pkg/tlv"
)

type ArgsSetStatus struct {
	DeviceId uint64
	Status   []protocol.SubData
}
type ReplySetStatus ReplyEmptyResult

type ArgsGetStatus ArgsDeviceId
type ReplyGetStatus struct {
	Status []protocol.SubData
}

type ArgsSendCommand struct {
	DeviceId  uint64
	SubDevice uint16
	No        uint16
	Priority  uint16
	WaitTime  uint32
	Params    []tlv.TLV
}
type ReplySendCommand ReplyEmptyResult
