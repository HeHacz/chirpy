package main

import (
	"fmt"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, req *http.Request) {
	http.NotFound(w, req)
	return
	//fmt.Fprint(w,"404 Not Found")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
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
