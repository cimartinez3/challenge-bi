package repository

import (
	"context"
	"errors"
	"time"

	"github.com/carlosmartinez/challenge-bi/internal/apperror"
	"github.com/carlosmartinez/challenge-bi/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

const pgErrUniqueViolation = "23505"

type transactionRepo struct {
	db *pgxpool.Pool
}

func NewTransactionRepo(db *pgxpool.Pool) *transactionRepo {
	return &transactionRepo{db: db}
}

func (r *transactionRepo) Create(ctx context.Context, tx *domain.Transaction) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO transactions
		 (id, account_id, type, amount, balance_after, reference, related_tx_id, status, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		tx.ID, tx.AccountID, tx.Type, tx.Amount, tx.BalanceAfter,
		tx.Reference, tx.RelatedTxID, tx.Status, tx.CreatedAt,
	)

	return mapTxError(err)
}

func (r *transactionRepo) GetByReference(ctx context.Context, ref string) (*domain.Transaction, error) {
	tx := &domain.Transaction{}
	err := r.db.QueryRow(ctx,
		`SELECT id, account_id, type, amount, balance_after,
		        reference, related_tx_id, status, created_at
		 FROM transactions WHERE reference = $1`,
		ref,
	).Scan(
		&tx.ID, &tx.AccountID, &tx.Type, &tx.Amount, &tx.BalanceAfter,
		&tx.Reference, &tx.RelatedTxID, &tx.Status, &tx.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (r *transactionRepo) ListByAccount(ctx context.Context, accountID uuid.UUID, limit int, before *time.Time) ([]domain.Transaction, error) {
	var (
		rows pgx.Rows
		err  error
	)

	if before != nil {
		rows, err = r.db.Query(ctx,
			`SELECT id, account_id, type, amount, balance_after,
			        reference, related_tx_id, status, created_at
			 FROM transactions
			 WHERE account_id = $1 AND created_at < $2
			 ORDER BY created_at DESC
			 LIMIT $3`,
			accountID, before, limit,
		)
	} else {
		rows, err = r.db.Query(ctx,
			`SELECT id, account_id, type, amount, balance_after,
			        reference, related_tx_id, status, created_at
			 FROM transactions
			 WHERE account_id = $1
			 ORDER BY created_at DESC
			 LIMIT $2`,
			accountID, limit,
		)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs []domain.Transaction
	for rows.Next() {
		var tx domain.Transaction
		if err := rows.Scan(
			&tx.ID, &tx.AccountID, &tx.Type, &tx.Amount, &tx.BalanceAfter,
			&tx.Reference, &tx.RelatedTxID, &tx.Status, &tx.CreatedAt,
		); err != nil {
			return nil, err
		}
		txs = append(txs, tx)
	}
	return txs, rows.Err()
}

func (r *transactionRepo) ExecuteTransfer(ctx context.Context, fromID, toID uuid.UUID, amount decimal.Decimal, out, in *domain.Transaction) (*domain.Account, *domain.Account, error) {
	dbTx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer dbTx.Rollback(ctx)

	fromAcc := &domain.Account{}
	err = dbTx.QueryRow(ctx,
		`UPDATE accounts
		 SET balance = balance - $1, updated_at = now()
		 WHERE id = $2 AND balance >= $1 AND status = 'active'
		 RETURNING id, customer_id, account_number, type, currency,
		           balance, status, created_at, updated_at`,
		amount, fromID,
	).Scan(
		&fromAcc.ID, &fromAcc.CustomerID, &fromAcc.AccountNumber, &fromAcc.Type, &fromAcc.Currency,
		&fromAcc.Balance, &fromAcc.Status, &fromAcc.CreatedAt, &fromAcc.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil, apperror.ErrInsufficientFunds
	}
	if err != nil {
		return nil, nil, err
	}

	toAcc := &domain.Account{}
	err = dbTx.QueryRow(ctx,
		`UPDATE accounts
		 SET balance = balance + $1, updated_at = now()
		 WHERE id = $2 AND status = 'active'
		 RETURNING id, customer_id, account_number, type, currency,
		           balance, status, created_at, updated_at`,
		amount, toID,
	).Scan(
		&toAcc.ID, &toAcc.CustomerID, &toAcc.AccountNumber, &toAcc.Type, &toAcc.Currency,
		&toAcc.Balance, &toAcc.Status, &toAcc.CreatedAt, &toAcc.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil, apperror.ErrAccountNotFound
	}
	if err != nil {
		return nil, nil, err
	}

	_, err = dbTx.Exec(ctx,
		`INSERT INTO transactions
		 (id, account_id, type, amount, balance_after, reference, related_tx_id, status, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, NULL, $7, $8)`,
		out.ID, out.AccountID, out.Type, out.Amount, out.BalanceAfter,
		out.Reference, out.Status, out.CreatedAt,
	)
	if err != nil {
		return nil, nil, mapTxError(err)
	}

	_, err = dbTx.Exec(ctx,
		`INSERT INTO transactions
		 (id, account_id, type, amount, balance_after, reference, related_tx_id, status, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		in.ID, in.AccountID, in.Type, in.Amount, in.BalanceAfter,
		in.Reference, out.ID, in.Status, in.CreatedAt,
	)
	if err != nil {
		return nil, nil, mapTxError(err)
	}

	_, err = dbTx.Exec(ctx,
		`UPDATE transactions SET related_tx_id = $1 WHERE id = $2`,
		in.ID, out.ID,
	)
	if err != nil {
		return nil, nil, err
	}

	if err := dbTx.Commit(ctx); err != nil {
		return nil, nil, err
	}
	return fromAcc, toAcc, nil
}

func mapTxError(err error) error {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) && pgErr.Code == pgErrUniqueViolation {
		return apperror.ErrDuplicateRef
	}
	return err
}
