package api

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/wfcornelissen/chirpy/internal/auth"
	"github.com/wfcornelissen/chirpy/internal/database"
	"github.com/wfcornelissen/chirpy/types"
)

func CreateChirp(cfg *types.ApiConfig) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		userToken, err := auth.GetBearerToken(req.Header)
		if err != nil {
			RespondWithInternalServerError(res, "Couldn't retrieve user token")
			return
		}
		if userToken == "" {
			RespondWithUnauthorised(res, "No token provided")
			return
		}

		userID, err := auth.ValidateJWT(userToken, cfg.Secret)
		if err != nil {
			log.Printf("JWT validation error: %v", err)
			RespondWithUnauthorised(res, "Token could not be validated in CreateChirp")
			return
		}

		_, err = cfg.Dbquery.FindUserByUUID(req.Context(), userID)
		if err != nil {
			RespondWithNotFound(res, "User not found in CreateChirp")
			return
		}

		chirp := types.Chirp{}

		//Decode
		decoder := json.NewDecoder(req.Body)
		err = decoder.Decode(&chirp)
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

		// Set the user ID from the validated token
		chirp.UserID = userID

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
		optAuth := req.URL.Query().Get("author_id")
		chirps, err := cfg.Dbquery.GetAllChirps(req.Context())
		if err != nil {
			RespondWithInternalServerError(res, "Error retrieving from DB.")
		}
		results := make([]database.Chirp, 0)

		if optAuth != "" {
			for chirp := range chirps {
				if chirps[chirp].UserID.String() == optAuth {
					results = append(results, chirps[chirp])
				}
			}
		} else {
			results = chirps
		}
		optSort := req.URL.Query().Get("sort")
		switch optSort {
		case "desc":
			sort.Slice(results, func(i, j int) bool { return results[j].CreatedAt.Time.Before(results[i].CreatedAt.Time) })
		case "asc", "":
			sort.Slice(results, func(i, j int) bool { return results[i].CreatedAt.Time.Before(results[j].CreatedAt.Time) })
		}

		allChirps, err := json.Marshal(results)
		if err != nil {
			RespondWithInternalServerError(res, "Trouble encoding")
		}

		res.WriteHeader(http.StatusOK)
		res.Write(allChirps)

	}
}

func GetChirp(cfg *types.ApiConfig) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		chirpID := req.URL.Path[len("/api/chirps/"):]
		result, err := cfg.Dbquery.GetChirp(req.Context(), uuid.MustParse(chirpID))
		if err != nil {
			RespondWithNotFound(res, "Chirp not found")
			return
		}

		allChirps, err := json.Marshal(result)
		if err != nil {
			RespondWithInternalServerError(res, "Trouble encoding")
		}

		res.WriteHeader(http.StatusOK)
		res.Write(allChirps)

	}
}

func DeleteChirp(cfg *types.ApiConfig) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		token, err := auth.GetBearerToken(req.Header)
		if err != nil {
			RespondWithUnauthorised(res, "Access Token not found")
			return
		}

		userID, err := auth.ValidateJWT(token, cfg.Secret)
		if err != nil {
			RespondWithUnauthorised(res, "Could not validate JWT")
			return
		}

		user, err := cfg.Dbquery.FindUserByUUID(req.Context(), userID)
		if err != nil {
			RespondWithNotFound(res, "User not found")
			return
		}

		chirpID := req.URL.Path[len("/api/chirps/"):]
		chirp, err := cfg.Dbquery.GetChirp(req.Context(), uuid.MustParse(chirpID))
		if err != nil {
			RespondWithNotFound(res, "Chirp not found")
			return
		}
		if chirp.UserID != user.ID {
			RespondWithForbidden(res, "Not authorized to delete this chirp")
			return
		}
		err = cfg.Dbquery.DeleteChirpByID(req.Context(), chirp.ID)
		if err != nil {
			RespondWithInternalServerError(res, "Could not delete chirp")
			return
		}
		res.WriteHeader(http.StatusNoContent)
	}
}
