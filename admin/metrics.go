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
	return func(res http.ResponseWriter, req *http.Request) {
		cfg.FileserverHits.Store(0)
		res.WriteHeader(http.StatusOK)
		res.Write([]byte("Hits reset to 0"))
	}
}

func MiddlewareMetricsInc(cfg *types.ApiConfig, next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		cfg.FileserverHits.Add(1)
		next.ServeHTTP(res, req)
	})
}
