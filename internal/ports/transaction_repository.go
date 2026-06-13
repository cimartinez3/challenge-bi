package ports

import (
	"context"
	"time"

	"github.com/carlosmartinez/challenge-bi/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TransactionRepository interface {
	Create(ctx context.Context, tx *domain.Transaction) error
	GetByReference(ctx context.Context, ref string) (*domain.Transaction, error)
	ListByAccount(ctx context.Context, accountID uuid.UUID, limit int, before *time.Time) ([]domain.Transaction, error)
	// ExecuteTransfer runs debit + credit + transaction pair inside a single DB transaction.
	// If any step fails the entire operation is rolled back — no partial state is possible.
	ExecuteTransfer(ctx context.Context, fromID, toID uuid.UUID, amount decimal.Decimal, out, in *domain.Transaction) (*domain.Account, *domain.Account, error)
}
