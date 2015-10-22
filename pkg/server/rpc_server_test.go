package server

import (
	"net/rpc"
	"testing"
	"time"
)

const (
	testRPCHost = "localhost:12346"
)

var testRPCArgs = &Args{100, 200}

type Args struct {
	A, B int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func validateRPCServer(t *testing.T, addr string, method string) {
	rpccli, err := rpc.Dial("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}

	var reply int

	err = rpccli.Call(method, testRPCArgs, &reply)
	if err != nil {
		t.Fatal(err)
	}

	if reply != testRPCArgs.A*testRPCArgs.B {
		t.Fatalf("rpc test faild, want %d, got %d", testRPCArgs.A*testRPCArgs.B, reply)
	}
}

func TestRPCServer(t *testing.T) {
	initLog("test", "debug")

	testRPC := new(Arith)

	err := rpc.Register(testRPC)
	if err != nil {
		t.Fatal(err)
	}

	handler := rpcHandler{}

	svr := &RPCServer{
		TCPServer{
			addr:    testRPCHost,
			handler: &handler,
			useTls:  false,
		},
	}

	err = svr.Start()
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Millisecond * 100)

	validateRPCServer(t, testRPCHost, "Arith.Multiply")
}
