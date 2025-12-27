package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestPasswordHash(t *testing.T) {
	passwordOne := "MyPassword123!"
	passwordTwo := "MySuperSecurePassword123!"
	hashOne, _ := HashPassword(passwordOne)
	hashTwo, _ := HashPassword(passwordTwo)

	tests := []struct {
		name          string
		password      string
		hash          string
		wantErr       bool
		matchPassword bool
	}{
		{
			name:          "Correct password",
			password:      passwordOne,
			hash:          hashOne,
			wantErr:       false,
			matchPassword: true,
		},
		{
			name:          "Password doesn't match different hash",
			password:      passwordOne,
			hash:          hashTwo,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Empty password",
			password:      "",
			hash:          hashOne,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Invalid hash",
			password:      passwordOne,
			hash:          "invalidhash",
			wantErr:       true,
			matchPassword: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && match != tt.matchPassword {
				t.Errorf("CheckPasswordHash() expects %v, got %v", tt.matchPassword, match)
			}
		})
	}
}

func TestJWTValidation(t *testing.T) {
	tokenSecret := "mysecrettoken"
	userID := uuid.New()
	token, _ := MakeJWT(userID, tokenSecret, time.Hour)

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			"Test JWT validation with a valid options",
			token,
			tokenSecret,
			userID,
			false,
		},
		{
			"Test JWT validation with a bad secret",
			token,
			"notgood",
			userID,
			true,
		},
		{
			"Test JWT validation with a bad token",
			"notavalidtoken",
			tokenSecret,
			userID,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("the token could not be validated when it should have: %v", err)
			}
			if tt.wantErr {
				return
			}

			if tt.wantUserID != gotUserID {
				t.Errorf("wanted %v from validated jwt but got %v", tt.wantUserID, gotUserID)
			}
		})
	}
}
