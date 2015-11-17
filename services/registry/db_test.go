package main

import (
	"testing"
)

func TestWrongDB(t *testing.T) {
	*confDBPass = "wrongpassword"

	_, err := getDB()
	if err == nil {
		t.Errorf("get db should fail with wrong db config")
	}
}
