package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

type apiConfig struct {
	fileserverHits atomic.Int32
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
	mux.HandleFunc("POST /reset", handlerReset(&apiCfg))
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
