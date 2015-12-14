package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/PandoCloud/pando-cloud/pkg/models"
	"github.com/PandoCloud/pando-cloud/pkg/server"
	"os"
	"strings"
)

func addVendor() error {
	args := models.Vendor{}

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("vendor name: ")
	name, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	args.VendorName = strings.Replace(name, "\n", "", -1)

	fmt.Printf("vendor description: ")
	desc, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	args.VendorDescription = strings.Replace(desc, "\n", "", -1)

	reply := &models.Vendor{}

	err = server.RPCCallByName("registry", "Registry.SaveVendor", &args, reply)
	if err != nil {
		return err
	}

	fmt.Println("=======> vendor created successfully:")
	printStruct(reply)
	fmt.Println("=======")

	return nil
}

func DoVendorCommand(args []string) error {
	if len(args) < 1 {
		return errors.New("command arguments not enough!")
	}

	op := strings.Replace(args[0], "\n", "", -1)

	switch op {
	case "add":
		if len(args) > 1 {
			return errors.New("unnecessary command arguments. just type 'vendor add'")
		}
		err := addVendor()
		if err != nil {
			return err
		}
	default:
		return errors.New("operation not suported:" + op)
	}

	return nil
}
