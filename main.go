package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	content, err := os.ReadFile("chirpyAdmin.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(string(content), cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

func main() {
	const port = "8080"

	apiCfg := &apiConfig{}
	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)
	mux.HandleFunc("GET /api/healthz", OkCheck)
	mux.HandleFunc("POST /api/validate_chirp", apiCfg.validateChirp)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Starting server on port %s", port)
	log.Fatal(srv.ListenAndServe())
}

func OkCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) validateChirp(w http.ResponseWriter, r *http.Request) {
	type chirpRequest struct {
		RequestBody      string `json:"body"`
		RequestCleanBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)

	requests := chirpRequest{}
	err := decoder.Decode(&requests)
	if err != nil {
		cfg.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if len(requests.RequestBody) > 140 {
		cfg.respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	if !cfg.validateClean(requests.RequestBody) {
		requests.RequestCleanBody = cfg.cleanBody(requests.RequestBody)
	} else {
		requests.RequestCleanBody = requests.RequestBody
	}

	cfg.respondWithJSON(w, http.StatusOK, map[string]string{"cleaned_body": requests.RequestCleanBody})
}

func (cfg *apiConfig) validateClean(body string) bool {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}

	for _, badWord := range badWords {
		if strings.Contains(body, badWord) {
			return false
		}
	}
	return true
}

func (cfg *apiConfig) cleanBody(body string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(body, " ")

	for i, word := range words {
		for _, badWord := range badWords {
			if strings.Contains(strings.ToLower(word), strings.ToLower(badWord)) {
				words[i] = "****"
			}
		}
	}

	return strings.Join(words, " ")
}

func (cfg *apiConfig) respondWithError(w http.ResponseWriter, status int, message string) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	w.Write([]byte(message))
}

func (cfg *apiConfig) respondWithJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
