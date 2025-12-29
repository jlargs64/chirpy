package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jlargs64/chirpy/internal/auth"
	"github.com/jlargs64/chirpy/internal/database"
	"github.com/jlargs64/chirpy/internal/utils"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type userReqParams struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (config *APIConfig) HandleCreateUser(w http.ResponseWriter, req *http.Request) {
	// Parse req
	decoder := json.NewDecoder(req.Body)

	createParams := &userReqParams{}
	err := decoder.Decode(createParams)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "create user params is not in a valid format", err)
		return
	}

	hashedPassword, err := auth.HashPassword(createParams.Password)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "there was an error hashing the password", err)
		return
	}

	dbUser, err := config.DBQueries.CreateUser(req.Context(), database.CreateUserParams{
		Email:          createParams.Email,
		HashedPassword: hashedPassword,
	})
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

func (config *APIConfig) HandleChangeUser(w http.ResponseWriter, req *http.Request) {
	// Check the auth req
	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "bearer token is missing", err)
		return
	}
	userID, err := auth.ValidateJWT(bearerToken, string(config.SigningKey))
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "uh uh uh you didn't say the magic word", err)
		return
	}

	// Parse update req
	decoder := json.NewDecoder(req.Body)
	userParams := &userReqParams{}
	err = decoder.Decode(userParams)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "update user params is not in a valid format", err)
		return
	}

	// Update the user
	hashedPassword, err := auth.HashPassword(userParams.Password)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "error hashing the password", err)
		return
	}
	updatedUser, err := config.DBQueries.UpdateUserById(req.Context(), database.UpdateUserByIdParams{
		Email:          userParams.Email,
		HashedPassword: hashedPassword,
		ID:             userID,
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "error updating the user", err)
		return
	}

	user := &User{
		ID:        updatedUser.ID,
		Email:     updatedUser.Email,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
	}
	utils.RespondWithJSON(w, http.StatusOK, user)
}
