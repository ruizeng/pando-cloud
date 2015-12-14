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

func addApplication() error {
	args := models.Application{}

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("application name: ")
	name, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	args.AppName = strings.Replace(name, "\n", "", -1)

	fmt.Printf("application description: ")
	desc, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	args.AppDescription = strings.Replace(desc, "\n", "", -1)

	fmt.Printf("application domain: ")
	domain, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	args.AppDomain = strings.Replace(domain, "\n", "", -1)

	fmt.Printf("application report url: ")
	url, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	args.ReportUrl = strings.Replace(url, "\n", "", -1)

	fmt.Printf("application token: ")
	token, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	args.AppToken = strings.Replace(token, "\n", "", -1)

	reply := &models.Application{}

	err = server.RPCCallByName("registry", "Registry.SaveApplication", &args, reply)
	if err != nil {
		return err
	}

	fmt.Println("=======> application created successfully:")
	printStruct(reply)
	fmt.Println("=======")

	return nil
}

func DoApplicationCommand(args []string) error {
	if len(args) < 1 {
		return errors.New("command arguments not enough!")
	}

	op := strings.Replace(args[0], "\n", "", -1)

	switch op {
	case "add":
		if len(args) > 1 {
			return errors.New("unnecessary command arguments. just type 'application add'")
		}
		err := addApplication()
		if err != nil {
			return err
		}
	default:
		return errors.New("operation not suported:" + op)
	}

	return nil
}
