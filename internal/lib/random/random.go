package random

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

func GenerateRandomString(lenght int) (string, error) {
	b := make([]byte, lenght)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(append(b, byte(time.Now().UnixNano())))[:lenght], nil
}
