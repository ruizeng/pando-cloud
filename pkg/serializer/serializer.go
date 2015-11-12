package serializer

import (
	"bytes"
	"encoding/gob"
)

// convert string to any kind of struct
func String2Struct(str string, target interface{}) error {
	bytes_buffer := bytes.NewBufferString(str)
	dec := gob.NewDecoder(bytes_buffer)
	err := dec.Decode(target)
	return err
}

// convert any kind of struct to string
func Struct2String(stru interface{}) (string, error) {
	var bytes_buffer bytes.Buffer
	enc := gob.NewEncoder(&bytes_buffer)
	err := enc.Encode(stru)
	return bytes_buffer.String(), err
}
