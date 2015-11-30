package main

import (
	"github.com/PandoCloud/pando-cloud/pkg/server"
)

func main() {
	// init server
	err := server.Init("mqttaccess")
	if err != nil {
		server.Log.Fatal(err)
		return
	}

	a, err := NewAccess()
	if err != nil {
		server.Log.Fatal(err)
		return
	}

	// register a rpc service
	err = server.RegisterRPCHandler(a)
	if err != nil {
		server.Log.Errorf("Register RPC service Error: %s", err)
		return
	}

	// register a tcp service for mqtt
	err = server.RegisterTCPHandler(a.MqttHandler)
	if err != nil {
		server.Log.Errorf("Register TCP service Error: %s", err)
		return
	}

	// start to run
	err = server.Run()
	if err != nil {
		server.Log.Fatal(err)
	}
}
