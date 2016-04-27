package main

import (
	"github.com/go-martini/martini"
	// "github.com/martini-contrib/binding"
)

// martini router
func route(m *martini.ClassicMartini) {
	// find a device by key
	m.Get("/application/v1/device/info", GetDeviceInfoByKey)

	// find a device by identifier
	m.Get("/application/v1/devices/:identifier/info", ApplicationAuthOnDeviceIdentifer, GetDeviceInfoByIdentifier)

	// get devie current status
	m.Get("/application/v1/devices/:identifier/status/current",
		ApplicationAuthOnDeviceIdentifer, CheckDeviceOnline, CheckProductConfig,
		GetDeviceCurrentStatus)

	// get devie latest status
	m.Get("/application/v1/devices/:identifier/status/latest",
		ApplicationAuthOnDeviceIdentifer, CheckDeviceOnline, CheckProductConfig,
		GetDeviceLatestStatus)

	// set device status
	m.Put("/application/v1/devices/:identifier/status",
		ApplicationAuthOnDeviceIdentifer, CheckDeviceOnline, CheckProductConfig,
		SetDeviceStatus)

	// send a command to device
	m.Post("/application/v1/devices/:identifier/commands",
		ApplicationAuthOnDeviceIdentifer, CheckDeviceOnline, CheckProductConfig,
		SendCommandToDevice)

	// and a rule to device
	m.Post("/application/v1/devices/:identifier/rules",
		ApplicationAuthOnDeviceIdentifer, CheckDeviceIdentifier,
		AddRule)

}
