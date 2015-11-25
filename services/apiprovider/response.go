package main

// common response fields
type Common struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type DeviceInfoData struct {
	Identifier  string `json:"identifier"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

type DeviceInfoResponse struct {
	Common
	Data DeviceInfoData `json:"data"`
}

type DeviceStatusData map[string][]interface{}

type DeviceStatusResponse struct {
	Common
	Data DeviceStatusData `json:"data"`
}
