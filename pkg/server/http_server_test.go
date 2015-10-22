package server

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

const (
	sayHi        = "hello pando"
	testHTTPHost = "localhost:12347"
)

type testHttpHandler struct{}

func (h testHttpHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(sayHi))
}

func validateHTTPServer(t *testing.T, url string) {
	response, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)

	if string(body) != sayHi {
		t.Fatalf("http server test error: want %s, got %s", sayHi, string(body))
	}
}

func validateHTTPSServer(t *testing.T, url string) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	response, err := client.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)

	if string(body) != sayHi {
		t.Errorf("https server test error: want %s, got %s", sayHi, string(body))
	}
}

func TestHTTPServer(t *testing.T) {
	initLog("test", "debug")

	svr := HTTPServer{
		addr:     testHTTPHost,
		handler:  testHttpHandler{},
		useHttps: false,
	}

	err := svr.Start()
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Millisecond * 100)

	validateHTTPServer(t, "http://"+testHTTPHost)
}
