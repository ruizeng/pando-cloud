package server

import (
	"testing"
)

func validateRPCClient(t *testing.T) {
	rpccli, err := NewRPCClient()
	if err != nil {
		t.Fatal(err)
	}

	args := &Args{100, 200}
	var reply int

	err = rpccli.Call("test", "Arith.Multiply", args, &reply)
	if err != nil {
		t.Fatal(err)
	}

	if reply != testRPCArgs.A*testRPCArgs.B {
		t.Fatalf("rpc client test faild, want %d, got %d", testRPCArgs.A*testRPCArgs.B, reply)
	}

	err = RPCCallByName("test", "Arith.Multiply", args, &reply)
	if err != nil {
		t.Fatal(err)
	}

	if reply != testRPCArgs.A*testRPCArgs.B {
		t.Fatalf("rpc client test faild, want %d, got %d", testRPCArgs.A*testRPCArgs.B, reply)
	}

	err = rpccli.Call("wrongtest", "Arith.Multiply", args, &reply)
	t.Log(err)
	if err == nil {
		t.Fatalf("rpc client should return error when no server found!")
	}
}
