package server

import (
	"crypto/tls"
	"fmt"
	"net"
	"testing"
	"time"
)

const (
	testTCPHost = "localhost:12345"
)

var testEchoData = "hello pando"

type testEchoHandler struct{}

func (h testEchoHandler) Handle(conn net.Conn) {
	buf := make([]byte, 1024)
	for {
		length, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
		}
		length, err = conn.Write(buf[:length])
		if err != nil {
			fmt.Println(err)
		}
	}
}

func validateTCPServer(t *testing.T, addr string) {
	cli, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}
	_, err = cli.Write([]byte(testEchoData))
	if err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 1024)
	length, err := cli.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	gotData := string(buf[:length])
	if gotData != testEchoData {
		t.Errorf("echo server test failed. want: %s, got: %s", testEchoData, gotData)
	}
}

func validateTLSServer(t *testing.T, addr string) {
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}
	cli, err := tls.Dial("tcp", addr, conf)
	if err != nil {
		t.Fatal(err)
	}
	_, err = cli.Write([]byte(testEchoData))
	if err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 1024)
	length, err := cli.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	gotData := string(buf[:length])
	if gotData != testEchoData {
		t.Errorf("echo server test failed. want: %s, got: %s", testEchoData, gotData)
	}
}

func TestTCPServer(t *testing.T) {
	initLog("test", "debug")

	h := testEchoHandler{}

	svr := &TCPServer{
		addr:    testTCPHost,
		handler: h,
		useTls:  false,
	}

	err := svr.Start()
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Millisecond * 100)

	validateTCPServer(t, testTCPHost)
}
