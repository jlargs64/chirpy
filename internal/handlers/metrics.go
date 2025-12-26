// Package handlers contains all http handlers for chirpy
package handlers

import (
	"fmt"
	"log"
	"net/http"
)

func (config *APIConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		config.FileserverHits.Add(1)
		next.ServeHTTP(w, req)
	})
}

func (config *APIConfig) HandlerMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	html := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, config.FileserverHits.Load())

	_, err := w.Write([]byte(html))
	if err != nil {
		log.Println("there was an error getting metrics:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (config *APIConfig) HandlerMetricsReset(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	config.FileserverHits.Store(0)
	_, err := w.Write([]byte(http.StatusText(http.StatusOK)))
	if err != nil {
		log.Println("there was an error resetting metrics:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
