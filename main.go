package main

import (
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"time"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, req)
	})
}

func handlerMetrics(cfg *apiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		hits := cfg.fileserverHits.Load()
		io.WriteString(w, fmt.Sprintf("Hits: %d\n", hits))
	}
}

func handlerReset(cfg *apiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		cfg.fileserverHits.Swap(0)
	}
}

func handlerHealthCheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "OK")
}

func main() {
	const rootPath = "."
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}
	mux := http.NewServeMux()
	handler := http.FileServer(http.Dir(rootPath))
	mux.Handle("/app/", http.StripPrefix("/app/", apiCfg.middlewareMetricsInc(handler)))
	mux.HandleFunc("GET /healthz", handlerHealthCheck)
	mux.HandleFunc("GET /metrics", handlerMetrics(&apiCfg))
	mux.Handle("POST /reset", handlerReset(&apiCfg))
	server := http.Server{
		Addr:              "127.0.0.1:8080",
		Handler:           mux,
		ReadTimeout:       time.Second * 60,
		ReadHeaderTimeout: time.Second * 60,
		WriteTimeout:      time.Second * 60,
		IdleTimeout:       time.Second * 300,
		MaxHeaderBytes:    1 << 20,
	}
	fmt.Printf("Startiing server on address: %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("Server error: %v", err)
	}
}
