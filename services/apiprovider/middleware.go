package main

import (
	"errors"
	"github.com/PandoCloud/pando-cloud/pkg/models"
	"github.com/PandoCloud/pando-cloud/pkg/productconfig"
	"github.com/PandoCloud/pando-cloud/pkg/rpcs"
	"github.com/PandoCloud/pando-cloud/pkg/server"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"strconv"
	"strings"
)

func checkAppDomain(domain string, identifier string) error {
	domainPieces := strings.Split(domain, "/")
	identifierPieces := strings.Split(identifier, "-")
	if len(domainPieces) == 0 {
		return errors.New("wrong app domain format.")
	}
	if len(identifierPieces) != 3 {
		return errors.New("wrong identifier format.")
	}
	devvendorid, err := strconv.ParseUint(identifierPieces[0], 16, 64)
	if err != nil {
		return errors.New("wrong vendor format.")
	}
	devproductid, err := strconv.ParseUint(identifierPieces[1], 16, 64)
	if err != nil {
		return errors.New("wrong product format.")
	}

	if len(domainPieces) == 1 {
		if domainPieces[0] != "*" {
			return errors.New("wrong app domain " + domainPieces[0])
		}
		return nil
	}

	if len(domainPieces) == 2 {
		id, err := strconv.ParseUint(domainPieces[1], 10, 64)
		if err != nil {
			return errors.New("wrong app domain format..")
		}
		if domainPieces[0] == "vendor" {
			if id != devvendorid {
				return errors.New("app has no access right on device.")
			}
		} else if domainPieces[0] == "product" {
			if id != devproductid {
				return errors.New("app has no access right on device.")
			}
		} else {
			return errors.New("wrong app domain" + domain)
		}
	}

	if len(domainPieces) > 2 {
		return errors.New("wrong app domain" + domainPieces[0])
	}

	return nil
}

// check if app has access right on device of given identifier( in url params )
func ApplicationAuthOnDeviceIdentifer(context martini.Context, params martini.Params, req *http.Request, r render.Render) {
	identifier := params["identifier"]
	key := req.Header.Get("App-Key")

	if identifier == "" || key == "" {
		r.JSON(http.StatusOK, renderError(ErrDeviceNotFound, errors.New("missing device identifier or app key.")))
		return
	}

	app := &models.Application{}
	err := server.RPCCallByName("registry", "Registry.ValidateApplication", key, app)
	if err != nil {
		r.JSON(http.StatusOK, renderError(ErrAccessDenied, err))
		return
	}

	err = checkAppDomain(app.AppDomain, identifier)
	if err != nil {
		r.JSON(http.StatusOK, renderError(ErrAccessDenied, err))
		return
	}

}

// check if device is online.
func CheckDeviceOnline(context martini.Context, params martini.Params, req *http.Request, r render.Render) {
	identifier := params["identifier"]

	device := &models.Device{}
	err := server.RPCCallByName("registry", "Registry.FindDeviceByIdentifier", identifier, device)
	if err != nil {
		r.JSON(http.StatusOK, renderError(ErrDeviceNotFound, err))
		return
	}

	onlineargs := rpcs.ArgsGetDeviceOnlineStatus{
		Id: uint64(device.ID),
	}
	onlinereply := rpcs.ReplyGetDeviceOnlineStatus{}
	err = server.RPCCallByName("devicemanager", "DeviceManager.GetDeviceOnlineStatus", onlineargs, &onlinereply)
	if err != nil {
		server.Log.Errorf("get devie online status error: %v", err)
		r.JSON(http.StatusOK, renderError(ErrDeviceNotOnline, err))
		return
	}

	context.Map(device)
}

// get device identifier
func CheckDeviceIdentifier(context martini.Context, params martini.Params, req *http.Request, r render.Render) {
	identifier := params["identifier"]

	device := &models.Device{}
	err := server.RPCCallByName("registry", "Registry.FindDeviceByIdentifier", identifier, device)
	if err != nil {
		r.JSON(http.StatusOK, renderError(ErrDeviceNotFound, err))
		return
	}

	context.Map(device)
}

// check if proudct is ok and map a product config to context, must by called after CheckDevice
func CheckProductConfig(context martini.Context, device *models.Device,
	params martini.Params, req *http.Request, r render.Render) {
	product := &models.Product{}
	err := server.RPCCallByName("registry", "Registry.FindProduct", device.ProductID, product)
	if err != nil {
		r.JSON(http.StatusOK, renderError(ErrProductNotFound, err))
		return
	}

	c, err := productconfig.New(product.ProductConfig)
	if err != nil {
		r.JSON(http.StatusOK, renderError(ErrWrongProductConfig, err))
		return
	}

	context.Map(c)
}
