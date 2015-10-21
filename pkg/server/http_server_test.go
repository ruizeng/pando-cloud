package server

import (
	"io/ioutil"
	"net/http"
	"testing"
)

const (
	sayHi = "hello pando"
)

type HttpHandler struct{}

func (h HttpHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(sayHi))
}

func TestHTTPServer(t *testing.T) {
	initLog("test", "debug")

	testAddr := "localhost:12347"

	svr := HTTPServer{
		addr:     testAddr,
		handler:  HttpHandler{},
		useHttps: false,
	}

	err := svr.Start()
	if err != nil {
		t.Fatal(err)
	}

	response, err := http.Get("http://" + testAddr)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)

	if string(body) != sayHi {
		t.Fatalf("http server test error: want %s, got %s", sayHi, string(body))
	}
}
