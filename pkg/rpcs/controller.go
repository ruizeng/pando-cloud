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

type ArgsOnEvent struct {
	DeviceId  uint64
	TimeStamp uint64
	SubDevice uint16
	No        uint16
	Priority  uint16
	Params    []tlv.TLV
}
type ReplyOnEvent ReplyEmptyResult
