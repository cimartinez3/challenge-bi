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

// CreateTransferPair inserts both transaction legs atomically (ADR-005 + ADR-006).
// Order: INSERT out (related_tx_id=NULL) → INSERT in → UPDATE out.related_tx_id.
func (r *transactionRepo) CreateTransferPair(ctx context.Context, out, in *domain.Transaction) error {
	dbTx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer dbTx.Rollback(ctx)

	// Step 1 — insert transfer_out without related_tx_id
	_, err = dbTx.Exec(ctx,
		`INSERT INTO transactions
		 (id, account_id, type, amount, balance_after, reference, related_tx_id, status, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, NULL, $7, $8)`,
		out.ID, out.AccountID, out.Type, out.Amount, out.BalanceAfter,
		out.Reference, out.Status, out.CreatedAt,
	)
	if err != nil {
		return mapTxError(err)
	}

	// Step 2 — insert transfer_in referencing out
	_, err = dbTx.Exec(ctx,
		`INSERT INTO transactions
		 (id, account_id, type, amount, balance_after, reference, related_tx_id, status, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		in.ID, in.AccountID, in.Type, in.Amount, in.BalanceAfter,
		in.Reference, out.ID, in.Status, in.CreatedAt,
	)
	if err != nil {
		return mapTxError(err)
	}

	// Step 3 — link transfer_out back to transfer_in
	_, err = dbTx.Exec(ctx,
		`UPDATE transactions SET related_tx_id = $1 WHERE id = $2`,
		in.ID, out.ID,
	)
	if err != nil {
		return err
	}

	return dbTx.Commit(ctx)
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
