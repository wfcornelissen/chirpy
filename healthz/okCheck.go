package healthz

import "net/http"

// OkCheck is a simple health check handler that returns HTTP 200 OK
func OkCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
