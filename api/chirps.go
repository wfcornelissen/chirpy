package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/wfcornelissen/chirpy/types"
)

func CreateChirp(res http.ResponseWriter, req *http.Request) {
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
	CleanedBody := ""
	words := strings.Split(chirp.Body, " ")
	for _, word := range words {
		switch strings.ToLower(word) {
		case "kerfuffle", "sharbert", "fornax":
			CleanedBody += "**** "
		default:
			CleanedBody += word + " "
		}
	}
	CleanedBody = strings.TrimSpace(CleanedBody)

	// 1/4
	chirp.Body = CleanedBody

	//Timestamps and UIDs
	// 2/4
	chirp.Id, err = uuid.NewUUID()
	if err != nil {
		RespondWithInternalServerError(res, "Couldn't create ID.")
	}
	time := time.Now()
	// 3/4
	chirp.CreatedAt = time
	// 4/4
	chirp.UpdatedAt = time

	result, err := json.Marshal(chirp)
	if err != nil {
		RespondWithInternalServerError(res, "Couldn't encode JSON")
		return
	}

	// DB logic
}
