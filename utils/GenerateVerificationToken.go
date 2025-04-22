package utils

import (
	"crypto/rand"
	"encoding/hex"
)

 func GenerateVerificationToken() (string, error) {
	// Генерация 32 случайных байт (256 бит)
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(tokenBytes), nil
}
