package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/wfcornelissen/chirpy/internal/database"
	"github.com/wfcornelissen/chirpy/types"
)

func CreateChirp(cfg *types.ApiConfig) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
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

		// Create DB entry
		dbChirp, err := cfg.Dbquery.CreateChirp(req.Context(), database.CreateChirpParams{Body: chirp.Body, UserID: chirp.UserID})
		if err != nil {
			RespondWithInternalServerError(res, "Error writing to database")
		}

		//Timestamps and UIDs
		// 2/4
		chirp.Id = dbChirp.ID
		// 3/4
		chirp.CreatedAt = dbChirp.CreatedAt.Time
		// 4/4
		chirp.UpdatedAt = dbChirp.UpdatedAt.Time

		result, err := json.Marshal(chirp)
		if err != nil {
			RespondWithInternalServerError(res, "Couldn't encode JSON")
			return
		}

		res.WriteHeader(http.StatusCreated)
		res.Write(result)

	}
}

func GetAllChirps(cfg *types.ApiConfig) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		result, err := cfg.Dbquery.GetAllChirps(req.Context())
		if err != nil {
			RespondWithInternalServerError(res, "Error retrieving from DB.")
		}

		allChirps, err := json.Marshal(result)
		if err != nil {
			RespondWithInternalServerError(res, "Trouble encoding")
		}

		res.WriteHeader(http.StatusOK)
		res.Write(allChirps)

	}
}

func GetChirp(cfg *types.ApiConfig) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		result, err := cfg.Dbquery.GetChirp(req.Context(), uuid.MustParse(req.URL.Query().Get("id")))
		if err != nil {
			RespondWithInternalServerError(res, "Error retrieving from DB.")
		}

		allChirps, err := json.Marshal(result)
		if err != nil {
			RespondWithInternalServerError(res, "Trouble encoding")
		}

		res.WriteHeader(http.StatusOK)
		res.Write(allChirps)

	}
}
