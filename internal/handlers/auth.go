package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jlargs64/chirpy/internal/auth"
	"github.com/jlargs64/chirpy/internal/utils"
)

type loginReq struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds *int   `json:"expires_in_seconds,omitempty"`
}

type loginResp struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	Updatedat time.Time `json:"updated_at"`
	Token     string    `json:"token"`
}

const (
	notAuthMsg = "the email or password do not match"
)

func (config *APIConfig) HandleLogin(w http.ResponseWriter, req *http.Request) {
	// Parse req
	decoder := json.NewDecoder(req.Body)

	loginReq := &loginReq{}
	err := decoder.Decode(loginReq)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "login req is not in a valid format", err)
		return
	}

	user, err := config.DBQueries.GetUserByEmail(req.Context(), loginReq.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.RespondWithError(w, http.StatusUnauthorized, notAuthMsg, errors.New(notAuthMsg))
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "there was a problem getting the user with that email", err)
		}
		return
	}

	if ok, err := auth.CheckPasswordHash(loginReq.Password, user.HashedPassword); !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, notAuthMsg, err)
		return
	}

	// Check expiration

	var expirationDuration time.Duration
	if loginReq.ExpiresInSeconds == nil {
		expirationDuration = 3600 * time.Second
	} else {
		expirationDuration = time.Duration(*loginReq.ExpiresInSeconds) * time.Second
		if expirationDuration*time.Second > 3600*time.Second {
			expirationDuration = 3600 * time.Second
		}
	}
	token, err := auth.MakeJWT(user.ID, string(config.SigningKey), expirationDuration)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "the jwt could not be generated", err)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, &loginResp{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		Updatedat: user.UpdatedAt,
		Token:     token,
	})
}
