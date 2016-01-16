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
	err = r.FindVendor(vendor.ID, vendor)
	if err != nil {
		t.Fatal(err)
	}

	updateVendor := &models.Vendor{
		ID:                vendor.ID,
		VendorName:        "testvendorupdate",
		VendorDescription: "this is a test vendor",
	}

	err = r.SaveVendor(updateVendor, updateVendor)
	if err != nil {
		t.Fatal(err)
	}
	vendorRow := &models.Vendor{}
	err = r.FindVendor(updateVendor.ID, vendorRow)
	if err != nil {
		t.Fatal(err)
	}

	if vendorRow.VendorName != updateVendor.VendorName {
		t.Errorf("expect vendorName:%v, got:%v", updateVendor.VendorName, vendorRow.VendorName)
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

	productRow := &models.Product{}
	err = r.FindProduct(product.ID, productRow)
	if err != nil {
		t.Fatal(err)
	}

	if product.ProductKey != productRow.ProductKey {
		t.Fatalf("expected %v, got %v", product.ProductKey, productRow.ProductKey)
	}

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

	if reply.ProductName != product.ProductName {
		t.Errorf("expected %v, got %v", product.ProductName, reply.ProductName)
	}

	testProductKey = product.ProductKey

	err = r.ValidateProduct("this is a wrong key , you know", reply)
	if err == nil {
		t.Error("wrong key should fail product key validation.")
	}
}

func testApplication(t *testing.T, r *Registry) {
	app := &models.Application{
		AppToken:       "test-token",
		ReportUrl:      "http://localhost://6060",
		AppName:        "test",
		AppDescription: "this is a test app",
		AppDomain:      "/*",
	}
	err := r.SaveApplication(app, app)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(app)

	appRow := &models.Application{}
	err = r.FindApplication(app.ID, appRow)
	if err != nil {
		t.Fatal(err)
	}

	app.AppName = "another desc."
	err = r.SaveApplication(app, app)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(app)
	err = r.FindApplication(app.ID, appRow)
	if err != nil {
		t.Fatal(err)
	}

	if appRow.AppName != app.AppName {
		t.Errorf("expected %v, got %v", app.AppName, appRow.AppName)
	}

	reply := &models.Application{}
	err = r.ValidateApplication(app.AppKey, reply)
	if err != nil {
		t.Error(err)
	}
	t.Log(reply)

	err = r.ValidateApplication("this is a wrong key , you know", reply)
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
	
	devRow := &models.Device{}
	err = r.FindDeviceByIdentifier(device.DeviceIdentifier, devRow)
	if err != nil {
		t.Error(err)
	}
	if devRow.DeviceName != device.DeviceName {
		t.Errorf(" want %v, got %v", device.DeviceName, devRow.DeviceName)
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
	testApplication(t, r)
	testDevice(t, r)
}
