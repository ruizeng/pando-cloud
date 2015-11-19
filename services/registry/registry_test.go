package main

import (
	"github.com/PandoCloud/pando-cloud/pkg/models"
	"github.com/PandoCloud/pando-cloud/pkg/mysql"
	"github.com/PandoCloud/pando-cloud/pkg/rpcs"
	"testing"
)

func testVendor(t *testing.T, r *Registry) {
	vendor := &models.Vendor{
		VendorName:        "testvendor",
		VendorDescription: "this is a test vendor",
	}
	err := r.SaveVendor(vendor, vendor)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(vendor)
}

var (
	testProductKey = ""
)

func testProduct(t *testing.T, r *Registry) {
	product := &models.Product{
		VendorID:           1,
		ProductName:        "test product.",
		ProductDescription: "this is a test product",
		ProductConfig:      "{}",
	}
	err := r.SaveProduct(product, product)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(product)

	product.ProductName = "test for update."
	err = r.SaveProduct(product, product)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(product)

	reply := &models.Product{}
	err = r.ValidateProduct(product.ProductKey, reply)
	if err != nil {
		t.Error(err)
	}
	t.Log(reply)

	testProductKey = product.ProductKey

	err = r.ValidateProduct("this is a wrong key , you know", reply)
	if err == nil {
		t.Error("wrong key should fail product key validation.")
	}
}

func testDevice(t *testing.T, r *Registry) {
	args := &rpcs.ArgsDeviceRegister{
		ProductKey:    testProductKey,
		DeviceCode:    "ffffaaeeccbb",
		DeviceVersion: "android-gprs-v1",
	}
	device := &models.Device{}
	err := r.RegisterDevice(args, device)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(device)

	founddev := &models.Device{}
	err = r.FindDeviceById(device.ID, founddev)
	if err != nil {
		t.Error(err)
	}
	if device.DeviceIdentifier != founddev.DeviceIdentifier {
		t.Errorf("FindDeviceById not match, want %v, got %v", device, founddev)
	}

	err = r.FindDeviceByIdentifier(device.DeviceIdentifier, founddev)
	if err != nil {
		t.Error(err)
	}
	if device.ID != founddev.ID {
		t.Errorf("FindDeviceByIdentifier not match, want %v, got %v", device, founddev)
	}

	device.DeviceDescription = "test change device info."

	args2 := &rpcs.ArgsDeviceUpdate{
		DeviceIdentifier:  device.DeviceIdentifier,
		DeviceName:        "testupdatename",
		DeviceDescription: "test change device info.",
	}
	err = r.UpdateDeviceInfo(args2, device)
	if err != nil {
		t.Error(err)
	}
	t.Log(device)
}

func TestRegistry(t *testing.T) {
	err := mysql.MigrateDatabase(defaultDBHost, defaultDBPort, defaultDBName, defaultDBUser, "")
	if err != nil {
		t.Fatal(err)
	}

	*confAESKey = "ABCDEFGHIJKLMNOPABCDEFGHIJKLMNOP"

	*confDBPass = ""
	r, err := NewRegistry()
	if err != nil {
		t.Fatal(err)
	}

	testVendor(t, r)
	testProduct(t, r)
	testDevice(t, r)
}
