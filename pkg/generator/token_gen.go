package generator

import (
	"crypto/rand"
)

const (
	ranTokendByteLength = 16
)

// gen random token bytes
func GenRandomToken() ([]byte, error) {
	ranbuf := make([]byte, ranTokendByteLength)
	_, err := rand.Read(ranbuf)
	if err != nil {
		return nil, err
	}

	return ranbuf, nil
}
