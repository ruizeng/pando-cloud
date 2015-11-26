package main

import (
	"errors"
	"fmt"
	"github.com/PandoCloud/pando-cloud/pkg/models"
	"github.com/PandoCloud/pando-cloud/pkg/server"
)

func DoVendorCommand(args []string) error {
	if len(args) < 1 {
		return errors.New("command arguments not enough!")
	}

	op := args[0]

	switch op {
	case "add":
		if len(args) != 2 {
			return errors.New("wrong command arguments, example:  vendor add pandocloud")
		}
		args := models.Vendor{
			VendorName: args[1],
		}
		reply := models.Vendor{}

		err := server.RPCCallByName("registry", "Registry.SaveVendor", &args, &reply)
		if err != nil {
			return err
		}

		fmt.Println("=======> vendor created successfully.")
		fmt.Printf("vendor id: %d\n", reply.ID)
		fmt.Printf("vendor key: %s\n", reply.VendorKey)
	default:
		return errors.New("operation not suported:" + op)
	}

	return nil
}
