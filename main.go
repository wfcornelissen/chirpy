package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/wfcornelissen/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
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
	if cfg.platform != "dev" {
		cfg.respondWithError(w, http.StatusForbidden, "Reset endpoint is only available in development environment")
		return
	}

	// Reset hits counter
	cfg.fileserverHits.Store(0)

	// Reset users table
	err := cfg.db.ResetUsers(r.Context())
	if err != nil {
		cfg.respondWithError(w, http.StatusInternalServerError, "Failed to reset users")
		return
	}

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits and users reset to 0"))
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	type User struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	dbQueries := database.New(dbConn)

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       os.Getenv("PLATFORM"),
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)
	mux.HandleFunc("GET /api/healthz", OkCheck)
	mux.HandleFunc("POST /api/validate_chirp", apiCfg.validateChirp)
	mux.HandleFunc("POST /api/users", apiCfg.createUser)

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

// UserResponse is used to control the JSON response format
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func databaseUserToResponse(dbUser database.User) UserResponse {
	return UserResponse{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}
}

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	type UserRequest struct {
		Email string `json:"email"`
	}
	var userReq UserRequest

	err := json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		cfg.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	dbUser, err := cfg.db.CreateUser(r.Context(), userReq.Email)
	if err != nil {
		cfg.respondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	cfg.respondWithJSON(w, http.StatusCreated, databaseUserToResponse(dbUser))
}
