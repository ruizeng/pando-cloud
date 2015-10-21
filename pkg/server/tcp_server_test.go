package server

import (
	"fmt"
	"net"
	"testing"
)

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

func TestTCPServer(t *testing.T) {
	initLog("test", "debug")

	testAddr := "localhost:12345"
	testData := "hello pando"

	h := testEchoHandler{}

	svr := &TCPServer{
		addr:    testAddr,
		handler: h,
		useTls:  false,
	}

	err := svr.Start()
	if err != nil {
		t.Fatal(err)
	}

	cli, err := net.Dial("tcp", testAddr)
	if err != nil {
		t.Fatal(err)
	}
	_, err = cli.Write([]byte(testData))
	if err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 1024)
	length, err := cli.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	gotData := string(buf[:length])
	if gotData != testData {
		t.Fatal("echo server test failed. want: %s, got: %s", testData, gotData)
	}

}
