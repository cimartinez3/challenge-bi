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

type TransferService struct {
	accountRepo ports.AccountRepository
	txRepo      ports.TransactionRepository
}

func NewTransferService(accountRepo ports.AccountRepository, txRepo ports.TransactionRepository) *TransferService {
	return &TransferService{accountRepo: accountRepo, txRepo: txRepo}
}

type TransferResult struct {
	Out       *domain.Transaction
	In        *domain.Transaction
	Duplicate bool
}

func (s *TransferService) Transfer(ctx context.Context, fromID, toID uuid.UUID, amount decimal.Decimal, reference string) (*TransferResult, error) {
	if amount.LessThanOrEqual(decimal.Zero) {
		return nil, apperror.ErrInvalidAmount
	}

	// Check for duplicate BEFORE touching any balance — prevents double debit/credit
	existing, err := s.txRepo.GetByReference(ctx, reference+"_out")
	if err != nil {
		return nil, err
	}
	if existing != nil {
		existingIn, err := s.txRepo.GetByReference(ctx, reference+"_in")
		if err != nil {
			return nil, err
		}
		return &TransferResult{Out: existing, In: existingIn, Duplicate: true}, nil
	}

	from, err := s.accountRepo.GetByID(ctx, fromID)
	if err != nil {
		return nil, err
	}
	if err := from.IsOperable(); err != nil {
		return nil, err
	}

	to, err := s.accountRepo.GetByID(ctx, toID)
	if err != nil {
		return nil, err
	}
	if err := to.IsOperable(); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	out := &domain.Transaction{
		ID:        uuid.New(),
		AccountID: fromID,
		Type:      domain.TxTransferOut,
		Amount:    amount,
		Reference: reference + "_out",
		Status:    domain.TxSuccess,
		CreatedAt: now,
	}
	in := &domain.Transaction{
		ID:        uuid.New(),
		AccountID: toID,
		Type:      domain.TxTransferIn,
		Amount:    amount,
		Reference: reference + "_in",
		Status:    domain.TxSuccess,
		CreatedAt: now,
	}

	updatedFrom, updatedTo, err := s.txRepo.ExecuteTransfer(ctx, fromID, toID, amount, out, in)
	if err != nil {
		return nil, err
	}

	out.BalanceAfter = updatedFrom.Balance
	in.BalanceAfter = updatedTo.Balance

	return &TransferResult{Out: out, In: in, Duplicate: false}, nil
}
