package server

import (
	"errors"
	"testing"
)

func TestErrorf(t *testing.T) {
	err := errorf("err %s %d", "1", 2)
	if err.Error() != "err 1 2" {
		t.Errorf("err is %v ,want %v", err, errors.New("err 1 2"))
	}
}
