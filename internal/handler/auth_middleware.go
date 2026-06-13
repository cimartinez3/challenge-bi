package handler

import (
	"net/http"
)

// These are hardoceded keys for test.
// In production, use AWS service like Systems Manager Parameter Store or Secrets Manager.
var validAPIKeys = map[string]bool{
	"novobanco-dev-key-123": true,
}

func APIKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-API-Key")
		if key == "" {
			writeJSON(w, http.StatusUnauthorized, map[string]string{
				"error":   "missing_api_key",
				"message": "X-API-Key header is required",
			})
			return
		}

		if !validAPIKeys[key] {
			writeJSON(w, http.StatusForbidden, map[string]string{
				"error":   "invalid_api_key",
				"message": "the provided API key is not authorized",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
