package generator

import (
	"testing"
)

func TestTokenGen(t *testing.T) {
	pass, err := GenRandomToken()
	if err != nil {
		t.Error(err)
	}
	t.Log(pass)
}
