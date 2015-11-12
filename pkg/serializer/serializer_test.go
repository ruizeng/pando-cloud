package serializer

import (
	"reflect"
	"testing"
)

type testStruct struct {
	Int1 int
	Str1 string
	Int2 int32
	Arr  []byte
}

func TestStringStructConvert(t *testing.T) {
	test := testStruct{0, "hello", 12, []byte{1, 0, 12}}
	str, err := Struct2String(test)
	if err != nil {
		t.Error(err)
	}
	stru := testStruct{}
	err = String2Struct(str, &stru)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(test, stru) {
		t.Errorf("wrong result %v, want %v", stru, test)
	}
}
