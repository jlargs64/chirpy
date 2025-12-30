package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jlargs64/chirpy/internal/auth"
	"github.com/jlargs64/chirpy/internal/database"
	"github.com/jlargs64/chirpy/internal/utils"
)

type createChirpRequest struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type chirpResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (config *APIConfig) HandleGetChirps(w http.ResponseWriter, req *http.Request) {
	chirps, err := config.DBQueries.GetChirps(req.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "there was an error getting all chirps from the database", err)
		return
	}

	chripResp := make([]chirpResponse, len(chirps))
	for i, chirp := range chirps {
		chripResp[i] = chirpResponse{
			ID:        chirp.ID,
			UserID:    chirp.UserID,
			Body:      chirp.Body,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
		}
	}
	utils.RespondWithJSON(w, http.StatusOK, chripResp)
}

func (config *APIConfig) HandleGetChirpByID(w http.ResponseWriter, req *http.Request) {
	chirpID := req.PathValue("chirpID")
	if len(chirpID) == 0 {
		utils.RespondWithError(
			w,
			http.StatusBadRequest,
			"no chirp id provided in req",
			errors.New("no chirp id provided in req"))
		return
	}

	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "chirp id is not a valid uuid", err)
		return
	}
	chirp, err := config.DBQueries.GetChirpById(req.Context(), chirpUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.RespondWithError(w, http.StatusNotFound, "chirp could not be found", err)
		} else {
			utils.RespondWithError(w, http.StatusBadRequest, "chirp could not be retrieved from the db", err)
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK,
		&chirpResponse{
			ID:        chirp.ID,
			UserID:    chirp.UserID,
			Body:      chirp.Body,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
		})
}

func (config *APIConfig) HandleCreateChrip(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	params := createChirpRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Bad JSON format", err)
		return
	}

	// Check user authorization
	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "user is unauthorized", err)
		return
	}
	userID, err := auth.ValidateJWT(bearerToken, string(config.SigningKey))
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "user is unauthorized", err)
		return
	}

	// Validate chirp
	if len(params.Body) > 140 {
		utils.RespondWithError(w, http.StatusBadRequest, "Chirp is too long", err)
		return
	}

	chirpSplit := strings.Split(params.Body, " ")
	for i, word := range chirpSplit {
		cleanedWord := strings.ToLower(word)
		if cleanedWord == "kerfuffle" || cleanedWord == "sharbert" || cleanedWord == "fornax" {
			chirpSplit[i] = "****"
		}
	}
	cleanedChirp := strings.Join(chirpSplit, " ")

	// Create the chirp
	chirpDBParams := database.CreateChirpParams{
		Body:   cleanedChirp,
		UserID: userID,
	}
	chirp, err := config.DBQueries.CreateChirp(req.Context(), chirpDBParams)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "could not create chirp in db", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, &chirpResponse{
		ID:        chirp.ID,
		UserID:    chirp.UserID,
		Body:      chirp.Body,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
	})
}

func (config *APIConfig) HandleDeleteChirps(w http.ResponseWriter, req *http.Request) {
	// Get chirp to delete id
	chirpID := req.PathValue("chirpID")
	if len(chirpID) == 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "missing chirp id", errors.New("missing chirp id"))
		return
	}
	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {

		utils.RespondWithError(w, http.StatusBadRequest, "bad chirp id provided", err)
		return
	}
	// Check user authorization
	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "user is unauthorized", err)
		return
	}
	userID, err := auth.ValidateJWT(bearerToken, string(config.SigningKey))
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "user is unauthorized", err)
		return
	}

	// Check if chirp exists
	_, err = config.DBQueries.GetChirpById(req.Context(), chirpUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.RespondWithError(w, http.StatusNotFound, "chirp not found", err)
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "db error could not get chirp", err)
		}
		return
	}

	rowsAffected, err := config.DBQueries.DeleteChirpById(req.Context(), database.DeleteChirpByIdParams{
		ID:     chirpUUID,
		UserID: userID,
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "unable to delete chirp", err)
		return
	}
	if rowsAffected == 0 {

		utils.RespondWithError(w, http.StatusForbidden, "chirp not found or not owned by user", errors.New("user does not own chirp"))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
