package helper

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func RandomString(length int) (string, error) {
	randomBytes := make([]byte, length)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("could not generate random bytes: %w", err)
	}

	return base64.URLEncoding.EncodeToString(randomBytes), nil
}
