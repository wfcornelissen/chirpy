package api

import (
	"net/http"
	"encoding/json"
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

	chirp.Valid = true
	bod, err := json.Marshal(chirp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Couldn't encode JSON"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bod)
	return
}
