package validation

import (
	"net/http"
)

func ReadyEndpoint(req *http.Request) http.Response {
	return http.Response{}
}
