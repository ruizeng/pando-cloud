package rpcs

import (
	"github.com/PandoCloud/pando-cloud/pkg/protocol"
)

type ArgsPutData struct {
	DeviceId  uint64
	Timestamp uint64
	Subdata   protocol.SubData
}
