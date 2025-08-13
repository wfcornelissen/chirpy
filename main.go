package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/wfcornelissen/chirpy/admin"
	"github.com/wfcornelissen/chirpy/api"
	"github.com/wfcornelissen/chirpy/internal/database"
	"github.com/wfcornelissen/chirpy/types"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("Error opening database: %v\n", err)
		return
	}
	defer db.Close()

	cfg := &types.ApiConfig{}
	cfg.Dbquery = database.New(db)
	cfg.Platform = os.Getenv("PLATFORM")
	mux := http.NewServeMux()
	fileserver := http.FileServer(http.Dir("."))
	mux.Handle("/app/", admin.MiddlewareMetricsInc(cfg, http.StripPrefix("/app", fileserver)))
	mux.Handle("/app/assets/logo.png", fileserver)
	mux.HandleFunc("GET /admin/healthz", admin.ReadyEndpoint)
	mux.HandleFunc("GET /admin/metrics", admin.HandlerMetrics(cfg))
	mux.HandleFunc("POST /admin/reset", admin.MetricsReset(cfg))
	mux.HandleFunc("POST /api/users", api.CreateUser)
	mux.HandleFunc("POST /api/chirps", api.CreateChirp)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Starting server on port 8080")
	server.ListenAndServe()

}
