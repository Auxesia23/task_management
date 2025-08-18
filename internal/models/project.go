package models

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	OwnerID     uuid.UUID `db:"owner_id"`
	CreatedAt   time.Time `db:"created_at"`
}
