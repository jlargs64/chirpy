package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	// TokenTypeAccess -
	TokenTypeAccess TokenType = "chirpy-access"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		Issuer:    string(TokenTypeAccess),
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedtoken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return signedtoken, nil
}

func ValidateJWT(token, tokenSecret string) (uuid.UUID, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	if claims, ok := parsedToken.Claims.(*jwt.RegisteredClaims); ok {
		userIDString, err := claims.GetSubject()
		if err != nil {
			return uuid.Nil, err
		}

		userID, err := uuid.Parse(userIDString)
		if err != nil {
			return uuid.Nil, err
		}
		return userID, nil
	} else {
		return uuid.Nil, errors.New("the claims could not be cast to *jwt.RegisteredClaims")
	}
}

func GetBearerToken(headers http.Header) (string, error) {
	bearerToken := headers.Get("Authorization")
	if len(bearerToken) == 0 || !strings.HasPrefix(bearerToken, "Bearer ") {
		return "", errors.New("no bearer token supplied in Authorization header")
	}
	if tokenSplit := strings.Split(bearerToken, " "); len(tokenSplit) == 2 {
		return tokenSplit[1], nil
	}
	return "", errors.New("no bearer token supplied in Authorization header")
}
