package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"

	log.Println("Starting server...")

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("."))))

	mux.HandleFunc("/healthz", handlerReadiness)
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
