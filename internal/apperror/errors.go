package apperror

import "net/http"

type AppError struct {
	Code       string
	Message    string
	HTTPStatus int
}

func (e *AppError) Error() string {
	return e.Message
}

func New(code, message string, status int) *AppError {
	return &AppError{Code: code, Message: message, HTTPStatus: status}
}

var (
	ErrAccountNotFound   = New("account_not_found", "account not found", http.StatusNotFound)
	ErrAccountBlocked    = New("account_blocked", "account is blocked and cannot perform operations", http.StatusUnprocessableEntity)
	ErrAccountClosed     = New("account_closed", "account is closed and cannot perform operations", http.StatusUnprocessableEntity)
	ErrInsufficientFunds = New("insufficient_funds", "account balance is insufficient for this operation", http.StatusUnprocessableEntity)
	ErrInvalidAmount     = New("invalid_amount", "amount must be greater than zero", http.StatusBadRequest)
	ErrDuplicateRef      = New("duplicate_reference", "a transaction with this reference already exists", http.StatusConflict)
	ErrCustomerNotFound  = New("customer_not_found", "customer not found", http.StatusNotFound)
	ErrInvalidStatus     = New("invalid_status", "invalid account status transition", http.StatusBadRequest)
)
