package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	missingAuthErr := errors.New("missing authorization header")
	if len(authHeader) == 0 || !strings.HasPrefix(authHeader, "ApiKey") {
		return "", missingAuthErr
	}

	apiKey := strings.Fields(authHeader)
	if len(apiKey) != 2 {
		return "", missingAuthErr
	}

	return apiKey[1], nil
}
