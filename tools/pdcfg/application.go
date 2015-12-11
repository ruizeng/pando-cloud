package main

import (
	"errors"
	"fmt"
	"github.com/PandoCloud/pando-cloud/pkg/models"
	"github.com/PandoCloud/pando-cloud/pkg/server"
	"strings"
)

func DoApplicationCommand(args []string) error {
	if len(args) < 1 {
		return errors.New("command arguments not enough!")
	}

	op := args[0]

	switch op {
	case "add":
		if len(args) != 3 {
			return errors.New("wrong command arguments, example:  application add appname appdomain.")
		}
		args := models.Application{
			AppName:   args[1],
			AppDomain: strings.Replace(args[2], "\n", "", -1),
		}
		reply := models.Application{}

		err := server.RPCCallByName("registry", "Registry.SaveApplication", &args, &reply)
		if err != nil {
			return err
		}

		fmt.Println("=======> app created successfully.")
		fmt.Printf("app id: %d\n", reply.ID)
		fmt.Printf("app key: %s\n", reply.AppKey)
	default:
		return errors.New("operation not suported:" + op)
	}

	return nil
}
