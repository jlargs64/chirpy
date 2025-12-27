package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jlargs64/chirpy/internal/auth"
	"github.com/jlargs64/chirpy/internal/database"
	"github.com/jlargs64/chirpy/internal/utils"
)

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResp struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"created_at"`
	Updatedat    time.Time `json:"updated_at"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

type refreshResp struct {
	Token string `json:"token"`
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
	// Create refrsh token
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "the refresh token could not be generated", err)
		return
	}
	dbRefreshToken, err := config.DBQueries.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 60),
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "the refresh token could not be saved to the database", err)
		return
	}

	// Create access token
	expirationDuration := time.Hour * 1
	token, err := auth.MakeJWT(user.ID, string(config.SigningKey), expirationDuration)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "the jwt could not be generated", err)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, &loginResp{
		ID:           user.ID,
		Email:        user.Email,
		CreatedAt:    user.CreatedAt,
		Updatedat:    user.UpdatedAt,
		Token:        token,
		RefreshToken: dbRefreshToken.Token,
	})
}

func (config *APIConfig) HandleRefreshToken(w http.ResponseWriter, req *http.Request) {
	// Get bearer from header
	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "the token was not found or was expired", err)
		return
	}
	// Get refresh token from db
	res, err := config.DBQueries.GetUserFromRefreshToken(req.Context(), bearerToken)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "could not get refresh token from the db", err)
		return
	}

	// Check if refresh token is expired or revoked
	utcNow := time.Now().UTC()
	if res.ExpiresAt.UTC().Before(utcNow) || res.RevokedAt.Valid {
		utils.RespondWithError(w, http.StatusUnauthorized, "the token was not found or was expired", err)
		return
	}

	// Generate new access token
	accessToken, err := auth.MakeJWT(res.UserID, string(config.SigningKey), time.Hour*1)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "could not get refresh access token", err)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, &refreshResp{
		Token: accessToken,
	})
}

func (config *APIConfig) HandleRefreshRevoke(w http.ResponseWriter, req *http.Request) {
	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "the token was not found or was expired", err)
		return
	}

	_, err = config.DBQueries.RevokeRefreshToken(req.Context(), bearerToken)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "could not revoke token in the db", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
