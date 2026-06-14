package service

import (
	"context"
	"errors"
	"testing"

	"github.com/carlosmartinez/challenge-bi/internal/domain"
	"github.com/carlosmartinez/challenge-bi/internal/mocks"
	"github.com/google/uuid"
)

func TestCreateCustomer_OK(t *testing.T) {
	repo := &mocks.CustomerRepositoryMock{
		CreateFn: func(_ context.Context, c *domain.Customer) error { return nil },
	}

	svc := NewCustomerService(repo)
	customer, err := svc.Create(context.Background(), "John Doe", "john@example.com")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if customer.FullName != "John Doe" {
		t.Errorf("expected FullName 'John Doe', got %v", customer.FullName)
	}
	if customer.Email != "john@example.com" {
		t.Errorf("expected Email 'john@example.com', got %v", customer.Email)
	}
	if customer.ID == uuid.Nil {
		t.Error("expected non-nil UUID")
	}
}

func TestCreateCustomer_RepoError(t *testing.T) {
	repoErr := errors.New("db error")
	repo := &mocks.CustomerRepositoryMock{
		CreateFn: func(_ context.Context, _ *domain.Customer) error { return repoErr },
	}

	svc := NewCustomerService(repo)
	customer, err := svc.Create(context.Background(), "Jane", "jane@example.com")

	if customer != nil {
		t.Error("expected nil customer")
	}
	if err != repoErr {
		t.Errorf("expected repoErr, got %v", err)
	}
}

func TestGetCustomerByID_OK(t *testing.T) {
	id := uuid.New()
	expected := &domain.Customer{ID: id, FullName: "John"}

	repo := &mocks.CustomerRepositoryMock{
		GetByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Customer, error) { return expected, nil },
	}

	svc := NewCustomerService(repo)
	customer, err := svc.GetByID(context.Background(), id)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if customer.ID != id {
		t.Errorf("expected ID %v, got %v", id, customer.ID)
	}
}
