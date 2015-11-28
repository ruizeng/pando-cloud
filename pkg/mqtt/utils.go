package mqtt

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

func Uint16ToByte(value uint16) []byte {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, value)

	return buf.Bytes()
}

func ByteToUint16(buf []byte) uint16 {
	tmpBuf := bytes.NewBuffer(buf)
	var value uint16
	binary.Read(tmpBuf, binary.BigEndian, &value)

	return value
}

func boolToByte(val bool) byte {
	if val {
		return byte(1)
	}
	return byte(0)
}

func encodeLength(length int32, buf *bytes.Buffer) {
	if length == 0 {
		buf.WriteByte(0)
		return
	}

	for length > 0 {
		digit := length & 0x7f
		length = length >> 7
		if length > 0 {
			digit = digit | 0x80
		}
		buf.WriteByte(byte(digit))
	}
}

func decodeLength(r io.Reader) (int32, error) {
	var v int32
	var buf [1]byte
	var shift uint
	for i := 0; i < 4; i++ {
		if _, err := io.ReadFull(r, buf[:]); err != nil {
			return 0, err
		}

		b := buf[0]
		v |= int32(b&0x7f) << shift

		if b&0x80 == 0 {
			return v, nil
		}
		shift += 7
	}

	return 0, errors.New("length decode error")
}

func setUint8(val uint8, buf *bytes.Buffer) {
	buf.WriteByte(byte(val))
}

func setUint16(val uint16, buf *bytes.Buffer) {
	buf.WriteByte(byte(val & 0xff00 >> 8))
	buf.WriteByte(byte(val & 0x00ff))
}

func setString(val string, buf *bytes.Buffer) {
	length := uint16(len(val))
	setUint16(length, buf)
	buf.WriteString(val)
}

func getUint8(r io.Reader, packetRemaining *int32) (uint8, error) {
	if *packetRemaining < 1 {
		return 0, errors.New("dataExceedPacketError")
	}

	var b [1]byte
	if _, err := io.ReadFull(r, b[:]); err != nil {
		return 0, err
	}
	*packetRemaining--

	return b[0], nil
}

func getUint16(r io.Reader, packetRemaining *int32) (uint16, error) {
	if *packetRemaining < 2 {
		return 0, errors.New("dataExceedPacketError")
	}

	var b [2]byte
	if _, err := io.ReadFull(r, b[:]); err != nil {
		return 0, err
	}
	*packetRemaining -= 2

	return uint16(b[0])<<8 | uint16(b[1]), nil
}

func getString(r io.Reader, packetRemaining *int32) (string, error) {
	var retString string
	len, err := getUint16(r, packetRemaining)
	if err != nil {
		return retString, err
	}
	strLen := int(len)

	if int(*packetRemaining) < strLen {
		return retString, errors.New("dataExceedPacketError")
	}

	b := make([]byte, strLen)
	if _, err := io.ReadFull(r, b); err != nil {
		return retString, err
	}
	*packetRemaining -= int32(strLen)

	return string(b), nil
}
