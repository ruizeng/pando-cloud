package main

import (
	"strconv"
)

func genDeviceIdentifier(vendor int32, product int32, device string) string {
	return strconv.FormatInt(int64(vendor), 16) + "-" + strconv.FormatInt(int64(product), 16) + "-" + device
}
