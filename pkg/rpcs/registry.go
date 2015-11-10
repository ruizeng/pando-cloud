package rpcs

// device register args
type DeviceRegisterArgs struct {
	ProductKey    string
	DeviceCode    string
	DeviceVersion string
}

// device update args
type DeviceUpdateArgs struct {
	DeviceIdentifier  string
	DeviceName        string
	DeviceDescription string
}
