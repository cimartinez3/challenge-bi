package ports

import (
	"context"

	"github.com/carlosmartinez/challenge-bi/internal/domain"
	"github.com/google/uuid"
)

type CustomerRepository interface {
	Create(ctx context.Context, c *domain.Customer) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Customer, error)
}
