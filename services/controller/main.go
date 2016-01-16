package main

import (
	"github.com/PandoCloud/pando-cloud/pkg/server"
)

func main() {
	// init server
	err := server.Init("controller")
	if err != nil {
		server.Log.Fatal(err)
		return
	}

	// register a rpc service
	controller, err := NewController(*confMongoHost, *confRabbitHost)
	if err != nil {
		server.Log.Errorf("NewController Error: %s", err)
		return
	}

	err = server.RegisterRPCHandler(controller)
	if err != nil {
		server.Log.Errorf("Register RPC service Error: %s", err)
		return
	}

	// start to run
	err = server.Run()
	if err != nil {
		server.Log.Fatal(err)
	}
}
