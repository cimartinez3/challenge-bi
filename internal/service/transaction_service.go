package service

import (
	"context"
	"time"

	"github.com/carlosmartinez/challenge-bi/internal/apperror"
	"github.com/carlosmartinez/challenge-bi/internal/domain"
	"github.com/carlosmartinez/challenge-bi/internal/ports"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TransactionService struct {
	accountRepo ports.AccountRepository
	txRepo      ports.TransactionRepository
}

func NewTransactionService(accountRepo ports.AccountRepository, txRepo ports.TransactionRepository) *TransactionService {
	return &TransactionService{accountRepo: accountRepo, txRepo: txRepo}
}

func (s *TransactionService) Deposit(ctx context.Context, accountID uuid.UUID, amount decimal.Decimal, reference string) (*domain.Transaction, error) {
	if amount.LessThanOrEqual(decimal.Zero) {
		return nil, apperror.ErrInvalidAmount
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

	return tx, nil
}

func (s *TransactionService) Withdraw(ctx context.Context, accountID uuid.UUID, amount decimal.Decimal, reference string) (*domain.Transaction, error) {
	if amount.LessThanOrEqual(decimal.Zero) {
		return nil, apperror.ErrInvalidAmount
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

	return tx, nil
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
