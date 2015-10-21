package server

import (
	"fmt"
)

// will print a log and return error
func errorf(format string, a ...interface{}) error {
	err := fmt.Errorf(format, a...)
	Log.Error(err)
	return err
}
