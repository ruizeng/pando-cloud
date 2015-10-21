// netif implements helper functions to read network interfaces.
// warning: ONLY suport standard Linux interface config.

package server

import (
	"errors"
	"net"
	"strings"
)

var (
	InternalIP string
	ExternalIP string
)

const (
	confInternalIP = "internal"
	confExternalIP = "external"
)

//to see if an IP is internal ip
func isInternalIP(ip string) bool {
	if ip == "127.0.0.1" {
		return true
	}

	ipSplit := strings.Split(ip, ".")

	if ipSplit[0] == "10" {
		return true
	}

	if ipSplit[0] == "172" && ipSplit[1] >= "16" && ipSplit[1] <= "31" {
		return true
	}

	if ipSplit[0] == "192" && ipSplit[1] == "168" {
		return true
	}

	return false
}

//read server IP
func readNetInterfaces() {
	interfaces, err := net.Interfaces()
	if err != nil {
		return
	}

	for _, inter := range interfaces {
		addr, err := inter.Addrs()
		if err != nil {
			continue
		}

		if !strings.Contains(inter.Name, "eth") {
			continue
		}

		if len(addr) == 0 {
			continue
		}

		ip := strings.Split(addr[0].String(), "/")[0]
		if isInternalIP(ip) {
			InternalIP = ip
		} else {
			ExternalIP = ip
		}
	}

	return
}

// fix host ip with "internal:port" or "external:port" format
func fixHostIp(addr string) (string, error) {
	if strings.Contains(addr, confInternalIP) {
		if InternalIP != "" {
			addr = strings.Replace(addr, confInternalIP, InternalIP, -1)
		} else {
			return addr, errors.New("server has no internal ip")
		}
	} else if strings.Contains(addr, confInternalIP) {
		if ExternalIP != "" {
			addr = strings.Replace(addr, confInternalIP, ExternalIP, -1)
		} else {
			return addr, errors.New("server has no external ip")
		}
	}
	return addr, nil
}
