package main

import (
	"errors"
	"github.com/PandoCloud/pando-cloud/pkg/models"
	"github.com/PandoCloud/pando-cloud/pkg/server"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
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
