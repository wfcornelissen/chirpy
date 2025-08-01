package admin

import (
	"fmt"
	"net/http"

	"github.com/wfcornelissen/chirpy/types"
)

func HandlerMetrics(cfg *types.ApiConfig) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "text/html")
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(fmt.Sprintf(`
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.FileserverHits.Load())))
	}
}

func MetricsReset(cfg *types.ApiConfig) http.HandlerFunc {
	if cfg.Platform != "dev" {
		return func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusForbidden)
			res.Write([]byte("You do not have the correct permissions"))
		}
	}
	return func(res http.ResponseWriter, req *http.Request) {
		cfg.Dbquery.DeleteAllUsers(req.Context())
		res.WriteHeader(http.StatusOK)
	}
}

func MiddlewareMetricsInc(cfg *types.ApiConfig, next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		cfg.FileserverHits.Add(1)
		next.ServeHTTP(res, req)
	})
}
