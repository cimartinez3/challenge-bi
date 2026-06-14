package service

import (
	"context"
	"strings"
	"testing"

	"github.com/carlosmartinez/challenge-bi/internal/apperror"
	"github.com/carlosmartinez/challenge-bi/internal/domain"
	"github.com/carlosmartinez/challenge-bi/internal/mocks"
	"github.com/google/uuid"
)


func TestCreateAccount_OK(t *testing.T) {
	customerID := uuid.New()
	customer := &domain.Customer{ID: customerID}

	accountRepo := &mocks.AccountRepositoryMock{
		CreateFn: func(_ context.Context, a *domain.Account) error { return nil },
	}
	customerRepo := &mocks.CustomerRepositoryMock{
		GetByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Customer, error) { return customer, nil },
	}

	svc := NewAccountService(accountRepo, customerRepo)
	acc, err := svc.Create(context.Background(), customerID, domain.TypeSavings)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if acc.CustomerID != customerID {
		t.Errorf("expected customerID %v, got %v", customerID, acc.CustomerID)
	}
	if acc.Type != domain.TypeSavings {
		t.Errorf("expected type savings, got %v", acc.Type)
	}
	if acc.Status != domain.StatusActive {
		t.Errorf("expected status active, got %v", acc.Status)
	}
	if acc.Currency != "USD" {
		t.Errorf("expected currency USD, got %v", acc.Currency)
	}
	if !strings.HasPrefix(acc.AccountNumber, "NB") {
		t.Errorf("expected account number to start with NB, got %v", acc.AccountNumber)
	}
}

func TestCreateAccount_CustomerNotFound(t *testing.T) {
	accountRepo := &mocks.AccountRepositoryMock{}
	customerRepo := &mocks.CustomerRepositoryMock{
		GetByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Customer, error) {
			return nil, apperror.ErrCustomerNotFound
		},
	}

	svc := NewAccountService(accountRepo, customerRepo)
	acc, err := svc.Create(context.Background(), uuid.New(), domain.TypeChecking)

	if acc != nil {
		t.Error("expected nil account")
	}
	if err != apperror.ErrCustomerNotFound {
		t.Errorf("expected ErrCustomerNotFound, got %v", err)
	}
}

func TestUpdateStatus_Invalid(t *testing.T) {
	svc := NewAccountService(&mocks.AccountRepositoryMock{}, &mocks.CustomerRepositoryMock{})
	err := svc.UpdateStatus(context.Background(), uuid.New(), domain.StatusActive)

	if err != apperror.ErrInvalidStatus {
		t.Errorf("expected ErrInvalidStatus, got %v", err)
	}
}

func TestUpdateStatus_OK(t *testing.T) {
	accountRepo := &mocks.AccountRepositoryMock{
		UpdateStatusFn: func(_ context.Context, _ uuid.UUID, _ domain.AccountStatus) error { return nil },
	}

	svc := NewAccountService(accountRepo, &mocks.CustomerRepositoryMock{})
	err := svc.UpdateStatus(context.Background(), uuid.New(), domain.StatusBlocked)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestGetAccountByID_OK(t *testing.T) {
	id := uuid.New()
	expected := &domain.Account{ID: id, Status: domain.StatusActive}

	accountRepo := &mocks.AccountRepositoryMock{
		GetByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Account, error) { return expected, nil },
	}

	svc := NewAccountService(accountRepo, &mocks.CustomerRepositoryMock{})
	acc, err := svc.GetByID(context.Background(), id)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if acc.ID != id {
		t.Errorf("expected ID %v, got %v", id, acc.ID)
	}
}

func TestCreateAccount_RepoError(t *testing.T) {
	customerID := uuid.New()
	repoErr := apperror.New("db_error", "db error", 500)

	accountRepo := &mocks.AccountRepositoryMock{
		CreateFn: func(_ context.Context, _ *domain.Account) error { return repoErr },
	}
	customerRepo := &mocks.CustomerRepositoryMock{
		GetByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Customer, error) {
			return &domain.Customer{ID: customerID}, nil
		},
	}

	svc := NewAccountService(accountRepo, customerRepo)
	acc, err := svc.Create(context.Background(), customerID, domain.TypeChecking)

	if acc != nil {
		t.Error("expected nil account")
	}
	if err != repoErr {
		t.Errorf("expected repoErr, got %v", err)
	}
}
