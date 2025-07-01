package api

import (
	"golang.org/x/time/rate"
	"net/http"
)

func RateLimit(next http.Handler) http.Handler {
	limiter := rate.NewLimiter(1000, 100)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			sendJSONResponse(w, http.StatusTooManyRequests, "too many requests")
			return
		}
		next.ServeHTTP(w, r)
	})
}
