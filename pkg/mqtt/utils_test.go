package mqtt

import (
	"bytes"
        "testing"
)

func TestUint16AndByte(t *testing.T) {
	value := uint16(1024)
	byteValue := Uint16ToByte(value)
	newValue := ByteToUint16(byteValue)

	if value != newValue {
		t.Errorf("uint16AndByte error, the origin:\n%x\n, now:\n%x\n", value, newValue)
	}
}

func TestUint8IO(t *testing.T) {
	rw := &TestReadWriter{
                w:   0,
                r:   0,
                buf: make([]byte, 10000),
        }

	buf := new(bytes.Buffer)
	value := uint8(100)
	setUint8(value, buf)

	_, err := rw.Write(buf.Bytes())
        if err != nil {
                t.Error(err)
        }

	packetRemaining := int32(rw.Size())
	newValue, err := getUint8(rw, &packetRemaining)
	if err != nil {
		t.Error(err)
	}

	if value != newValue { 
		t.Errorf("uint8IO error, the origin:\n%x\n, now:\n%x\n", value, newValue)
	}
}

func TestUint16IO(t *testing.T) {
	rw := &TestReadWriter{
                w:   0,
                r:   0,
                buf: make([]byte, 10000),
        }

	buf := new(bytes.Buffer)
	value := uint16(100)
	setUint16(value, buf)

	_, err := rw.Write(buf.Bytes())
        if err != nil {
                t.Error(err)
        }

	packetRemaining := int32(rw.Size())
	newValue, err := getUint16(rw, &packetRemaining)
	if err != nil {
		t.Error(err)
	}

	if value != newValue { 
		t.Errorf("uint16IO error, the origin:\n%x\n, now:\n%x\n", value, newValue)
	}
}

func TestStringIO(t *testing.T) {
	rw := &TestReadWriter{
                w:   0,
                r:   0,
                buf: make([]byte, 10000),
        }

	buf := new(bytes.Buffer)
	value := "itachili"
	setString(value, buf)

	_, err := rw.Write(buf.Bytes())
        if err != nil {
                t.Error(err)
        }

	packetRemaining := int32(rw.Size())
	newValue, err := getString(rw, &packetRemaining)
	if err != nil {
		t.Error(err)
	}

	if value != newValue { 
		t.Errorf("StringIO error, the origin:\n%x\n, now:\n%x\n", value, newValue)
	}
}

func TestLengthIO(t *testing.T) {
	rw := &TestReadWriter{
                w:   0,
                r:   0,
                buf: make([]byte, 10000),
        }

	buf := new(bytes.Buffer)
	value := int32(100)
	encodeLength(value, buf)

	_, err := rw.Write(buf.Bytes())
        if err != nil {
                t.Error(err)
        }

	newValue, err := decodeLength(rw)
	if err != nil {
		t.Error(err)
	}

	if value != newValue { 
		t.Errorf("LengthIO error, the origin:\n%x\n, now:\n%x\n", value, newValue)
	}
}
