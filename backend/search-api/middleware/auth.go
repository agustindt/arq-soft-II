package middleware

import (
	"net/http"
)

// RequireServiceToken is a middleware that validates service-to-service authentication
// For now, it's a simple pass-through, but it should be enhanced with actual token validation
func RequireServiceToken(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement proper service token validation
		next(w, r)
	}
}
