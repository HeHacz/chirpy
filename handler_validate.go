package main

import (
	"encoding/json"
	"net/http"
)

type parameters struct {
	Body string `json:"body"`
}

func handlerValidate(w http.ResponseWriter, req *http.Request) {
	type validate struct {
		Cleaned_body string `json:"cleaned_body"`
	}
	const chirpLenght = 140
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		responseWithError(w, http.StatusInternalServerError, err, "Coudnt decode parameters")
		return
	}
	if len(params.Body) > chirpLenght {
		responseWithError(w, http.StatusBadRequest, nil, "Chirp is too long")
		return
	}
	responseWithJSON(w, http.StatusOK, validate{
		Cleaned_body: censorMessage(params.Body),
	})
}
