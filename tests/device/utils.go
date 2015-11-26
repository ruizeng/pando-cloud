package main

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
)

/**
  Params:
    argUrl: reqeust url
    argReq: reqeust contents
    argType: reqeust type
    argHead: reqeust head
  Retrun: reqesut result body
*/
func SendHttpRequest(argUrl string, argReq string, argType string, argHead map[string]string) ([]byte, error) {

	fmt.Print(">============<\n")
	fmt.Printf("[request url]:%v\n", argUrl)
	fmt.Printf("[request content]:%v\n", argReq)
	fmt.Printf("[request type]:%v\n", argType)
	fmt.Printf("[request head]:%+v\n", argHead)
	//*/
	bReq := []byte(argReq)
	req, err := http.NewRequest(argType, argUrl, bytes.NewBuffer(bReq))
	if err != nil {
		panic(err)
	}
	for key, vaule := range argHead {
		req.Header.Set(key, vaule)
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("[respons body]:%v\n", string(body))
	return body, nil
}

/**
assert the https response code
*/
func AssertHttpsCode(resp interface{}) {
	res := reflect.ValueOf(resp)
	// struct
	if res.Kind() == reflect.Struct {
		// exported field
		f := res.FieldByName("Code")
		if f.IsValid() {
			if f.Interface() == 0 {
				return
			} else {
				err := errors.New("[Https Response Error]code is not equal to 0!")
				panic(err)
			}
		}
	}
}
