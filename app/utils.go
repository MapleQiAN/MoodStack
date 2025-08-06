package app

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateID generates a random ID
func GenerateID() (string, error) {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
