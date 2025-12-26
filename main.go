package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	_ "github.com/lib/pq"

	"github.com/jlargs64/chirpy/internal/database"
	"github.com/jlargs64/chirpy/internal/handlers"
)

func main() {
	// Init vars
	const port = "8080"
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Could not access database:", err)
	}

	apiCfg := handlers.APIConfig{
		FileserverHits: atomic.Int32{},
		DBQueries:      database.New(db),
	}

	// Start server
	log.Println("Starting server...")

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Create app routes
	mux.Handle("/app/", http.StripPrefix("/app/",
		apiCfg.MiddlewareMetricsInc(http.FileServer(http.Dir(".")))))
	// Create API routes
	mux.HandleFunc("GET /api/healthz", handlers.HandlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlers.HandlerValidateChirp)

	// Create Admin routes
	mux.Handle("GET /admin/metrics", http.HandlerFunc(apiCfg.HandlerMetrics))
	mux.HandleFunc("POST /admin/reset", apiCfg.HandlerMetricsReset)

	log.Println("Serving on port:", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("server could not start: ", err)
	}
}
