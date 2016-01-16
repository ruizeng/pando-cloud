package utils

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net/http"
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
	bReq := []byte(argReq)
	req, err := http.NewRequest(argType, argUrl, bytes.NewBuffer(bReq))
	if err != nil {
		return []byte{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	if argHead != nil {
		for key, vaule := range argHead {
			req.Header.Set(key, vaule)
		}
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}
