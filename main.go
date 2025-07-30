package main

import (
	"fmt"
	"net/http"

	"github.com/wfcornelissen/chirpy/admin"
	"github.com/wfcornelissen/chirpy/types"
	"github.com/wfcornelissen/chirpy/api"
)

func main() {
	cfg := &types.ApiConfig{}
	mux := http.NewServeMux()
	fileserver := http.FileServer(http.Dir("."))
	mux.Handle("/app/", admin.MiddlewareMetricsInc(cfg, http.StripPrefix("/app", fileserver)))
	mux.Handle("/app/assets/logo.png", fileserver)
	mux.HandleFunc("GET /admin/healthz", admin.ReadyEndpoint)
	mux.HandleFunc("GET /admin/metrics", admin.HandlerMetrics(cfg))
	mux.HandleFunc("POST /admin/reset", admin.MetricsReset(cfg))
	mux.HandleFunc("POST /api/validate_chirp", api.ValidateChirp)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Starting server on port 8080")
	server.ListenAndServe()

}
