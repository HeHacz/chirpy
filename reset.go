package main

import (
	"io"
	"net/http"
)

func handlerReset(cfg *apiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		cfg.fileserverHits.Store(0)
		io.WriteString(w, "Hits reset to 0.")
	}
}
