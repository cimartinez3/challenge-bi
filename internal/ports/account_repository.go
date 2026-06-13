package ports

import (
	"context"

	"github.com/carlosmartinez/challenge-bi/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AccountRepository interface {
	Create(ctx context.Context, a *domain.Account) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Account, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.AccountStatus) error
	Debit(ctx context.Context, id uuid.UUID, amount decimal.Decimal) (*domain.Account, error)
	Credit(ctx context.Context, id uuid.UUID, amount decimal.Decimal) (*domain.Account, error)
}
