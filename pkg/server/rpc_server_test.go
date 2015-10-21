package server

import (
	"net/rpc"
	"testing"
)

type Args struct {
	A, B int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func TestRPCServer(t *testing.T) {
	initLog("test", "debug")

	testAddr := "localhost:12346"
	testArgs := &Args{100, 200}
	testRPC := new(Arith)

	err := rpc.Register(testRPC)
	if err != nil {
		t.Fatal(err)
	}

	handler := rpcHandler{}

	svr := &RPCServer{
		TCPServer{
			addr:    testAddr,
			handler: &handler,
			useTls:  false,
		},
	}

	err = svr.Start()
	if err != nil {
		t.Fatal(err)
	}

	rpccli, err := rpc.Dial("tcp", testAddr)
	if err != nil {
		t.Fatal(err)
	}

	var reply int

	err = rpccli.Call("Arith.Multiply", testArgs, &reply)
	if err != nil {
		t.Fatal(err)
	}

	if reply != testArgs.A*testArgs.B {
		t.Fatalf("rpc test faild, want %d, got %d", testArgs.A*testArgs.B, reply)
	}

}
