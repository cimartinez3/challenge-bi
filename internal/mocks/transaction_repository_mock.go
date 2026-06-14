package mocks

import (
	"context"
	"time"

	"github.com/carlosmartinez/challenge-bi/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TransactionRepositoryMock struct {
	CreateFn          func(ctx context.Context, tx *domain.Transaction) error
	GetByReferenceFn  func(ctx context.Context, ref string) (*domain.Transaction, error)
	ListByAccountFn   func(ctx context.Context, accountID uuid.UUID, limit int, before *time.Time) ([]domain.Transaction, error)
	ExecuteTransferFn func(ctx context.Context, fromID, toID uuid.UUID, amount decimal.Decimal, out, in *domain.Transaction) (*domain.Account, *domain.Account, error)
}

func (m *TransactionRepositoryMock) Create(ctx context.Context, tx *domain.Transaction) error {
	return m.CreateFn(ctx, tx)
}

func (m *TransactionRepositoryMock) GetByReference(ctx context.Context, ref string) (*domain.Transaction, error) {
	return m.GetByReferenceFn(ctx, ref)
}

func (m *TransactionRepositoryMock) ListByAccount(ctx context.Context, accountID uuid.UUID, limit int, before *time.Time) ([]domain.Transaction, error) {
	return m.ListByAccountFn(ctx, accountID, limit, before)
}

func (m *TransactionRepositoryMock) ExecuteTransfer(ctx context.Context, fromID, toID uuid.UUID, amount decimal.Decimal, out, in *domain.Transaction) (*domain.Account, *domain.Account, error) {
	return m.ExecuteTransferFn(ctx, fromID, toID, amount, out, in)
}
