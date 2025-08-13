package main

import (
	"net/http"

	"github.com/wfcornelissen/chirpy/api"
	"github.com/wfcornelissen/chirpy/types"
)

func MetricsReset(cfg *types.ApiConfig) http.HandlerFunc {
	if cfg.Platform != "dev" {
		return func(res http.ResponseWriter, req *http.Request) {
			api.RespondWithForbidden(res, "You do not have dev permissions")
		}
	}
	return func(res http.ResponseWriter, req *http.Request) {
		cfg.Dbquery.DeleteAllUsers(req.Context())
		res.WriteHeader(http.StatusOK)
	}
}
