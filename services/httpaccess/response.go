package main

// common response fields
type Common struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// device register response data field
type DeviceRegisterData struct {
	DeviceId         int64  `json:"device_id"`
	DeviceSecret     string `json:"device_secret"`
	DeviceKey        string `json:"device_key"`
	DeviceIdentifier string `json:"device_identifier"`
}

// device register response
type DeviceRegisterResponse struct {
	Common
	Data DeviceRegisterData `json:"data"`
}

// device auth response data field
type DeviceAuthData struct {
	AccessToken string `json:"access_token"`
	AccessAddr  string `json:"access_addr"`
}

// device auth response
type DeviceAuthResponse struct {
	Common
	Data DeviceAuthData `json:"data"`
}
