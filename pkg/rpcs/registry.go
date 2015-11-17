package rpcs

// device register args
type ArgsDeviceRegister struct {
	ProductKey    string
	DeviceCode    string
	DeviceVersion string
}

// device update args
type ArgsDeviceUpdate struct {
	DeviceIdentifier  string
	DeviceName        string
	DeviceDescription string
}
