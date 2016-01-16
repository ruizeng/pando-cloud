package utils

import (
	"testing"
)

func TestSendHttpRequest(t *testing.T) {
	headers := make(map[string]string)
	headers["test"] = "test"

	res, err := SendHttpRequest("http://www.baidu.com", "", "GET", headers)

	if err != nil {
		t.Error(err)
	}

	t.Log(len(res))
}
