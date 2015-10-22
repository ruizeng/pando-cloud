package server

import (
	"testing"
	"time"
)

type Arith2 Arith

func (t *Arith2) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

type testTimer struct{}

func (t *testTimer) DoTask() {
	Log.Info("timer task fires.")
}

func TestServer(t *testing.T) {
	err := Init("simpletest")
	if err != nil {
		t.Error(err)
	}

	*confHTTPHost = "localhost:59000"
	*confTCPHost = "localhost:59001"
	*confRPCHost = "localhost:59002"
	*confUseHttps = true
	*confUseTls = true
	*confCAFile = "testdata/cert.pem"
	*confKeyFile = "testdata/key.pem"

	err = Init("test")
	if err != nil {
		t.Fatalf("%s", err)
	}

	// test TCP
	testtcp := &testEchoHandler{}
	err = RegisterTCPHandler(testtcp)
	if err != nil {
		t.Errorf("RegisterTCPHandler : %s", err)
	}

	// test RPC
	testrpc := new(Arith2)
	err = RegisterRPCHandler(testrpc)
	if err != nil {
		t.Errorf("RegisterRPCService : %s", err)
	}

	// test HTTP
	testhttp := &testHttpHandler{}
	err = RegisterHTTPHandler(testhttp)
	if err != nil {
		t.Errorf("RegisterHTTPServer : %s", err)
	}

	// test timer
	timer := &testTimer{}
	err = RegisterTimerTask(timer)
	if err != nil {
		t.Errorf("RegisterTimerTask : %s", err)
	}

	go func() {
		err = Run()
		if err != nil {
			t.Errorf("Run Server error : %s", err)
		}
	}()

	time.Sleep(time.Millisecond * 100)

	validateHTTPSServer(t, "https://"+*confHTTPHost)
	validateTLSServer(t, *confTCPHost)
	validateRPCServer(t, *confRPCHost, "Arith2.Multiply")
}
