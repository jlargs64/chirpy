package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() (string, error) {
	rawRefreshToken := make([]byte, 32)
	_, err := rand.Read(rawRefreshToken)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(rawRefreshToken), nil
}
