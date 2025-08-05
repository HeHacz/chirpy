package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func responseWithError(w http.ResponseWriter, code int, err error, msg string) {
	if err != nil {
		log.Printf("%v", err)
	}
	if code > 499 {
		log.Printf("Responding with status code 5xx.\nMessage: %s", msg)
	}
	type errorBody struct {
		Error string `json:"error"`
	}
	responseWithJSON(w, code, errorBody{
		Error: msg,
	})

}

func responseWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, fmt.Sprintf("Error marshaling JSON: %v", err))
		return
	}
	w.WriteHeader(code)
	w.Write(data)
}
