package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jlargs64/chirpy/internal/utils"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type createUserParams struct {
	Email string `json:"email"`
}

func (config *APIConfig) HandleCreateUser(w http.ResponseWriter, req *http.Request) {
	// Parse req
	decoder := json.NewDecoder(req.Body)

	createParams := &createUserParams{}
	err := decoder.Decode(createParams)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "create user params is not in a valid format", err)
		return
	}

	dbUser, err := config.DBQueries.CreateUser(req.Context(), createParams.Email)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "the database encountered an error when creating a user", err)
		return
	}

	user := &User{
		ID:        dbUser.ID,
		Email:     dbUser.Email,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}
	utils.RespondWithJSON(w, http.StatusCreated, user)
}
