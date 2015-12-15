package main

import (
	"flag"
	"fmt"
)

var (
	TestUrl        = flag.String("url", "https://localhost", "login url")
	TestProductKey = flag.String("productkey", "", "product key")
)

func main() {
	flag.Parse()

	if *TestProductKey == "" {
		fmt.Println("product key not provided. use -productkey flag")
		return
	}

	dev := NewDevice(*TestUrl, *TestProductKey, "ffe34e", "version")

	err := dev.DoRegister()
	if err != nil {
		fmt.Errorf("device register error %s", err)
		return
	}

	err = dev.DoLogin()
	if err != nil {
		fmt.Errorf("device login error %s", err)
		return
	}

	err = dev.DoAccess()
	if err != nil {
		fmt.Errorf("device access error %s", err)
		return
	}

}
