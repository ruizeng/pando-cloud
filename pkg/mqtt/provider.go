package mqtt

type Provider interface {
	ValidateDeviceToken(deviceid uint64, token []byte) error
	OnDeviceOnline(deviceid uint64, interval uint16)
	OnDeviceOffline(deviceid uint64)
	OnDeviceHeartBeat(deviceid uint64)
	OnDeviceMessage(deviceid uint64, message []byte)
}
