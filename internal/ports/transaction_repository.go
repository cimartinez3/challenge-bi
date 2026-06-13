package ports

import (
	"context"
	"time"

	"github.com/carlosmartinez/challenge-bi/internal/domain"
	"github.com/google/uuid"
)

type TransactionRepository interface {
	Create(ctx context.Context, tx *domain.Transaction) error
	GetByReference(ctx context.Context, ref string) (*domain.Transaction, error)
	ListByAccount(ctx context.Context, accountID uuid.UUID, limit int, before *time.Time) ([]domain.Transaction, error)
	CreateTransferPair(ctx context.Context, out, in *domain.Transaction) error
}
