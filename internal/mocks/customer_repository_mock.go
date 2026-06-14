package mocks

import (
	"context"

	"github.com/carlosmartinez/challenge-bi/internal/domain"
	"github.com/google/uuid"
)

type CustomerRepositoryMock struct {
	CreateFn  func(ctx context.Context, c *domain.Customer) error
	GetByIDFn func(ctx context.Context, id uuid.UUID) (*domain.Customer, error)
}

func (m *CustomerRepositoryMock) Create(ctx context.Context, c *domain.Customer) error {
	return m.CreateFn(ctx, c)
}

func (m *CustomerRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*domain.Customer, error) {
	return m.GetByIDFn(ctx, id)
}
