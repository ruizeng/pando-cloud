package main

import (
	"testing"
)

func checkPair(t *testing.T, domain string, identifier string, shoudpass bool) {
	err := checkAppDomain(domain, identifier)
	if shoudpass {
		if err != nil {
			t.Errorf("check domain should pass, but failed: domain: %v, identifier: %v", domain, identifier)
		}
	} else {
		if err == nil {
			t.Errorf("check domain should fail, but passed: domain: %v, identifier: %v", domain, identifier)
		}
	}
}

func TestCheckAppDomain(t *testing.T) {
	checkPair(t, "", "", false)
	checkPair(t, "/*", "1-2-3333", true)
}
