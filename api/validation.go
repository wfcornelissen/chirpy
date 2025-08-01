package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/wfcornelissen/chirpy/types"
)

func ValidateChirp(w http.ResponseWriter, req *http.Request) {
	chirp := types.Chirp{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&chirp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Couldn't decode JSON"))
		return
	}

	if len(chirp.Body) > 140 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Chirp is too long"))
		return
	}

	if len(chirp.Body) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No Chirp Content"))
		return
	}

	chirp.Cleaned_body = ""
	words := strings.Split(chirp.Body, " ")
	for _, word := range words {
		switch strings.ToLower(word) {
		case "kerfuffle", "sharbert", "fornax":
			chirp.Cleaned_body += "**** "
		default:
			chirp.Cleaned_body += word + " "
		}
	}
	chirp.Cleaned_body = strings.TrimSpace(chirp.Cleaned_body)

	bod, err := json.Marshal(chirp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Couldn't encode JSON"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bod)

}
