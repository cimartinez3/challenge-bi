package repository

import (
	"context"
	"errors"

	"github.com/carlosmartinez/challenge-bi/internal/apperror"
	"github.com/carlosmartinez/challenge-bi/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type accountRepo struct {
	db *pgxpool.Pool
}

func NewAccountRepo(db *pgxpool.Pool) *accountRepo {
	return &accountRepo{db: db}
}

func (r *accountRepo) Create(ctx context.Context, a *domain.Account) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO accounts
		 (id, customer_id, account_number, type, currency, balance, status, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		a.ID, a.CustomerID, a.AccountNumber, a.Type, a.Currency,
		a.Balance, a.Status, a.CreatedAt, a.UpdatedAt,
	)
	return err
}

func (r *accountRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	a := &domain.Account{}
	err := r.db.QueryRow(ctx,
		`SELECT id, customer_id, account_number, type, currency,
		        balance, status, created_at, updated_at
		 FROM accounts WHERE id = $1`,
		id,
	).Scan(
		&a.ID, &a.CustomerID, &a.AccountNumber, &a.Type, &a.Currency,
		&a.Balance, &a.Status, &a.CreatedAt, &a.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperror.ErrAccountNotFound
	}

	if err != nil {
		return nil, err
	}

	return a, nil
}

func (r *accountRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.AccountStatus) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE accounts SET status = $1, updated_at = now() WHERE id = $2`,
		status, id,
	)

	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return apperror.ErrAccountNotFound
	}

	return nil
}

// Debit decrements the balance atomically (ADR-003).
// The WHERE balance >= amount clause makes it a single atomic read-modify-write;
// if no row is updated the balance was insufficient.
func (r *accountRepo) Debit(ctx context.Context, id uuid.UUID, amount decimal.Decimal) (*domain.Account, error) {
	a := &domain.Account{}
	err := r.db.QueryRow(ctx,
		`UPDATE accounts
		 SET balance = balance - $1, updated_at = now()
		 WHERE id = $2 AND balance >= $1 AND status = 'active'
		 RETURNING id, customer_id, account_number, type, currency,
		           balance, status, created_at, updated_at`,
		amount, id,
	).Scan(
		&a.ID, &a.CustomerID, &a.AccountNumber, &a.Type, &a.Currency,
		&a.Balance, &a.Status, &a.CreatedAt, &a.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperror.ErrInsufficientFunds
	}
	if err != nil {
		return nil, err
	}
	return a, nil
}

// Credit increments the balance atomically.
func (r *accountRepo) Credit(ctx context.Context, id uuid.UUID, amount decimal.Decimal) (*domain.Account, error) {
	a := &domain.Account{}
	err := r.db.QueryRow(ctx,
		`UPDATE accounts
		 SET balance = balance + $1, updated_at = now()
		 WHERE id = $2 AND status = 'active'
		 RETURNING id, customer_id, account_number, type, currency,
		           balance, status, created_at, updated_at`,
		amount, id,
	).Scan(
		&a.ID, &a.CustomerID, &a.AccountNumber, &a.Type, &a.Currency,
		&a.Balance, &a.Status, &a.CreatedAt, &a.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperror.ErrAccountNotFound
	}
	if err != nil {
		return nil, err
	}
	return a, nil
}
