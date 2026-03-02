package main

import (
	"fmt"
	"io"
	"net/http"
)

func handlerReset(cfg *apiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if cfg.platform != "dev" {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusForbidden)
			response := fmt.Sprintf("Forbidden: Reset is only allowed in dev environment. %s environment detected.", cfg.platform)
			io.WriteString(w, response)
			return
		}
		if err := cfg.dbQueries.DeleteUsers(req.Context()); err != nil {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, "Internal Server Error: Unable to reset users.")
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		cfg.fileserverHits.Store(0)
		io.WriteString(w, "Hits reset to 0.")
	}
}
