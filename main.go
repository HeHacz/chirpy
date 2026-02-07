package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/hehacz/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	dbQueries      *database.Queries
	fileserverHits atomic.Int32
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("cant open database")
	}
	const rootPath = "."
	apiCfg := apiConfig{
		dbQueries:      database.New(db),
		fileserverHits: atomic.Int32{},
	}
	mux := http.NewServeMux()
	handler := http.FileServer(http.Dir(rootPath))
	mux.Handle("/app/", http.StripPrefix("/app/", apiCfg.middlewareMetricsInc(handler)))
	mux.HandleFunc("GET /api/healthz", handlerHealthCheck)
	mux.HandleFunc("GET /admin/metrics", handlerMetrics(&apiCfg))
	mux.HandleFunc("POST /admin/reset", handlerReset(&apiCfg))
	mux.HandleFunc("POST /api/validate_chirp", handlerValidate)
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
