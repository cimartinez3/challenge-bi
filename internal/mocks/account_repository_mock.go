package mocks

import (
	"context"

	"github.com/carlosmartinez/challenge-bi/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AccountRepositoryMock struct {
	CreateFn       func(ctx context.Context, a *domain.Account) error
	GetByIDFn      func(ctx context.Context, id uuid.UUID) (*domain.Account, error)
	UpdateStatusFn func(ctx context.Context, id uuid.UUID, status domain.AccountStatus) error
	DebitFn        func(ctx context.Context, id uuid.UUID, amount decimal.Decimal) (*domain.Account, error)
	CreditFn       func(ctx context.Context, id uuid.UUID, amount decimal.Decimal) (*domain.Account, error)
}

func (m *AccountRepositoryMock) Create(ctx context.Context, a *domain.Account) error {
	return m.CreateFn(ctx, a)
}

func (m *AccountRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	return m.GetByIDFn(ctx, id)
}

func (m *AccountRepositoryMock) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.AccountStatus) error {
	return m.UpdateStatusFn(ctx, id, status)
}

func (m *AccountRepositoryMock) Debit(ctx context.Context, id uuid.UUID, amount decimal.Decimal) (*domain.Account, error) {
	return m.DebitFn(ctx, id, amount)
}

func (m *AccountRepositoryMock) Credit(ctx context.Context, id uuid.UUID, amount decimal.Decimal) (*domain.Account, error) {
	return m.CreditFn(ctx, id, amount)
}
