package middleware

import (
	"net/http"
)

func HeadMiddleware(nextHandlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers","Content-Type")
		w.Header().Set("content-type", "application/json")
		nextHandlerFunc(w, r)
	}
}