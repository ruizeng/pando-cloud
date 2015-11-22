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
	m.Get("/application/v1/devices/:identifier/info", GetDeviceInfoByIdentifier)

	// get devie current status
	m.Get("/application/v1/devices/:identifier/status/current", GetDeviceCurrentStatus)

	// get devie latest status
	m.Get("/application/v1/devices/:identifier/status/latest", GetDeviceLatestStatus)

	// set devie status
	m.Put("/application/v1/devices/:identifier/status", SetDeviceStatus)

}
