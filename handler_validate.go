package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/hehacz/chirpy/internal/database"
)

type parameters struct {
	Body   string `json:"body"`
	UserID string `json:"user_id"`
}

type chirpResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, req *http.Request) {
	const chirpLength = 140
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		responseWithError(w, http.StatusInternalServerError, err, "Couldn't decode parameters")
		return
	}
	if len(params.Body) > chirpLength {
		responseWithError(w, http.StatusBadRequest, nil, "Chirp is too long")
		return
	}
	userID, err := uuid.Parse(params.UserID)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, err, "Invalid user_id")
		return
	}

	createParams := database.CreateChirpParams{
		Body:   params.Body,
		UserID: userID,
	}

	resp, err := cfg.dbQueries.CreateChirp(req.Context(), createParams)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, err, "Couldn't create chirp")
		return
	}
	apiResp := chirpResponse{
		ID:        resp.ID,
		CreatedAt: resp.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: resp.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Body:      resp.Body,
		UserID:    resp.UserID,
	}
	responseWithJSON(w, http.StatusCreated, apiResp)
}
