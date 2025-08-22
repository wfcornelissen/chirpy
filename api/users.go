package api

import (
	"encoding/json"
	"net/http"

	"github.com/wfcornelissen/chirpy/internal/auth"
	"github.com/wfcornelissen/chirpy/internal/database"
	"github.com/wfcornelissen/chirpy/types"
)

func CreateUser(cfg *types.ApiConfig) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		user := types.UserRequest{
			Email:    "",
			Password: "",
		}

		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&user)
		if err != nil {
			RespondWithInternalServerError(res, "failed to decode during CreateUser")
		}

		if user.Email == "" {
			RespondWithBadRequest(res, "Email is required to create user in CreateUser.")
			return
		}
		if user.Password == "" {
			RespondWithBadRequest(res, "Password is required in CreateUser.")
			return
		}
		user.Password, err = auth.HashPassword(user.Password)
		if err != nil {
			RespondWithInternalServerError(res, "Failed to hash password during CreateUser.")
			return
		}

		params := database.CreateUserParams{
			Email:        user.Email,
			PasswordHash: user.Password,
		}
		dbUser, err := cfg.Dbquery.CreateUser(req.Context(), params)
		newUser := types.User{
			UserID:    dbUser.ID,
			CreatedAt: dbUser.CreatedAt.Time,
			UpdatedAt: dbUser.UpdatedAt.Time,
			EAddress:  dbUser.Email,
			Password:  "",
		}

		result, err := json.Marshal(newUser)
		if err != nil {
			RespondWithInternalServerError(res, "Internal Error: Couldn't Marshal JSON")
			return
		}
		res.WriteHeader(http.StatusCreated)
		res.Write(result)
	}
}

func UserLogin(cfg *types.ApiConfig) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		userRequest := types.UserRequest{
			Email:    "",
			Password: "",
		}

		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&userRequest)
		if err != nil {
			RespondWithInternalServerError(res, "failed to decode during UserLogin")
		}
		if userRequest.Email == "" {
			RespondWithBadRequest(res, "No email supplied.")
			return
		}

		dbUser, err := cfg.Dbquery.FindUser(req.Context(), userRequest.Email)
		if err != nil {
			RespondWithNotFound(res, "User not found in UserLogin.")
			return
		}

		err = auth.CompareHashAndPassword(userRequest.Password, dbUser.PasswordHash)
		if err != nil {
			RespondWithUnauthorised(res, "Incorrect email or password")
			return
		}
		tokenString, err := auth.MakeJWT(dbUser.ID, cfg.Secret)
		if err != nil {
			RespondWithInternalServerError(res, "Failed to create token in UserLogin")
		}
		refreshToken, err := auth.MakeRefreshToken(cfg, dbUser.ID, req)
		if err != nil {
			RespondWithInternalServerError(res, "Failed to create refresh token in UserLogin")
		}

		newUser := types.User{
			UserID:       dbUser.ID,
			CreatedAt:    dbUser.CreatedAt.Time,
			UpdatedAt:    dbUser.UpdatedAt.Time,
			EAddress:     dbUser.Email,
			Password:     "",
			AccessToken:  tokenString,
			RefreshToken: refreshToken,
		}

		result, err := json.Marshal(newUser)
		if err != nil {
			RespondWithInternalServerError(res, "Error Marshalling user in UserLogin")
			return
		}

		res.WriteHeader(http.StatusOK)
		res.Write(result)
	}
}
