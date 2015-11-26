package main

import (
	"bufio"
	"fmt"
	"github.com/PandoCloud/pando-cloud/pkg/server"
	"os"
	"strings"
)

func main() {
	// init server
	err := server.Init("pdcfg")
	if err != nil {
		fmt.Printf("pdcfg init error : %s\n", err)
		return
	}

	cmdHandler := NewCommandHander()
	cmdHandler.SetHandler("vendor", DoVendorCommand)
	cmdHandler.SetHandler("product", DoProductCommand)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}

		fragments := strings.Split(line, " ")
		if len(fragments) < 1 {
			fmt.Println("wrong command.")
			continue
		}
		cmd := fragments[0]
		handler, err := cmdHandler.GetHandler(cmd)
		if err != nil {
			fmt.Println(err)
			continue
		}

		args := fragments[1:]
		err = handler(args)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}
