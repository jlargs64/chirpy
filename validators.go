package main

import (
	"encoding/json"
	"net/http"
)

type validationReq struct {
	Body string `json:"body"`
}
type validationSuccessRep struct {
	Valid bool `json:"valid"`
}

func HandlerValidateChirp(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	params := validationReq{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Bad JSON format", err)
		return
	}

	// Validate chirp
	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", err)
		return
	}

	successResp := validationSuccessRep{Valid: true}
	respondWithJSON(w, http.StatusOK, successResp)
}
