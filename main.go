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

	log.Println("Serving on port:", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("server could not start: ", err)
	}
}
