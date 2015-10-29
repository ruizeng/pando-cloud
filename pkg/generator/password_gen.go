package generator

import (
	"crypto/rand"
	"encoding/base64"
)

const (
	ranPasswordByteLength = 24
)

// gen random base64 encoded password
func GenRandomPassword() (string, error) {
	ranbuf := make([]byte, ranPasswordByteLength)
	_, err := rand.Read(ranbuf)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(ranbuf), nil
}
