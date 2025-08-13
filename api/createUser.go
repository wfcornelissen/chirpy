package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/wfcornelissen/chirpy/internal/auth"
	"github.com/wfcornelissen/chirpy/types"
)

func CreateUser(cfg *types.ApiConfig) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		user := types.User{
			UserID:    uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			EAddress:  "",
			Password:  "",
		}
		pass, err := req.URL.User.Password()
		hashedPass, err := auth.HashPassword(pass)

		decoder := json.NewDecoder(req.Body)
		tempUser := types.User{}
		err := decoder.Decode(&tempUser)
		if err != nil {
			RespondWithInternalServerError(res, "Internal Error: Couldn't Decode JSON")
			return
		}

		if tempUser.EAddress != "" {
			user.EAddress = tempUser.EAddress
		}

		dbUser, err := cfg.Dbquery.CreateUser(req.Context(), user.EAddress)
		user.UserID = dbUser.ID
		user.CreatedAt = dbUser.CreatedAt.Time
		user.UpdatedAt = dbUser.UpdatedAt.Time
		user.EAddress = dbUser.Email

		NewUser, err := json.Marshal(user)
		if err != nil {
			RespondWithInternalServerError(res, "Internal Error: Couldn't Marshal JSON")
			return
		}
		res.WriteHeader(http.StatusCreated)
		res.Write(NewUser)
	}
}
