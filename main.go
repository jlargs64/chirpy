package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const port = "8080"
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	log.Println("Starting server...")

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	mux.Handle("/app/", http.StripPrefix("/app/",
		apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
	// API
	mux.HandleFunc("GET /api/healthz", HandlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", HandlerValidateChirp)

	// Admin
	mux.Handle("GET /admin/metrics", http.HandlerFunc(apiCfg.handlerMetrics))
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerMetricsReset)

	log.Println("Serving on port:", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("server could not start: ", err)
	}
}
