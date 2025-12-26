package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jlargs64/chirpy/internal/utils"
)

type validationReq struct {
	Body string `json:"body"`
}
type validationSuccessRep struct {
	CleanedBody string `json:"cleaned_body"`
}

func HandlerValidateChirp(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	params := validationReq{}
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

	successResp := validationSuccessRep{CleanedBody: cleanedChirp}
	utils.RespondWithJSON(w, http.StatusOK, successResp)
}
