package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/jlargs64/chirpy/internal/database"
	"github.com/jlargs64/chirpy/internal/handlers"
)

func main() {
	// Init env
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("could not read env file", err)
	}
	// Init vars
	const port = "8080"
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	signingKey := []byte(os.Getenv("SIGNING_KEY"))
	polkaAPIKey := os.Getenv("POLKA_API_KEY")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Could not access database:", err)
	}

	apiCfg := handlers.APIConfig{
		FileserverHits: atomic.Int32{},
		DBQueries:      database.New(db),
		Platform:       platform,
		SigningKey:     signingKey,
		PolkaAPIKey:    polkaAPIKey,
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

	// Users
	mux.HandleFunc("POST /api/users", apiCfg.HandleCreateUser)
	mux.HandleFunc("PUT /api/users", apiCfg.HandleChangeUser)
	mux.HandleFunc("POST /api/login", apiCfg.HandleLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.HandleRefreshToken)
	mux.HandleFunc("POST /api/revoke", apiCfg.HandleRefreshRevoke)

	// Webhooks
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.HandlePolkaWebhook)

	// Chirps
	mux.HandleFunc("GET /api/chirps", apiCfg.HandleGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.HandleGetChirpByID)
	mux.HandleFunc("POST /api/chirps", apiCfg.HandleCreateChrip)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.HandleDeleteChirps)

	// Create Admin routes
	mux.Handle("GET /admin/metrics", http.HandlerFunc(apiCfg.HandlerMetrics))
	mux.HandleFunc("POST /admin/reset", apiCfg.HandlerReset)

	log.Println("Serving on port:", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("server could not start: ", err)
	}
}
