package handler

import (
	"encoding/json"
	"net/http"

	"github.com/carlosmartinez/challenge-bi/internal/apperror"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, err error) {
	var appErr *apperror.AppError
	switch e := err.(type) {
	case *apperror.AppError:
		appErr = e
	default:
		appErr = apperror.New("internal_error", "an unexpected error occurred", http.StatusInternalServerError)
	}
	writeJSON(w, appErr.HTTPStatus, map[string]string{
		"error":   appErr.Code,
		"message": appErr.Message,
	})
}
