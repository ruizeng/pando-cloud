package mqtt

import (
	"github.com/PandoCloud/pando-cloud/pkg/rpcs"
)

type Provider interface {
	ValidateDeviceToken(deviceid uint64, token []byte) error
	OnDeviceOnline(args rpcs.ArgsGetOnline)
	OnDeviceOffline(deviceid uint64)
	OnDeviceHeartBeat(deviceid uint64)
	OnDeviceMessage(deviceid uint64, msgtype string, message []byte)
}
