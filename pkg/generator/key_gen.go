package generator

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
)

const (
	maxEncodeLen = 32
)

type KeyGenerator struct {
	AESKey string
}

func encryptAESCFB(msg, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(msg))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], msg)

	return ciphertext, nil
}

func decryptAESCFB(msg, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(msg) < aes.BlockSize {
		return nil, errors.New("decrypt message too short")
	}
	iv := msg[:aes.BlockSize]
	msg = msg[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(msg, msg)
	return msg, nil
}

func NewKeyGenerator(key string) (*KeyGenerator, error) {
	l := len(key)
	if l != 16 && l != 24 && l != 32 {
		return nil, errors.New("invalid aes key length, should be 16, 24 or 32 bytes.")
	}
	return &KeyGenerator{
		AESKey: key,
	}, nil
}

func (g *KeyGenerator) GenRandomKey(id int64) (string, error) {
	buf := make([]byte, maxEncodeLen-binary.Size(id)-aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		return "", nil
	}

	binid := bytes.NewBuffer([]byte{})
	binary.Write(binid, binary.BigEndian, id)

	buf = append(buf, binid.Bytes()...)

	binkey, err := encryptAESCFB(buf, []byte(g.AESKey))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(binkey), nil
}

// get id from encrypt strings
func (g *KeyGenerator) DecodeIdFromRandomKey(encrypted string) (int64, error) {
	buf, err := hex.DecodeString(encrypted)
	if err != nil {
		return 0, err
	}

	raw, err := decryptAESCFB(buf, []byte(g.AESKey))
	if err != nil {
		return 0, err
	}

	var id int64

	if len(raw) > maxEncodeLen || len(raw) < maxEncodeLen-aes.BlockSize-binary.Size(id) {
		return 0, errors.New("invalid key format.")
	}

	binbuf := bytes.NewBuffer(raw[maxEncodeLen-aes.BlockSize-binary.Size(id):])
	binary.Read(binbuf, binary.BigEndian, &id)

	return id, nil

}
