package domain

import (
	"time"

	"github.com/carlosmartinez/challenge-bi/internal/apperror"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AccountStatus string
type AccountType string

const (
	StatusActive  AccountStatus = "active"
	StatusBlocked AccountStatus = "blocked"
	StatusClosed  AccountStatus = "closed"

	TypeSavings  AccountType = "savings"
	TypeChecking AccountType = "checking"
)

type Account struct {
	ID            uuid.UUID
	CustomerID    uuid.UUID
	AccountNumber string
	Type          AccountType
	Currency      string
	Balance       decimal.Decimal
	Status        AccountStatus
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (a *Account) IsOperable() error {
	switch a.Status {
	case StatusBlocked:
		return apperror.ErrAccountBlocked
	case StatusClosed:
		return apperror.ErrAccountClosed
	}
	return nil
}
