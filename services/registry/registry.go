package main

import (
	"flag"
	"github.com/PandoCloud/pando-cloud/pkg/generator"
	"github.com/PandoCloud/pando-cloud/pkg/models"
)

const (
	flagAESKey = "aeskey"
)

var confAESKey = flag.String(flagAESKey, "", "use your own aes encryting key.")

type Registry struct {
	keygen *generator.KeyGenerator
}

func NewRegistry() (*Registry, error) {
	gen, err := generator.NewKeyGenerator(*confAESKey)
	if err != nil {
		return nil, err
	}
	return &Registry{
		keygen: gen,
	}, nil
}

func (r *Registry) CreateVendor(vendor *models.Vendor, reply *models.Vendor) error {
	db, err := getDB()
	if err != nil {
		return err
	}

	err = db.Save(vendor).Error
	if err != nil {
		return err
	}

	key, err := r.keygen.GenRandomKey(int64(vendor.ID))
	if err != nil {
		return err
	}

	vendor.VendorKey = key

	err = db.Save(vendor).Error
	if err != nil {
		return err
	}

	reply = vendor
	return nil
}

func (r *Registry) ValidateProduct(key string, reply *models.Product) error {

	return nil
}
