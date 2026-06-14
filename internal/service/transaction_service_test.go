package service

import (
	"context"
	"testing"
	"time"

	"github.com/carlosmartinez/challenge-bi/internal/apperror"
	"github.com/carlosmartinez/challenge-bi/internal/domain"
	"github.com/carlosmartinez/challenge-bi/internal/mocks"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func activeAccount(id uuid.UUID) *domain.Account {
	return &domain.Account{ID: id, Status: domain.StatusActive, Balance: decimal.NewFromInt(1000)}
}

// ── Deposit ───────────────────────────────────────────────────────────────────

func TestDeposit_OK(t *testing.T) {
	accID := uuid.New()
	credited := &domain.Account{ID: accID, Balance: decimal.NewFromInt(1100)}

	txRepo := &mocks.TransactionRepositoryMock{
		GetByReferenceFn: func(_ context.Context, _ string) (*domain.Transaction, error) { return nil, nil },
		CreateFn:         func(_ context.Context, _ *domain.Transaction) error { return nil },
	}
	accountRepo := &mocks.AccountRepositoryMock{
		GetByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Account, error) { return activeAccount(accID), nil },
		CreditFn:  func(_ context.Context, _ uuid.UUID, _ decimal.Decimal) (*domain.Account, error) { return credited, nil },
	}

	svc := NewTransactionService(accountRepo, txRepo)
	result, err := svc.Deposit(context.Background(), accID, decimal.NewFromInt(100), "ref-001")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Duplicate {
		t.Error("expected Duplicate=false")
	}
	if result.Transaction.Type != domain.TxDeposit {
		t.Errorf("expected type deposit, got %v", result.Transaction.Type)
	}
	if result.Transaction.Status != domain.TxSuccess {
		t.Errorf("expected status success, got %v", result.Transaction.Status)
	}
}

func TestDeposit_InvalidAmount(t *testing.T) {
	svc := NewTransactionService(&mocks.AccountRepositoryMock{}, &mocks.TransactionRepositoryMock{})
	_, err := svc.Deposit(context.Background(), uuid.New(), decimal.Zero, "ref-x")
	if err != apperror.ErrInvalidAmount {
		t.Errorf("expected ErrInvalidAmount, got %v", err)
	}
}

func TestDeposit_AccountBlocked(t *testing.T) {
	accID := uuid.New()
	blocked := &domain.Account{ID: accID, Status: domain.StatusBlocked}

	txRepo := &mocks.TransactionRepositoryMock{
		GetByReferenceFn: func(_ context.Context, _ string) (*domain.Transaction, error) { return nil, nil },
	}
	accountRepo := &mocks.AccountRepositoryMock{
		GetByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Account, error) { return blocked, nil },
	}

	svc := NewTransactionService(accountRepo, txRepo)
	_, err := svc.Deposit(context.Background(), accID, decimal.NewFromInt(50), "ref-002")

	if err != apperror.ErrAccountBlocked {
		t.Errorf("expected ErrAccountBlocked, got %v", err)
	}
}

func TestDeposit_AccountClosed(t *testing.T) {
	accID := uuid.New()
	closed := &domain.Account{ID: accID, Status: domain.StatusClosed}

	txRepo := &mocks.TransactionRepositoryMock{
		GetByReferenceFn: func(_ context.Context, _ string) (*domain.Transaction, error) { return nil, nil },
	}
	accountRepo := &mocks.AccountRepositoryMock{
		GetByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Account, error) { return closed, nil },
	}

	svc := NewTransactionService(accountRepo, txRepo)
	_, err := svc.Deposit(context.Background(), accID, decimal.NewFromInt(50), "ref-003")

	if err != apperror.ErrAccountClosed {
		t.Errorf("expected ErrAccountClosed, got %v", err)
	}
}

func TestDeposit_DuplicateReference(t *testing.T) {
	existing := &domain.Transaction{Reference: "ref-dup", Status: domain.TxSuccess}

	txRepo := &mocks.TransactionRepositoryMock{
		GetByReferenceFn: func(_ context.Context, _ string) (*domain.Transaction, error) { return existing, nil },
	}

	svc := NewTransactionService(&mocks.AccountRepositoryMock{}, txRepo)
	result, err := svc.Deposit(context.Background(), uuid.New(), decimal.NewFromInt(50), "ref-dup")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !result.Duplicate {
		t.Error("expected Duplicate=true")
	}
	if result.Transaction != existing {
		t.Error("expected existing transaction to be returned")
	}
}

// ── Withdraw ──────────────────────────────────────────────────────────────────

func TestWithdraw_OK(t *testing.T) {
	accID := uuid.New()
	debited := &domain.Account{ID: accID, Balance: decimal.NewFromInt(900)}

	txRepo := &mocks.TransactionRepositoryMock{
		GetByReferenceFn: func(_ context.Context, _ string) (*domain.Transaction, error) { return nil, nil },
		CreateFn:         func(_ context.Context, _ *domain.Transaction) error { return nil },
	}
	accountRepo := &mocks.AccountRepositoryMock{
		GetByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Account, error) { return activeAccount(accID), nil },
		DebitFn:   func(_ context.Context, _ uuid.UUID, _ decimal.Decimal) (*domain.Account, error) { return debited, nil },
	}

	svc := NewTransactionService(accountRepo, txRepo)
	result, err := svc.Withdraw(context.Background(), accID, decimal.NewFromInt(100), "ref-004")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Transaction.Type != domain.TxWithdrawal {
		t.Errorf("expected type withdrawal, got %v", result.Transaction.Type)
	}
}

func TestWithdraw_InsufficientFunds(t *testing.T) {
	accID := uuid.New()

	txRepo := &mocks.TransactionRepositoryMock{
		GetByReferenceFn: func(_ context.Context, _ string) (*domain.Transaction, error) { return nil, nil },
	}
	accountRepo := &mocks.AccountRepositoryMock{
		GetByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Account, error) { return activeAccount(accID), nil },
		DebitFn:   func(_ context.Context, _ uuid.UUID, _ decimal.Decimal) (*domain.Account, error) { return nil, apperror.ErrInsufficientFunds },
	}

	svc := NewTransactionService(accountRepo, txRepo)
	_, err := svc.Withdraw(context.Background(), accID, decimal.NewFromInt(9999), "ref-005")

	if err != apperror.ErrInsufficientFunds {
		t.Errorf("expected ErrInsufficientFunds, got %v", err)
	}
}

func TestDeposit_AccountNotFound(t *testing.T) {
	txRepo := &mocks.TransactionRepositoryMock{
		GetByReferenceFn: func(_ context.Context, _ string) (*domain.Transaction, error) { return nil, nil },
	}
	accountRepo := &mocks.AccountRepositoryMock{
		GetByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Account, error) {
			return nil, apperror.ErrAccountNotFound
		},
	}

	svc := NewTransactionService(accountRepo, txRepo)
	_, err := svc.Deposit(context.Background(), uuid.New(), decimal.NewFromInt(50), "ref-anf")

	if err != apperror.ErrAccountNotFound {
		t.Errorf("expected ErrAccountNotFound, got %v", err)
	}
}

func TestDeposit_CreateTxError(t *testing.T) {
	accID := uuid.New()
	repoErr := apperror.New("db_error", "db error", 500)

	txRepo := &mocks.TransactionRepositoryMock{
		GetByReferenceFn: func(_ context.Context, _ string) (*domain.Transaction, error) { return nil, nil },
		CreateFn:         func(_ context.Context, _ *domain.Transaction) error { return repoErr },
	}
	accountRepo := &mocks.AccountRepositoryMock{
		GetByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Account, error) { return activeAccount(accID), nil },
		CreditFn: func(_ context.Context, _ uuid.UUID, _ decimal.Decimal) (*domain.Account, error) {
			return activeAccount(accID), nil
		},
	}

	svc := NewTransactionService(accountRepo, txRepo)
	_, err := svc.Deposit(context.Background(), accID, decimal.NewFromInt(50), "ref-cterr")

	if err != repoErr {
		t.Errorf("expected repoErr, got %v", err)
	}
}

func TestDeposit_CheckDuplicateError(t *testing.T) {
	repoErr := apperror.New("db_error", "db error", 500)
	txRepo := &mocks.TransactionRepositoryMock{
		GetByReferenceFn: func(_ context.Context, _ string) (*domain.Transaction, error) { return nil, repoErr },
	}

	svc := NewTransactionService(&mocks.AccountRepositoryMock{}, txRepo)
	_, err := svc.Deposit(context.Background(), uuid.New(), decimal.NewFromInt(50), "ref-err")

	if err != repoErr {
		t.Errorf("expected repoErr, got %v", err)
	}
}

func TestWithdraw_AccountNotFound(t *testing.T) {
	txRepo := &mocks.TransactionRepositoryMock{
		GetByReferenceFn: func(_ context.Context, _ string) (*domain.Transaction, error) { return nil, nil },
	}
	accountRepo := &mocks.AccountRepositoryMock{
		GetByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Account, error) {
			return nil, apperror.ErrAccountNotFound
		},
	}

	svc := NewTransactionService(accountRepo, txRepo)
	_, err := svc.Withdraw(context.Background(), uuid.New(), decimal.NewFromInt(50), "ref-w-anf")

	if err != apperror.ErrAccountNotFound {
		t.Errorf("expected ErrAccountNotFound, got %v", err)
	}
}

func TestWithdraw_CreateTxError(t *testing.T) {
	accID := uuid.New()
	repoErr := apperror.New("db_error", "db error", 500)

	txRepo := &mocks.TransactionRepositoryMock{
		GetByReferenceFn: func(_ context.Context, _ string) (*domain.Transaction, error) { return nil, nil },
		CreateFn:         func(_ context.Context, _ *domain.Transaction) error { return repoErr },
	}
	accountRepo := &mocks.AccountRepositoryMock{
		GetByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Account, error) { return activeAccount(accID), nil },
		DebitFn: func(_ context.Context, _ uuid.UUID, _ decimal.Decimal) (*domain.Account, error) {
			return activeAccount(accID), nil
		},
	}

	svc := NewTransactionService(accountRepo, txRepo)
	_, err := svc.Withdraw(context.Background(), accID, decimal.NewFromInt(50), "ref-w-cterr")

	if err != repoErr {
		t.Errorf("expected repoErr, got %v", err)
	}
}

func TestWithdraw_CheckDuplicateError(t *testing.T) {
	repoErr := apperror.New("db_error", "db error", 500)
	txRepo := &mocks.TransactionRepositoryMock{
		GetByReferenceFn: func(_ context.Context, _ string) (*domain.Transaction, error) { return nil, repoErr },
	}

	svc := NewTransactionService(&mocks.AccountRepositoryMock{}, txRepo)
	_, err := svc.Withdraw(context.Background(), uuid.New(), decimal.NewFromInt(50), "ref-err")

	if err != repoErr {
		t.Errorf("expected repoErr, got %v", err)
	}
}

func TestWithdraw_AccountBlocked(t *testing.T) {
	accID := uuid.New()
	blocked := &domain.Account{ID: accID, Status: domain.StatusBlocked}

	txRepo := &mocks.TransactionRepositoryMock{
		GetByReferenceFn: func(_ context.Context, _ string) (*domain.Transaction, error) { return nil, nil },
	}
	accountRepo := &mocks.AccountRepositoryMock{
		GetByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Account, error) { return blocked, nil },
	}

	svc := NewTransactionService(accountRepo, txRepo)
	_, err := svc.Withdraw(context.Background(), accID, decimal.NewFromInt(50), "ref-006")

	if err != apperror.ErrAccountBlocked {
		t.Errorf("expected ErrAccountBlocked, got %v", err)
	}
}

func TestListByAccount_OK(t *testing.T) {
	accID := uuid.New()
	expected := []domain.Transaction{{AccountID: accID}}

	txRepo := &mocks.TransactionRepositoryMock{
		ListByAccountFn: func(_ context.Context, _ uuid.UUID, _ int, _ *time.Time) ([]domain.Transaction, error) {
			return expected, nil
		},
	}
	accountRepo := &mocks.AccountRepositoryMock{
		GetByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Account, error) { return activeAccount(accID), nil },
	}

	svc := NewTransactionService(accountRepo, txRepo)
	txs, err := svc.ListByAccount(context.Background(), accID, 10, nil)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(txs) != 1 {
		t.Errorf("expected 1 transaction, got %d", len(txs))
	}
}

func TestListByAccount_AccountNotFound(t *testing.T) {
	txRepo := &mocks.TransactionRepositoryMock{}
	accountRepo := &mocks.AccountRepositoryMock{
		GetByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Account, error) {
			return nil, apperror.ErrAccountNotFound
		},
	}

	svc := NewTransactionService(accountRepo, txRepo)
	_, err := svc.ListByAccount(context.Background(), uuid.New(), 10, nil)

	if err != apperror.ErrAccountNotFound {
		t.Errorf("expected ErrAccountNotFound, got %v", err)
	}
}

func TestGetByReference_OK(t *testing.T) {
	expected := &domain.Transaction{Reference: "ref-get-01"}

	txRepo := &mocks.TransactionRepositoryMock{
		GetByReferenceFn: func(_ context.Context, _ string) (*domain.Transaction, error) { return expected, nil },
	}

	svc := NewTransactionService(&mocks.AccountRepositoryMock{}, txRepo)
	tx, err := svc.GetByReference(context.Background(), "ref-get-01")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if tx.Reference != "ref-get-01" {
		t.Errorf("expected reference ref-get-01, got %v", tx.Reference)
	}
}

func TestGetByReference_NotFound(t *testing.T) {
	txRepo := &mocks.TransactionRepositoryMock{
		GetByReferenceFn: func(_ context.Context, _ string) (*domain.Transaction, error) { return nil, nil },
	}

	svc := NewTransactionService(&mocks.AccountRepositoryMock{}, txRepo)
	_, err := svc.GetByReference(context.Background(), "ref-missing")

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetByReference_RepoError(t *testing.T) {
	repoErr := apperror.New("db_error", "db error", 500)
	txRepo := &mocks.TransactionRepositoryMock{
		GetByReferenceFn: func(_ context.Context, _ string) (*domain.Transaction, error) { return nil, repoErr },
	}

	svc := NewTransactionService(&mocks.AccountRepositoryMock{}, txRepo)
	_, err := svc.GetByReference(context.Background(), "ref-err")

	if err != repoErr {
		t.Errorf("expected repoErr, got %v", err)
	}
}

func TestListByAccount_ZeroLimit(t *testing.T) {
	accID := uuid.New()
	txRepo := &mocks.TransactionRepositoryMock{
		ListByAccountFn: func(_ context.Context, _ uuid.UUID, limit int, _ *time.Time) ([]domain.Transaction, error) {
			if limit != 20 {
				return nil, apperror.New("bad_limit", "expected default limit 20", 400)
			}
			return []domain.Transaction{}, nil
		},
	}
	accountRepo := &mocks.AccountRepositoryMock{
		GetByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Account, error) { return activeAccount(accID), nil },
	}

	svc := NewTransactionService(accountRepo, txRepo)
	_, err := svc.ListByAccount(context.Background(), accID, 0, nil)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
