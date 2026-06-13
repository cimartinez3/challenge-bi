package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/carlosmartinez/challenge-bi/internal/apperror"
	"github.com/carlosmartinez/challenge-bi/internal/domain"
	"github.com/carlosmartinez/challenge-bi/internal/ports"
)

type TransactionService struct {
	accountRepo ports.AccountRepository
	txRepo      ports.TransactionRepository
}

func NewTransactionService(accountRepo ports.AccountRepository, txRepo ports.TransactionRepository) *TransactionService {
	return &TransactionService{accountRepo: accountRepo, txRepo: txRepo}
}

type TxResult struct {
	Transaction *domain.Transaction
	Duplicate   bool // true when the reference already existed (idempotent replay)
}

func (s *TransactionService) Deposit(ctx context.Context, accountID uuid.UUID, amount decimal.Decimal, reference string) (*TxResult, error) {
	if amount.LessThanOrEqual(decimal.Zero) {
		return nil, apperror.ErrInvalidAmount
	}

	// Check for duplicate BEFORE touching the balance — prevents double credit
	if result, err := s.checkDuplicate(ctx, reference); result != nil || err != nil {
		return result, err
	}

	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return nil, err
	}
	if err := account.IsOperable(); err != nil {
		return nil, err
	}

	updated, err := s.accountRepo.Credit(ctx, accountID, amount)
	if err != nil {
		return nil, err
	}

	tx := &domain.Transaction{
		ID:           uuid.New(),
		AccountID:    accountID,
		Type:         domain.TxDeposit,
		Amount:       amount,
		BalanceAfter: updated.Balance,
		Reference:    reference,
		Status:       domain.TxSuccess,
		CreatedAt:    time.Now().UTC(),
	}
	if err := s.txRepo.Create(ctx, tx); err != nil {
		return nil, err
	}
	return &TxResult{Transaction: tx, Duplicate: false}, nil
}

func (s *TransactionService) Withdraw(ctx context.Context, accountID uuid.UUID, amount decimal.Decimal, reference string) (*TxResult, error) {
	if amount.LessThanOrEqual(decimal.Zero) {
		return nil, apperror.ErrInvalidAmount
	}

	// Check for duplicate BEFORE touching the balance — prevents double debit
	if result, err := s.checkDuplicate(ctx, reference); result != nil || err != nil {
		return result, err
	}

	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return nil, err
	}
	if err := account.IsOperable(); err != nil {
		return nil, err
	}

	updated, err := s.accountRepo.Debit(ctx, accountID, amount)
	if err != nil {
		return nil, err
	}

	tx := &domain.Transaction{
		ID:           uuid.New(),
		AccountID:    accountID,
		Type:         domain.TxWithdrawal,
		Amount:       amount,
		BalanceAfter: updated.Balance,
		Reference:    reference,
		Status:       domain.TxSuccess,
		CreatedAt:    time.Now().UTC(),
	}
	if err := s.txRepo.Create(ctx, tx); err != nil {
		return nil, err
	}
	return &TxResult{Transaction: tx, Duplicate: false}, nil
}

func (s *TransactionService) ListByAccount(ctx context.Context, accountID uuid.UUID, limit int, before *time.Time) ([]domain.Transaction, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if _, err := s.accountRepo.GetByID(ctx, accountID); err != nil {
		return nil, err
	}
	return s.txRepo.ListByAccount(ctx, accountID, limit, before)
}

func (s *TransactionService) GetByReference(ctx context.Context, ref string) (*domain.Transaction, error) {
	tx, err := s.txRepo.GetByReference(ctx, ref)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, apperror.New("transaction_not_found", "transaction not found", 404)
	}
	return tx, nil
}

// checkDuplicate returns the existing TxResult if the reference already exists, nil otherwise.
func (s *TransactionService) checkDuplicate(ctx context.Context, reference string) (*TxResult, error) {
	existing, err := s.txRepo.GetByReference(ctx, reference)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return &TxResult{Transaction: existing, Duplicate: true}, nil
	}
	return nil, nil
}
