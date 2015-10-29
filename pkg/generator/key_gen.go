package generator

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strconv"
	"strings"
)

const (
	ranKeyByteLength = 20
	maxEncodeLen     = 32
)

type KeyGenerator struct {
	AESKey string
}

func encryptAESCFB(dst, src, key []byte) error {
	// Using IV same as key is probably bad
	iv := []byte(key)[:aes.BlockSize]
	aesBlockEncrypter, err := aes.NewCipher([]byte(key))
	if err != nil {
		return err
	}
	aesEncrypter := cipher.NewCFBEncrypter(aesBlockEncrypter, iv)
	aesEncrypter.XORKeyStream(dst, src)
	return nil
}

func decryptAESCFB(dst, src, key []byte) error {
	// Using IV same as key is probably bad
	iv := []byte(key)[:aes.BlockSize]
	aesBlockDecrypter, err := aes.NewCipher([]byte(key))
	if err != nil {
		return err
	}
	aesDecrypter := cipher.NewCFBDecrypter(aesBlockDecrypter, iv)
	aesDecrypter.XORKeyStream(dst, src)
	return nil
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
	id_str := strconv.FormatInt(id, 16)
	split_str := "#"
	ranbuf := make([]byte, maxEncodeLen-len(id_str)-len(split_str))
	_, err := rand.Read(ranbuf)
	if err != nil {
		return "", err
	}
	msg := string(ranbuf) + split_str + strconv.FormatInt(id, 16)
	encrypted := make([]byte, len(msg))
	err = encryptAESCFB(encrypted, []byte(msg), []byte(g.AESKey))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(encrypted), nil
}

// get id from encrypt strings
func (g *KeyGenerator) DecodeIdFromRandomKey(encrypted string) (int64, error) {
	decrypted := make([]byte, maxEncodeLen)
	byteArray, err := hex.DecodeString(encrypted)
	err = decryptAESCFB(decrypted, byteArray, []byte(g.AESKey))
	if err != nil {
		return 0, err
	}

	res := string(decrypted)
	split_index := strings.LastIndex(res, "#")
	if split_index == -1 {
		return 0, errors.New("invalid key format.")
	}
	device_id := res[split_index+1:]
	return strconv.ParseInt(string(device_id), 16, 64)

}
