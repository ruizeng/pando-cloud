package main

import (
	"fmt"
)

const (
	TestBroker     = "https://localhost"
	TestProductKey = "aec003c9018b9a572ceb19720e589c375ead1a2b1fbd0d089d067128611754ff"
)

func main() {
	dev := NewDevice(TestBroker, TestProductKey, "ffe34e", "version")

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

}
