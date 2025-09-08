package models

import (
	"time"

	"github.com/google/uuid"
)

type Invitation struct {
	ID        uuid.UUID `db:"id"`
	ProjectID uuid.UUID `db:"project_id"`
	UserID    uuid.UUID `db:"user_id"`
	InviterID uuid.UUID `db:"inviter_id"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
}
