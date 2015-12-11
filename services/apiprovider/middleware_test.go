package main

import (
	"testing"
)

func checkPair(t *testing.T, domain string, identifier string, shoudpass bool) {
	err := checkAppDomain(domain, identifier)
	if shoudpass {
		if err != nil {
			t.Errorf("check domain should pass, but failed: domain: %v, identifier: %v, err: %v", domain, identifier, err)
		}
	} else {
		if err == nil {
			t.Errorf("check domain should fail, but passed: domain: %v, identifier: %v, err: %v", domain, identifier, err)
		}
	}
}

func TestCheckAppDomain(t *testing.T) {
	checkPair(t, "", "", false)
	checkPair(t, "*", "1-2-3333", true)
	checkPair(t, "vendor/1", "1-2-3333", true)
	checkPair(t, "product/2", "1-2-3333", true)
	checkPair(t, "product/2", "1-a-3333", false)
	checkPair(t, "product/10", "1-a-3333", true)
	checkPair(t, "vendor/11", "b-a-3333", true)
	checkPair(t, "fff/product/2", "1-a-3333", false)
	checkPair(t, "product/10", "1-a-3333-11111", false)
}
