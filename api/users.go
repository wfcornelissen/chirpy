package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/wfcornelissen/chirpy/internal/auth"
	"github.com/wfcornelissen/chirpy/internal/database"
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

		decoder := json.NewDecoder(req.Body)
		tempUser := types.User{}
		err := decoder.Decode(&tempUser)
		if err != nil {
			RespondWithInternalServerError(res, "failed to decode during CreateUser")
		}

		if tempUser.EAddress == "" {
			user.EAddress = tempUser.EAddress
		}
		pass, hasPass := req.URL.User.Password()
		if !hasPass {
			RespondWithInternalServerError(res, "Failed to retrieve password during CreateUser.")
		}
		user.Password, err = auth.HashPassword(pass)
		if err != nil {
			RespondWithInternalServerError(res, "Failed to hash password during CreateUser.")
		}

		params := database.CreateUserParams{
			Email:        user.EAddress,
			PasswordHash: user.Password,
		}
		dbUser, err := cfg.Dbquery.CreateUser(req.Context(), params)
		user.UserID = dbUser.ID
		user.CreatedAt = dbUser.CreatedAt.Time
		user.UpdatedAt = dbUser.UpdatedAt.Time
		user.EAddress = dbUser.Email
		user.Password = ""

		NewUser, err := json.Marshal(user)
		if err != nil {
			RespondWithInternalServerError(res, "Internal Error: Couldn't Marshal JSON")
			return
		}
		res.WriteHeader(http.StatusCreated)
		res.Write(NewUser)
	}
}

func UserLogin(cfg *types.ApiConfig) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		user := types.User{
			UserID:    uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			EAddress:  "",
			Password:  "",
		}

		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&user)
		if err != nil {
			RespondWithInternalServerError(res, "failed to decode during UserLogin")
		}

		if user.EAddress == "" {
			RespondWithBadRequest(res, "No email supplied.")
		}

		dbUser, err := cfg.Dbquery.FindUser(req.Context(), user.EAddress)
		if err != nil {
			RespondWithNotFound(res, "User not found in UserLogin.")
		}

		err = auth.CompareHashAndPassword(user.Password, dbUser.PasswordHash)
		if err != nil {
			RespondWithUnauthorised(res, "Incorrect email or password")
		}

		dbUser.PasswordHash = ""
		result, err := json.Marshal(dbUser)
		if err != nil {
			RespondWithInternalServerError(res, "Error Marshalling user in UserLogin")
			return
		}

		res.WriteHeader(http.StatusOK)
		res.Write(result)
	}
}
