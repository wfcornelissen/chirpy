package main

import (
	"log"
	"net/http"

	"github.com/wfcornelissen/chirpy/healthz"
)

func main() {
	const port = "8080"

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(".")))
	mux.Handle("/healthz", http.HandlerFunc(healthz.OkCheck))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Starting server on port %s", port)
	log.Fatal(srv.ListenAndServe())
}
