package server

import (
	"reflect"
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

func validateGetServerHosts(t *testing.T, flag string, want string) {
	hosts, err := serverInstance.svrmgr.GetServerHosts("test", flag)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(hosts, []string{want}) {
		t.Errorf("error get server hosts, want: %v, got %v", []string{want}, hosts)
	}
}

func validateServerManager(t *testing.T) {
	validateGetServerHosts(t, FlagTCPHost, *confTCPHost)
	validateGetServerHosts(t, FlagRPCHost, *confRPCHost)
	validateGetServerHosts(t, FlagHTTPHost, *confHTTPHost)
}

func registerBadHandlers(t *testing.T) {
	// test TCP
	testtcp := &testEchoHandler{}
	err := RegisterTCPHandler(testtcp)
	if err == nil {
		t.Errorf("RegisterTCPHandler shoud fail when server is not initialized.")
	}

	// test RPC
	testrpc := new(Arith2)
	err = RegisterRPCHandler(testrpc)
	if err == nil {
		t.Errorf("RegisterRPCService shoud fail when server is not initialized.")
	}

	// test HTTP
	testhttp := &testHttpHandler{}
	err = RegisterHTTPHandler(testhttp)
	if err == nil {
		t.Errorf("RegisterHTTPServer shoud fail when server is not initialized.")
	}

	// test timer
	timer := &testTimer{}
	err = RegisterTimerTask(timer)
	if err == nil {
		t.Errorf("RegisterTimerTask shoud fail when server is not initialized.")
	}
}

func registerHandlers(t *testing.T) {
	// test TCP
	testtcp := &testEchoHandler{}
	err := RegisterTCPHandler(testtcp)
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
}

func TestServer(t *testing.T) {
	*confHTTPHost = "localhost:59000"
	*confTCPHost = "localhost:59001"
	*confRPCHost = "localhost:59002"
	*confUseHttps = true
	*confUseTls = true
	*confCAFile = "testdata/cert.pem"
	*confKeyFile = "testdata/key.pem"
	*confEtcd = "http://localhost:2379"

	// before init , should all fail
	registerBadHandlers(t)

	err := Init("test")
	if err != nil {
		t.Fatalf("%s", err)
	}

	registerHandlers(t)

	go func() {
		err = Run()
		if err != nil {
			t.Errorf("Run Server error : %s", err)
		}
	}()

	time.Sleep(time.Second)

	validateHTTPSServer(t, "https://"+*confHTTPHost)
	validateTLSServer(t, *confTCPHost)
	validateRPCServer(t, *confRPCHost, "Arith2.Multiply")
	validateRPCClient(t)
	validateServerManager(t)
}
