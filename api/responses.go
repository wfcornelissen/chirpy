package api

import "net/http"

func RespondWithBadRequest(res http.ResponseWriter, message string) {
	res.WriteHeader(http.StatusBadRequest)
	res.Write([]byte(message))
}

func RespondWithForbidden(res http.ResponseWriter, message string) {
	res.WriteHeader(http.StatusBadRequest)
	res.Write([]byte(message))
}

func RespondWithInternalServerError(res http.ResponseWriter, message string) {
	res.WriteHeader(http.StatusBadRequest)
	res.Write([]byte(message))
}

func RespondWithNotFound(res http.ResponseWriter, message string) {
	res.WriteHeader(http.StatusNotFound)
	res.Write([]byte(message))
}
