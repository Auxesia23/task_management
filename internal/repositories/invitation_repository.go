package repositories

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type InvitationRepository interface {
	Create(ctx context.Context, projectId, userId, inviterId *uuid.UUID) error
}

type invitationRepository struct {
	db *sqlx.DB
}

func NewInvitationRepository(db *sqlx.DB) InvitationRepository {
	return &invitationRepository{
		db: db,
	}
}

func (r *invitationRepository) Create(ctx context.Context, projectId, userId, inviterId *uuid.UUID) error {
	query := `
			INSERT INTO invitations (project_id, user_id, inviter_id)
			VALUES ($1, $2, $3)
		`
	_, err := r.db.ExecContext(ctx, query, projectId, userId, inviterId)
	if err != nil {
		return errors.New("invitation creation failed")
	}
	return nil
}
