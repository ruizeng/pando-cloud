package server

import (
	"testing"
)

func TestLog(t *testing.T) {
	Log = nil
	err := initLog("wrongtest", "wronglevel")
	if err == nil {
		t.Errorf("init log should return error when level is wrong.")
	}

	err = initLog("test", "error")
	if err != nil {
		t.Error(err)
	}

	Log.Error("test log.")
}
