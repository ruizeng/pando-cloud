package main

import (
	"testing"
)

func TestGenDeviceIdentifier(t *testing.T) {
	vendorid, productid, devicecode := int32(100), int32(100), "ffaf4fffeeaa"
	identifier := genDeviceIdentifier(vendorid, productid, devicecode)
	if identifier != "64-64-ffaf4fffeeaa" {
		t.Errorf("gen identifier error, need %s, got %s ", "64-64-ffaf4fffeeaa", identifier)
	}
}
