package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	// w.WriteHeader(http.StatusOK)
	html := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())

	w.Write([]byte(html))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	// w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	// w.WriteHeader(http.StatusMethodNotAllowed)
	cfg.fileserverHits.Store(0)
	// w.Write([]byte("OK\n"))

	platform := os.Getenv("PLATFORM")
	if platform != "dev" {
		w.WriteHeader(403)
	}

	err := cfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		log.Printf("Error deleting users: %s", err)
		w.WriteHeader(500)
		return
	}

	resp := struct {
		Message string `json:"message"`
	}{
		Message: "reset successful",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

}
