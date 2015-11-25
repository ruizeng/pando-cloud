package main

import (
	"github.com/PandoCloud/pando-cloud/pkg/server"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

func main() {
	// init server
	err := server.Init("apiprovidor")
	if err != nil {
		server.Log.Fatal(err)
		return
	}

	// martini setup
	martini.Env = martini.Prod
	handler := martini.Classic()
	handler.Use(render.Renderer())
	route(handler)

	// register a http handler
	err = server.RegisterHTTPHandler(handler)
	if err != nil {
		server.Log.Errorf("RegisterHTTPHandler Error: %s", err)
		return
	}

	// go
	err = server.Run()
	if err != nil {
		server.Log.Fatal(err)
	}
}
