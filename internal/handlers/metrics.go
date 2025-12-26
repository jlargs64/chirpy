// Package handlers contains all http handlers for chirpy
package handlers

import (
	"fmt"
	"net/http"

	"github.com/jlargs64/chirpy/internal/utils"
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
		utils.RespondWithError(w, http.StatusInternalServerError, "could not get metriccs", err)
		return
	}
}
