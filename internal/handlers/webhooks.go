package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jlargs64/chirpy/internal/auth"
	"github.com/jlargs64/chirpy/internal/utils"
)

type webhookEvent struct {
	Event string         `json:"event"`
	Data  map[string]any `json:"data"`
}

func (config *APIConfig) HandlePolkaWebhook(w http.ResponseWriter, req *http.Request) {
	// Check authorization
	apiKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "missing api key", err)
		return
	}
	if apiKey != config.PolkaAPIKey {
		utils.RespondWithError(w, http.StatusUnauthorized, "bad api key", err)
		return
	}
	// Parse webhook response
	decoder := json.NewDecoder(req.Body)
	var event webhookEvent
	err = decoder.Decode(&event)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "the webhook event schema was not in the expected format", err)
		return
	}

	// Handle event types
	eventType := event.Event
	switch eventType {
	case "user.upgraded":
		upgradeUserToChirpyRed(config, w, req, event.Data)
	default:
		w.WriteHeader(http.StatusNoContent)
	}
}

func upgradeUserToChirpyRed(config *APIConfig, w http.ResponseWriter, req *http.Request, data map[string]any) {
	userID, ok := data["user_id"].(string)
	if !ok {
		utils.RespondWithError(w, http.StatusBadRequest, "missing user id", errors.New("missing user id"))
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid user id", err)
		return
	}
	_, err = config.DBQueries.UpgradeUserToChirpyRed(req.Context(), userUUID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "user could not be upgraded in the database", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
