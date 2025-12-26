package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jlargs64/chirpy/internal/database"
	"github.com/jlargs64/chirpy/internal/utils"
)

type createChirpRequest struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type createChirpResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (config *APIConfig) HandleCreateChrip(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	params := createChirpRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Bad JSON format", err)
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
		UserID: params.UserID,
	}
	chirp, err := config.DBQueries.CreateChirp(req.Context(), chirpDBParams)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "could not create chirp in db", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, &createChirpResponse{
		ID:        chirp.ID,
		UserID:    chirp.UserID,
		Body:      chirp.Body,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
	})
}
