package dto

import (
	"time"

	"github.com/google/uuid"
)

type ProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ProjectResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OwnerId     uuid.UUID `json:"owner_id"`
	CreatedAt   time.Time `json:"created_at"`
}
