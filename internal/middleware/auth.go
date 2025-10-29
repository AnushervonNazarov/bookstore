package middleware

import (
	"net/http"
	"strings"
)

// Middleware function to handle authentication
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization token required", http.StatusUnauthorized)
			return
		}

		token := strings.Split(authHeader, "Bearer ")
		if len(token) != 2 {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		// Here you would typically validate the token.
		// For this example, we assume a simple token validation.

		// If token is valid, proceed with the next handler
		next.ServeHTTP(w, r)
	})
}
