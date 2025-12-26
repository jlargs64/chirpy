package handlers

import (
	"errors"
	"net/http"

	"github.com/jlargs64/chirpy/internal/utils"
)

func (config *APIConfig) HandlerReset(w http.ResponseWriter, req *http.Request) {
	if config.Platform != "dev" {
		utils.RespondWithError(w, http.StatusForbidden, "not allowed in prod", errors.New("not allowed in prod"))
		return
	}
	w.WriteHeader(http.StatusOK)
	// Reset metrics
	config.FileserverHits.Store(0)
	// Reset users
	err := config.DBQueries.ResetUsers(req.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusForbidden, "users could not be reset", err)
		return
	}

	// Write message back to admin
	_, err = w.Write([]byte(http.StatusText(http.StatusOK)))
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "could not write a response back to user", err)
		return
	}
}
