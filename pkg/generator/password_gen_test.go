package generator

import (
	"testing"
)

func TestPasswordGen(t *testing.T) {
	pass, err := GenRandomPassword()
	if err != nil {
		t.Error(err)
	}
	t.Log(pass)
}
