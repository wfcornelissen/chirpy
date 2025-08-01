package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/wfcornelissen/chirpy/types"
)

func CreateChirp(res http.ResponseWriter, req *http.Request) {

}

func ValidateChirp(res http.ResponseWriter, req *http.Request) {
	chirp := types.Chirp{}

	//Decode
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&chirp)
	if err != nil {
		RespondWithInternalServerError(res, "Couldn't decode JSON")
		return
	}

	//Lencheck
	if len(chirp.Body) > 140 {
		RespondWithBadRequest(res, "Chirp too long")
		return
	}
	if len(chirp.Body) < 1 {
		RespondWithBadRequest(res, "No content")
		return
	}

	//The Profane
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

	//Response
	bod, err := json.Marshal(chirp)
	if err != nil {
		RespondWithInternalServerError(res, "Couldn't encode JSON")
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(bod)

}
