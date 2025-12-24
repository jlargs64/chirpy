package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

func (config *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		config.fileServerHits.Add(1)
		next.ServeHTTP(w, req)
	})
}

func (config *apiConfig) handlerMetrics(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	msg := fmt.Sprintf("Hits: %v", config.fileServerHits.Load())

	_, err := w.Write([]byte(msg))
	if err != nil {
		log.Println("there was an error getting metrics:", err)
	}
}

func (config *apiConfig) handlerMetricsReset(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	config.fileServerHits.Store(0)
	_, err := w.Write([]byte(http.StatusText(http.StatusOK)))
	if err != nil {
		log.Println("there was an error resetting metrics:", err)
	}
}

func main() {
	const port = "8080"
	apiConfig := &apiConfig{}

	log.Println("Starting server...")

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	mux.Handle("/app/", http.StripPrefix("/app/",
		apiConfig.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.Handle("GET /api/metrics", http.HandlerFunc(apiConfig.handlerMetrics))
	mux.HandleFunc("POST /api/reset", apiConfig.handlerMetricsReset)

	log.Println("Serving on port:", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("server could not start: ", err)
	}
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(http.StatusText(http.StatusOK)))
	if err != nil {
		log.Println("handlerReadiness failed to write:", err)
	}
}
