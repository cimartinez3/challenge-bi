package handler

import (
	"net/http"

	"github.com/carlosmartinez/challenge-bi/internal/apperror"
)

func errBadRequest(msg string) *apperror.AppError {
	return apperror.New("bad_request", msg, http.StatusBadRequest)
}
