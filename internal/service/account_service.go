package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/carlosmartinez/challenge-bi/internal/apperror"
	"github.com/carlosmartinez/challenge-bi/internal/domain"
	"github.com/carlosmartinez/challenge-bi/internal/ports"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AccountService struct {
	repo         ports.AccountRepository
	customerRepo ports.CustomerRepository
}

func NewAccountService(repo ports.AccountRepository, customerRepo ports.CustomerRepository) *AccountService {
	return &AccountService{repo: repo, customerRepo: customerRepo}
}

func (s *AccountService) Create(ctx context.Context, customerID uuid.UUID, accountType domain.AccountType) (*domain.Account, error) {
	if _, err := s.customerRepo.GetByID(ctx, customerID); err != nil {
		return nil, err
	}

	number, err := s.generateAccountNumber()
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	a := &domain.Account{
		ID:            uuid.New(),
		CustomerID:    customerID,
		AccountNumber: number,
		Type:          accountType,
		Currency:      "USD",
		Balance:       decimal.Zero,
		Status:        domain.StatusActive,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := s.repo.Create(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *AccountService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *AccountService) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.AccountStatus) error {
	if status != domain.StatusBlocked && status != domain.StatusClosed {
		return apperror.ErrInvalidStatus
	}

	return s.repo.UpdateStatus(ctx, id, status)
}

// generateAccountNumber produces a unique account number in the format NB<ms><random6>.
func (s *AccountService) generateAccountNumber() (string, error) {
	ms := time.Now().UnixMilli()
	rnd := rand.Intn(999999)
	return fmt.Sprintf("NB%d%06d", ms, rnd), nil
}
