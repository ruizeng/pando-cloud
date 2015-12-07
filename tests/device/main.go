package main

import (
	"flag"
	"fmt"
)

var (
	TestUrl        = flag.String("url", "https://localhost", "login url")
	TestProductKey = flag.String("productkey", "aec003c9018b9a572ceb19720e589c375ead1a2b1fbd0d089d067128611754ff", "product key")
)

func main() {
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
