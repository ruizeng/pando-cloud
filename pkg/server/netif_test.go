package server

import (
	"testing"
)

func TestNetIf(t *testing.T) {
	readNetInterfaces()
	t.Logf("internal ip: %s", InternalIP)
	t.Logf("external ip: %s", ExternalIP)
}

func TestIsInternalIP(t *testing.T) {
	testIPs := []string{"127.0.0.1", "192.168.5.234", "10.23.45.56", "172.17.2.4"}
	for _, ip := range testIPs {
		if isInternalIP(ip) == false {
			t.Errorf("test internal ip failed: %s", ip)
		}
	}
}

func TestFixHostIP(t *testing.T) {
	InternalIP = "10.1.1.1"
	ExternalIP = "5.1.1.1"
	fixedIP, err := fixHostIp("internal:40")
	if err != nil || fixedIP != "10.1.1.1:40" {
		t.Errorf("test fix host ip failed: %s, %s", fixedIP, err)
	}
	fixedIP, err = fixHostIp("external:40")
	if err != nil || fixedIP != "5.1.1.1:40" {
		t.Errorf("test fix host ip failed: %s, %s", fixedIP, err)
	}
}
