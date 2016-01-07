package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/PandoCloud/pando-cloud/pkg/generator"
	"github.com/PandoCloud/pando-cloud/pkg/models"
	"github.com/PandoCloud/pando-cloud/pkg/rpcs"
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

// SaveVendor will create a vendor if the ID field is not initialized
// if ID field is initialized, it will update the conresponding vendor.
func (r *Registry) SaveVendor(vendor *models.Vendor, reply *models.Vendor) error {
	db, err := getDB()
	if err != nil {
		return err
	}

	if vendor.ID == 0 {
		// if ID field is not initialized, will generate key first
		err = db.Save(vendor).Error
		if err != nil {
			return err
		}

		key, err := r.keygen.GenRandomKey(int64(vendor.ID))
		if err != nil {
			return err
		}

		vendor.VendorKey = key
	}

	err = db.Save(vendor).Error
	if err != nil {
		return err
	}

	cache := getCache()
	cacheKey := fmt.Sprintf("Vendor:%v", vendor.ID)
	if _, ok := cache.Get(cacheKey); ok {
		cache.Delete(cacheKey)
	}

	reply.ID = vendor.ID
	reply.VendorName = vendor.VendorName
	reply.VendorDescription = vendor.VendorDescription
	reply.VendorKey = vendor.VendorKey
	reply.CreatedAt = vendor.CreatedAt
	reply.UpdatedAt = vendor.UpdatedAt

	return nil
}

// SaveProduct will create a product if the ID field is not initialized
// if ID field is initialized, it will update the conresponding product.
func (r *Registry) SaveProduct(product *models.Product, reply *models.Product) error {
	db, err := getDB()
	if err != nil {
		return err
	}

	if product.ID == 0 {
		// create product
		err = db.Save(product).Error
		if err != nil {
			return err
		}

		key, err := r.keygen.GenRandomKey(int64(product.ID))
		if err != nil {
			return err
		}

		product.ProductKey = key
	}

	err = db.Save(product).Error
	if err != nil {
		return err
	}

	cache := getCache()
	cacheKey := fmt.Sprintf("Product:%v", product.ID)
	if _, ok := cache.Get(cacheKey); ok {
		cache.Delete(cacheKey)
	}

	reply.ID = product.ID
	reply.ProductName = product.ProductName
	reply.ProductDescription = product.ProductDescription
	reply.ProductKey = product.ProductKey
	reply.ProductConfig = product.ProductConfig
	reply.CreatedAt = product.CreatedAt
	reply.UpdatedAt = product.UpdatedAt

	return nil
}

// SaveApplication will create a application if the ID field is not initialized
// if ID field is initialized, it will update the conresponding application.
func (r *Registry) SaveApplication(app *models.Application, reply *models.Application) error {
	db, err := getDB()
	if err != nil {
		return err
	}

	if app.ID == 0 {
		err = db.Save(app).Error
		if err != nil {
			return err
		}

		key, err := r.keygen.GenRandomKey(int64(app.ID))
		if err != nil {
			return err
		}

		app.AppKey = key
	}

	err = db.Save(app).Error
	if err != nil {
		return err
	}

	cache := getCache()
	cacheKey := fmt.Sprintf("Application:%v", app.ID)
	if _, ok := cache.Get(cacheKey); ok {
		cache.Delete(cacheKey)
	}

	reply.ID = app.ID
	reply.AppName = app.AppName
	reply.AppDescription = app.AppDescription
	reply.AppKey = app.AppKey
	reply.ReportUrl = app.ReportUrl
	reply.AppToken = app.AppToken
	reply.AppDomain = app.AppDomain
	reply.CreatedAt = app.CreatedAt
	reply.UpdatedAt = app.UpdatedAt

	return nil
}

// ValidateApplication try to validate the given app key.
// if success, it will reply the corresponding application
func (r *Registry) ValidateApplication(key string, reply *models.Application) error {
	db, err := getDB()
	if err != nil {
		return err
	}

	id, err := r.keygen.DecodeIdFromRandomKey(key)
	if err != nil {
		return err
	}

	cache := getCache()
	cacheKey := fmt.Sprintf("Application:%v", id)
	if cacheValue, ok := cache.Get(cacheKey); ok {
		reply = cacheValue.(*models.Application)
	} else {
		err = db.First(reply, id).Error
		if err != nil {
			return err
		}
		var storage models.Application
		storage = *reply
		cache.Set(cacheKey, &storage)
	}

	if reply.AppKey != key {
		return errors.New("app key not match.")
	}

	return nil
}

// FindVendor will find product by specified ID
func (r *Registry) FindVendor(id int32, reply *models.Vendor) error {
	db, err := getDB()
	if err != nil {
		return err
	}

	cache := getCache()
	cacheKey := fmt.Sprintf("Vendor:%v", id)
	if cacheValue, ok := cache.Get(cacheKey); ok {
		reply = cacheValue.(*models.Vendor)
	} else {
		err = db.First(reply, id).Error
		if err != nil {
			return err
		}
		var storage models.Vendor
		storage = *reply
		cache.Set(cacheKey, &storage)
	}

	return nil
}

// GetVendors will get all vendors in the platform.
func (r *Registry) GetVendors(noarg int, reply *[]models.Vendor) error {
	db, err := getDB()
	if err != nil {
		return err
	}

	return db.Find(reply).Error
}

// GetProducts will get all products in the platform.
func (r *Registry) GetProducts(noarg int, reply *[]models.Product) error {
	db, err := getDB()
	if err != nil {
		return err
	}

	return db.Find(reply).Error
}

// GetApplications will get all applications in the platform.
func (r *Registry) GetApplications(noarg int, reply *[]models.Application) error {
	db, err := getDB()
	if err != nil {
		return err
	}

	return db.Find(reply).Error
}

// FindProduct will find product by specified ID
func (r *Registry) FindProduct(id int32, reply *models.Product) error {
	db, err := getDB()
	if err != nil {
		return err
	}

	cache := getCache()
	cacheKey := fmt.Sprintf("Product:%v", id)
	if cacheValue, ok := cache.Get(cacheKey); ok {
		reply = cacheValue.(*models.Product)
	} else {
		err = db.First(reply, id).Error
		if err != nil {
			return err
		}
		var storage models.Product
		storage = *reply
		cache.Set(cacheKey, &storage)
	}

	return nil
}

// FindAppliation will find product by specified ID
func (r *Registry) FindApplication(id int32, reply *models.Application) error {
	db, err := getDB()
	if err != nil {
		return err
	}

	cache := getCache()
	cacheKey := fmt.Sprintf("Application:%v", id)
	if cacheValue, ok := cache.Get(cacheKey); ok {
		reply = cacheValue.(*models.Application)
	} else {
		err = db.First(reply, id).Error
		if err != nil {
			return err
		}
		var storage models.Application
		storage = *reply
		cache.Set(cacheKey, &storage)
	}

	return nil
}

// ValidProduct try to validate the given product key.
// if success, it will reply the corresponding product
func (r *Registry) ValidateProduct(key string, reply *models.Product) error {
	db, err := getDB()
	if err != nil {
		return err
	}

	id, err := r.keygen.DecodeIdFromRandomKey(key)
	if err != nil {
		return err
	}

	cache := getCache()
	cacheKey := fmt.Sprintf("Product:%v", id)
	if cacheValue, ok := cache.Get(cacheKey); ok {
		reply = cacheValue.(*models.Product)
	} else {
		err = db.First(reply, id).Error
		if err != nil {
			return err
		}
		var storage models.Product
		storage = *reply
		cache.Set(cacheKey, &storage)
	}

	if reply.ProductKey != key {
		return errors.New("product key not match.")
	}

	return nil
}

// RegisterDevice try to register a device to our platform.
// if the device has already been registered,
// the registration will success return the registered device before.
func (r *Registry) RegisterDevice(args *rpcs.ArgsDeviceRegister, reply *models.Device) error {
	db, err := getDB()
	if err != nil {
		return err
	}

	product := &models.Product{}
	err = r.ValidateProduct(args.ProductKey, product)
	if err != nil {
		return err
	}

	identifier := genDeviceIdentifier(product.VendorID, product.ID, args.DeviceCode)
	if db.Where(&models.Device{DeviceIdentifier: identifier}).First(reply).RecordNotFound() {
		// device is not registered yet.
		reply.ProductID = product.ID
		reply.DeviceIdentifier = identifier
		reply.DeviceName = product.ProductName // product name as default device name.
		reply.DeviceDescription = product.ProductDescription
		reply.DeviceVersion = args.DeviceVersion
		err = db.Save(reply).Error
		if err != nil {
			return err
		}
		// generate a random device key with hex encoding.
		reply.DeviceKey, err = r.keygen.GenRandomKey(reply.ID)
		if err != nil {
			return err
		}
		// generate a random password with base64 encoding.
		reply.DeviceSecret, err = generator.GenRandomPassword()
		if err != nil {
			return err
		}

		err = db.Save(reply).Error
		if err != nil {
			return err
		}
	} else {

		//delete cache
		cache := getCache()
		cacheKey := fmt.Sprintf("Device:%v", identifier)
		if _, ok := cache.Get(cacheKey); ok {
			cache.Delete(cacheKey)
		}
		
		// device has aleady been saved. just update version info.
		reply.DeviceVersion = args.DeviceVersion
		err = db.Save(reply).Error
		if err != nil {
			return err
		}

	}

	return nil
}

// FindDeviceByIdentifier will find the device by indentifier
func (r *Registry) FindDeviceByIdentifier(identifier string, reply *models.Device) error {
	db, err := getDB()
	if err != nil {
		return err
	}

	cache := getCache()
	cacheKey := fmt.Sprintf("Device:%v", identifier)
	if cacheValue, ok := cache.Get(identifier); ok {
		reply = cacheValue.(*models.Device)
	} else {
		err = db.Where(&models.Device{
			DeviceIdentifier: identifier,
		}).First(reply).Error
		if err != nil {
			return err
		}
		var storage models.Device
		storage = *reply
		cache.Set(cacheKey, &storage)
	}

	return nil
}

// FindDeviceById will find the device with given id
func (r *Registry) FindDeviceById(id int64, reply *models.Device) error {
	db, err := getDB()
	if err != nil {
		return err
	}

	err = db.Where(&models.Device{
		ID: id,
	}).First(reply).Error

	if err != nil {
		return err
	}
	return nil
}

// ValidateDevice will validate a device key and return the model if success.
func (r *Registry) ValidateDevice(key string, device *models.Device) error {
	id, err := r.keygen.DecodeIdFromRandomKey(key)
	if err != nil {
		return err
	}

	err = r.FindDeviceById(id, device)
	if err != nil {
		return err
	}

	if device.DeviceKey != key {
		return errors.New("device key not match.")
	}

	return nil
}

// UpdateDevice will update a device info by identifier
func (r *Registry) UpdateDeviceInfo(args *rpcs.ArgsDeviceUpdate, reply *models.Device) error {
	db, err := getDB()
	if err != nil {
		return err
	}

	err = r.FindDeviceByIdentifier(args.DeviceIdentifier, reply)
	if err != nil {
		return err
	}
	
	//delete cache
	cache := getCache()
	cacheKey := fmt.Sprintf("Device:%v", args.DeviceIdentifier)
	if _, ok := cache.Get(cacheKey); ok {
		cache.Delete(cacheKey)
	}

	reply.DeviceName = args.DeviceName
	reply.DeviceDescription = args.DeviceDescription

	err = db.Save(reply).Error
	if err != nil {
		return err
	}

	return nil
}
