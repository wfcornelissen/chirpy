package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	fileserver := http.FileServer(http.Dir("."))
	mux.Handle("/app/", http.StripPrefix("/app", fileserver))
	mux.Handle("/assets/logo.png", fileserver)
	mux.HandleFunc("/healthz", readyEndpoint)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Starting server on port 8080")
	server.ListenAndServe()

}

func readyEndpoint(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(http.StatusText(http.StatusOK)))
}
