package server

import (
	"testing"
)

func TestNetIf(t *testing.T) {
	readNetInterfaces()
	t.Logf("internal ip: %s", InternalIP)
	t.Logf("external ip: %s", ExternalIP)
}
