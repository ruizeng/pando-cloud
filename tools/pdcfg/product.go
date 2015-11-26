package main

import (
	"errors"
	"fmt"
	"github.com/PandoCloud/pando-cloud/pkg/models"
	"github.com/PandoCloud/pando-cloud/pkg/server"
	"strconv"
)

func DoProductCommand(args []string) error {
	if len(args) < 1 {
		return errors.New("command arguments not enough!")
	}

	op := args[0]

	switch op {
	case "add":
		if len(args) != 3 {
			return errors.New("wrong command arguments, example:  product add vendorid productname")
		}
		vendorid, err := strconv.Atoi(args[1])
		if err != nil {
			return err
		}
		args := models.Product{
			VendorID:    int32(vendorid),
			ProductName: args[2],
		}
		reply := models.Product{}

		err = server.RPCCallByName("registry", "Registry.SaveProduct", &args, &reply)
		if err != nil {
			return err
		}

		fmt.Println("=======> product created successfully.")
		fmt.Printf("product id: %d\n", reply.ID)
		fmt.Printf("product key: %s\n", reply.ProductKey)
	default:
		return errors.New("operation not suported:" + op)
	}

	return nil
}
