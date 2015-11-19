package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
)

// martini router
func route(m *martini.ClassicMartini) {
	// regist a device
	m.Post("/v1/devices/registration", binding.Json(DeviceRegisterArgs{}), RegisterDevice)

	// auth device
	m.Post("/v1/devices/authentication", binding.Json(DeviceAuthArgs{}), AuthDevice)

}
