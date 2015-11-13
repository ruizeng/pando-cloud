package token

import (
	"testing"
)

func TestTokenHelper(t *testing.T) {
	helper := NewHelper("localhost:6379")

	testid := uint64(123)

	token, err := helper.GenerateToken(testid)
	if err != nil {
		t.Error(err)
	}

	err = helper.ValidateToken(testid, token)
	if err != nil {
		t.Error(err)
	}

	err = helper.ClearToken(testid)
	if err != nil {
		t.Error(err)
	}

}
