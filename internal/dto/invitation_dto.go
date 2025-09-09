package dto

import (
	"time"

	"github.com/google/uuid"
)

type InvitationRequest struct {
	UserID uuid.UUID `json:"user_id"`
}

type InvitationResponse struct {
	ID           uuid.UUID `json:"id"`
	ProjectName  string    `json:"project_name"`
	UserEmail    string    `json:"user_email"`
	InviterEmail string    `json:"inviter_email"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}
