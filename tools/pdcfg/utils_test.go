package main

import (
	"testing"
)

type testStruct struct {
	aa int
	bb string
	cc float32
}

func TestPrintStruct(t *testing.T) {
	args := &testStruct{1, "2222", 1.23}
	printStruct(args)
}
