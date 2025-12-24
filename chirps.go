package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type validationErrResp struct {
	Error string `json:"error"`
}

type validationReq struct {
	Body string `json:"body"`
}
type validationSuccessRep struct {
	Valid bool `json:"valid"`
}

func sendValidationError(w http.ResponseWriter, statusCode int, errMsg string) {
	w.WriteHeader(statusCode)
	errBytes, err := json.Marshal(&validationErrResp{Error: errMsg})
	if err != nil {
		log.Println("could not create err resp:", err)
		return
	}
	_, err = w.Write(errBytes)
	if err != nil {
		log.Println("could not create err resp:", err)
		return
	}
}

func HandlerValidateChirp(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	params := validationReq{}
	err := decoder.Decode(&params)
	if err != nil {
		sendValidationError(w, http.StatusInternalServerError, "Body could not be read correctly")
		return
	}

	// Validate chirp
	if len(params.Body) > 140 {
		sendValidationError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	w.WriteHeader(http.StatusOK)
	successResp := validationSuccessRep{Valid: true}
	successBytes, err := json.Marshal(&successResp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("could not create success resp:", err)
		return
	}
	_, err = w.Write(successBytes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("could not create success resp:", err)
		return
	}
}
