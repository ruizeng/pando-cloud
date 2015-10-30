package main

import (
	"github.com/PandoCloud/pando-cloud/pkg/models"
	"github.com/PandoCloud/pando-cloud/pkg/mysql"
	"testing"
)

func testVendor(t *testing.T, r *Registry) {
	vendor := &models.Vendor{
		VendorName:        "testvendor",
		VendorDescription: "this is a test vendor",
	}
	err := r.CreateVendor(vendor, vendor)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(vendor)
}

func TestVendor(t *testing.T) {
	mysql.MigrateDatabase(defaultDBHost, defaultDBPort, defaultDBName, defaultDBUser, "")

	*confAESKey = "ABCDEFGHIJKLMNOPABCDEFGHIJKLMNOP"

	r, err := NewRegistry()
	if err != nil {
		t.Fatal(err)
	}

	testVendor(t, r)
}
