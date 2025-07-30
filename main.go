package main

import (
	"fmt"
	"net/http"

	"github.com/wfcornelissen/chirpy/admin"
)

func main() {
	cfg := &admin.ApiConfig{}
	mux := http.NewServeMux()
	fileserver := http.FileServer(http.Dir("."))
	mux.Handle("/app/", admin.MiddlewareMetricsInc(cfg, http.StripPrefix("/app", fileserver)))
	mux.Handle("/app/assets/logo.png", fileserver)
	mux.HandleFunc("GET /admin/healthz", admin.ReadyEndpoint)
	mux.HandleFunc("GET /admin/metrics", admin.HandlerMetrics(cfg))
	mux.HandleFunc("POST /admin/reset", admin.MetricsReset(cfg))

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Starting server on port 8080")
	server.ListenAndServe()

}
