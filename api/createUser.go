package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/wfcornelissen/chirpy/types"
)

func CreateUser(res http.ResponseWriter, req *http.Request) {
	user := types.User{
		UserID:    uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		EAddress:  "",
	}
	decoder := json.NewDecoder(req.Body)
	tempUser := types.User{}
	err := decoder.Decode(&tempUser)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte("Couldn Decode JSON"))
		return
	}

	if tempUser.EAddress != "" {
		user.EAddress = tempUser.EAddress
	}

	NewUser, err := json.Marshal(user)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte("Couldn't Marshal Response"))
		return
	}
	res.WriteHeader(http.StatusCreated)
	res.Write(NewUser)
}
