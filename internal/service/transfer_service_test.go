package service

import (
	"context"
	"testing"

	"github.com/carlosmartinez/challenge-bi/internal/apperror"
	"github.com/carlosmartinez/challenge-bi/internal/domain"
	"github.com/carlosmartinez/challenge-bi/internal/mocks"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func TestTransfer_OK(t *testing.T) {
	fromID, toID := uuid.New(), uuid.New()
	from := &domain.Account{ID: fromID, Status: domain.StatusActive}
	to := &domain.Account{ID: toID, Status: domain.StatusActive}
	updatedFrom := &domain.Account{ID: fromID, Balance: decimal.NewFromInt(400)}
	updatedTo := &domain.Account{ID: toID, Balance: decimal.NewFromInt(300)}

	txRepo := &mocks.TransactionRepositoryMock{
		GetByReferenceFn: func(_ context.Context, _ string) (*domain.Transaction, error) { return nil, nil },
		ExecuteTransferFn: func(_ context.Context, _, _ uuid.UUID, _ decimal.Decimal, _, _ *domain.Transaction) (*domain.Account, *domain.Account, error) {
			return updatedFrom, updatedTo, nil
		},
	}
	accountRepo := &mocks.AccountRepositoryMock{
		GetByIDFn: func(_ context.Context, id uuid.UUID) (*domain.Account, error) {
			if id == fromID {
				return from, nil
			}
			return to, nil
		},
	}

	svc := NewTransferService(accountRepo, txRepo)
	result, err := svc.Transfer(context.Background(), fromID, toID, decimal.NewFromInt(100), "ref-t01")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Duplicate {
		t.Error("expected Duplicate=false")
	}
	if result.Out.Type != domain.TxTransferOut {
		t.Errorf("expected TxTransferOut, got %v", result.Out.Type)
	}
	if result.In.Type != domain.TxTransferIn {
		t.Errorf("expected TxTransferIn, got %v", result.In.Type)
	}
	if !result.Out.BalanceAfter.Equal(updatedFrom.Balance) {
		t.Errorf("expected BalanceAfter %v, got %v", updatedFrom.Balance, result.Out.BalanceAfter)
	}
}

func TestTransfer_InvalidAmount(t *testing.T) {
	svc := NewTransferService(&mocks.AccountRepositoryMock{}, &mocks.TransactionRepositoryMock{})
	_, err := svc.Transfer(context.Background(), uuid.New(), uuid.New(), decimal.Zero, "ref-t02")
	if err != apperror.ErrInvalidAmount {
		t.Errorf("expected ErrInvalidAmount, got %v", err)
	}
}

func TestTransfer_FromAccountBlocked(t *testing.T) {
	fromID, toID := uuid.New(), uuid.New()
	blocked := &domain.Account{ID: fromID, Status: domain.StatusBlocked}

	txRepo := &mocks.TransactionRepositoryMock{
		GetByReferenceFn: func(_ context.Context, _ string) (*domain.Transaction, error) { return nil, nil },
	}
	accountRepo := &mocks.AccountRepositoryMock{
		GetByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Account, error) { return blocked, nil },
	}

	svc := NewTransferService(accountRepo, txRepo)
	_, err := svc.Transfer(context.Background(), fromID, toID, decimal.NewFromInt(50), "ref-t03")

	if err != apperror.ErrAccountBlocked {
		t.Errorf("expected ErrAccountBlocked, got %v", err)
	}
}

func TestTransfer_ToAccountClosed(t *testing.T) {
	fromID, toID := uuid.New(), uuid.New()
	from := &domain.Account{ID: fromID, Status: domain.StatusActive}
	closed := &domain.Account{ID: toID, Status: domain.StatusClosed}

	txRepo := &mocks.TransactionRepositoryMock{
		GetByReferenceFn: func(_ context.Context, _ string) (*domain.Transaction, error) { return nil, nil },
	}
	accountRepo := &mocks.AccountRepositoryMock{
		GetByIDFn: func(_ context.Context, id uuid.UUID) (*domain.Account, error) {
			if id == fromID {
				return from, nil
			}
			return closed, nil
		},
	}

	svc := NewTransferService(accountRepo, txRepo)
	_, err := svc.Transfer(context.Background(), fromID, toID, decimal.NewFromInt(50), "ref-t04")

	if err != apperror.ErrAccountClosed {
		t.Errorf("expected ErrAccountClosed, got %v", err)
	}
}

func TestTransfer_InsufficientFunds(t *testing.T) {
	fromID, toID := uuid.New(), uuid.New()
	from := &domain.Account{ID: fromID, Status: domain.StatusActive}
	to := &domain.Account{ID: toID, Status: domain.StatusActive}

	txRepo := &mocks.TransactionRepositoryMock{
		GetByReferenceFn: func(_ context.Context, _ string) (*domain.Transaction, error) { return nil, nil },
		ExecuteTransferFn: func(_ context.Context, _, _ uuid.UUID, _ decimal.Decimal, _, _ *domain.Transaction) (*domain.Account, *domain.Account, error) {
			return nil, nil, apperror.ErrInsufficientFunds
		},
	}
	accountRepo := &mocks.AccountRepositoryMock{
		GetByIDFn: func(_ context.Context, id uuid.UUID) (*domain.Account, error) {
			if id == fromID {
				return from, nil
			}
			return to, nil
		},
	}

	svc := NewTransferService(accountRepo, txRepo)
	_, err := svc.Transfer(context.Background(), fromID, toID, decimal.NewFromInt(9999), "ref-t05")

	if err != apperror.ErrInsufficientFunds {
		t.Errorf("expected ErrInsufficientFunds, got %v", err)
	}
}

func TestTransfer_FromAccountNotFound(t *testing.T) {
	txRepo := &mocks.TransactionRepositoryMock{
		GetByReferenceFn: func(_ context.Context, _ string) (*domain.Transaction, error) { return nil, nil },
	}
	accountRepo := &mocks.AccountRepositoryMock{
		GetByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Account, error) {
			return nil, apperror.ErrAccountNotFound
		},
	}

	svc := NewTransferService(accountRepo, txRepo)
	_, err := svc.Transfer(context.Background(), uuid.New(), uuid.New(), decimal.NewFromInt(50), "ref-t-anf")

	if err != apperror.ErrAccountNotFound {
		t.Errorf("expected ErrAccountNotFound, got %v", err)
	}
}

func TestTransfer_ToAccountNotFound(t *testing.T) {
	fromID, toID := uuid.New(), uuid.New()
	from := &domain.Account{ID: fromID, Status: domain.StatusActive}

	txRepo := &mocks.TransactionRepositoryMock{
		GetByReferenceFn: func(_ context.Context, _ string) (*domain.Transaction, error) { return nil, nil },
	}
	accountRepo := &mocks.AccountRepositoryMock{
		GetByIDFn: func(_ context.Context, id uuid.UUID) (*domain.Account, error) {
			if id == fromID {
				return from, nil
			}
			return nil, apperror.ErrAccountNotFound
		},
	}

	svc := NewTransferService(accountRepo, txRepo)
	_, err := svc.Transfer(context.Background(), fromID, toID, decimal.NewFromInt(50), "ref-t-to-anf")

	if err != apperror.ErrAccountNotFound {
		t.Errorf("expected ErrAccountNotFound, got %v", err)
	}
}

func TestTransfer_DuplicateInLookupError(t *testing.T) {
	fromID, toID := uuid.New(), uuid.New()
	existingOut := &domain.Transaction{Reference: "ref-tdup2_out"}
	repoErr := apperror.New("db_error", "db error", 500)

	txRepo := &mocks.TransactionRepositoryMock{
		GetByReferenceFn: func(_ context.Context, ref string) (*domain.Transaction, error) {
			if ref == "ref-tdup2_out" {
				return existingOut, nil
			}
			return nil, repoErr
		},
	}

	svc := NewTransferService(&mocks.AccountRepositoryMock{}, txRepo)
	_, err := svc.Transfer(context.Background(), fromID, toID, decimal.NewFromInt(50), "ref-tdup2")

	if err != repoErr {
		t.Errorf("expected repoErr, got %v", err)
	}
}

func TestTransfer_CheckDuplicateError(t *testing.T) {
	repoErr := apperror.New("db_error", "db error", 500)
	txRepo := &mocks.TransactionRepositoryMock{
		GetByReferenceFn: func(_ context.Context, _ string) (*domain.Transaction, error) { return nil, repoErr },
	}

	svc := NewTransferService(&mocks.AccountRepositoryMock{}, txRepo)
	_, err := svc.Transfer(context.Background(), uuid.New(), uuid.New(), decimal.NewFromInt(50), "ref-err")

	if err != repoErr {
		t.Errorf("expected repoErr, got %v", err)
	}
}

func TestTransfer_Duplicate(t *testing.T) {
	existing := &domain.Transaction{Reference: "ref-t06_out"}

	txRepo := &mocks.TransactionRepositoryMock{
		GetByReferenceFn: func(_ context.Context, ref string) (*domain.Transaction, error) {
			if ref == "ref-t06_out" {
				return existing, nil
			}
			return nil, nil
		},
	}

	svc := NewTransferService(&mocks.AccountRepositoryMock{}, txRepo)
	result, err := svc.Transfer(context.Background(), uuid.New(), uuid.New(), decimal.NewFromInt(50), "ref-t06")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !result.Duplicate {
		t.Error("expected Duplicate=true")
	}
}
