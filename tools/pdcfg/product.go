package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/PandoCloud/pando-cloud/pkg/models"
	"github.com/PandoCloud/pando-cloud/pkg/productconfig"
	"github.com/PandoCloud/pando-cloud/pkg/server"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func addProduct() error {
	args := models.Product{}

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("vendor ID: ")
	id, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	vendor := strings.Replace(id, "\n", "", -1)
	vendorid, err := strconv.Atoi(vendor)
	if err != nil {
		return err
	}
	args.VendorID = int32(vendorid)

	fmt.Printf("product name: ")
	name, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	args.ProductName = strings.Replace(name, "\n", "", -1)

	fmt.Printf("product description: ")
	desc, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	args.ProductDescription = strings.Replace(desc, "\n", "", -1)

	fmt.Printf("product config json file: ")
	file, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	jsonfile := strings.Replace(file, "\n", "", -1)
	fi, err := os.Open(jsonfile)
	if err != nil {
		return err
	}
	content, err := ioutil.ReadAll(fi)
	config := string(content)
	fi.Close()
	_, err = productconfig.New(config)
	if err != nil {
		return err
	}
	args.ProductConfig = config

	reply := &models.Product{}

	err = server.RPCCallByName("registry", "Registry.SaveProduct", &args, reply)
	if err != nil {
		return err
	}

	fmt.Println("=======> product created successfully:")
	printStruct(reply)
	fmt.Println("=======")

	return nil
}

func DoProductCommand(args []string) error {
	if len(args) < 1 {
		return errors.New("command arguments not enough!")
	}

	op := strings.Replace(args[0], "\n", "", -1)

	switch op {
	case "add":
		if len(args) > 1 {
			return errors.New("unnecessary command arguments. just type 'product add'")
		}
		err := addProduct()
		if err != nil {
			return err
		}
	default:
		return errors.New("operation not suported:" + op)
	}

	return nil
}
