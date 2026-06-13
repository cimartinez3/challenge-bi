package repository

import (
	"context"
	"errors"

	"github.com/carlosmartinez/challenge-bi/internal/apperror"
	"github.com/carlosmartinez/challenge-bi/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type customerRepo struct {
	db *pgxpool.Pool
}

func NewCustomerRepo(db *pgxpool.Pool) *customerRepo {
	return &customerRepo{db: db}
}

func (r *customerRepo) Create(ctx context.Context, c *domain.Customer) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO customers (id, full_name, email, created_at)
		 VALUES ($1, $2, $3, $4)`,
		c.ID, c.FullName, c.Email, c.CreatedAt,
	)

	return err
}

func (r *customerRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Customer, error) {
	c := &domain.Customer{}
	err := r.db.QueryRow(ctx,
		`SELECT id, full_name, email, created_at
		 FROM customers WHERE id = $1`,
		id,
	).Scan(&c.ID, &c.FullName, &c.Email, &c.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperror.ErrCustomerNotFound
	}

	if err != nil {
		return nil, err
	}

	return c, nil
}
