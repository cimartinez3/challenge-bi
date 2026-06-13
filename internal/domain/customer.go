package domain

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID        uuid.UUID
	FullName  string
	Email     string
	CreatedAt time.Time
}
