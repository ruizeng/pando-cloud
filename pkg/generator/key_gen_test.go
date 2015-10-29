package generator

import (
	"testing"
)

func TestKeyGen(t *testing.T) {
	generator, err := NewKeyGenerator("INVALIDKEY")
	if err == nil {
		t.Error("should return error when key length is invalid")
	}
	testid := int64(10000)
	generator, err = NewKeyGenerator("ABCDEFGHIJKLMNOPABCDEFGHIJKLMNOP")
	if err != nil {
		t.Fatal(err)
	}
	key, err := generator.GenRandomKey(testid)
	if err != nil {
		t.Error(err)
	}
	t.Log(key)
	id, err := generator.DecodeIdFromRandomKey(key)
	if err != nil {
		t.Error(err)
	}
	if id != testid {
		t.Errorf("wrong id %d, want %d", id, testid)
	}
}
