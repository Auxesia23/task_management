package models

import (
	"time"

	"github.com/google/uuid"
)

type Invitation struct {
	ID           uuid.UUID `db:"id"`
	ProjectName  string    `db:"project_name"`
	UserEmail    string    `db:"user_email"`
	InviterEmail string    `db:"inviter_email"`
	Status       string    `db:"status"`
	CreatedAt    time.Time `db:"created_at"`
}
