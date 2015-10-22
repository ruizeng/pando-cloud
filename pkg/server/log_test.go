package server

import (
	"testing"
)

func TestLog(t *testing.T) {
	err := initLog("test", "error")
	if err != nil {
		t.Error(err)
	}
	Log.Error("test log.")
}
