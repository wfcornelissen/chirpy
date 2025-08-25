package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/wfcornelissen/chirpy/types"
)

func UpgradeToRed(cfg *types.ApiConfig) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		request := types.Webhook{}
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&request)
		if err != nil {
			RespondWithBadRequest(res, "Request does not conform to expectated format.")
			return
		}
		if request.Event != "user.upgraded" {
			res.WriteHeader(http.StatusNoContent)
			return
		}
		userIDString, err := uuid.Parse(request.Data["user_id"])
		if err != nil {
			RespondWithInternalServerError(res, "Failed to parse user ID in UpgradeToRed")
		}

		user, err := cfg.Dbquery.FindUserByUUID(req.Context(), userIDString)
		if err != nil {
			RespondWithNotFound(res, "User not found")
		}

		err = cfg.Dbquery.UpgradeToRed(req.Context(), user.ID)
		if err != nil {
			RespondWithInternalServerError(res, "Failed to upgrade")
		}

		res.WriteHeader(http.StatusNoContent)
	}
}
