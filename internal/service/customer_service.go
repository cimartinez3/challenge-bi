package service

import (
	"context"
	"time"

	"github.com/carlosmartinez/challenge-bi/internal/domain"
	"github.com/carlosmartinez/challenge-bi/internal/ports"
	"github.com/google/uuid"
)

type CustomerService struct {
	repo ports.CustomerRepository
}

func NewCustomerService(repo ports.CustomerRepository) *CustomerService {
	return &CustomerService{repo: repo}
}

func (s *CustomerService) Create(ctx context.Context, fullName, email string) (*domain.Customer, error) {
	c := &domain.Customer{
		ID:        uuid.New(),
		FullName:  fullName,
		Email:     email,
		CreatedAt: time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}

	return c, nil
}

func (s *CustomerService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Customer, error) {
	return s.repo.GetByID(ctx, id)
}
