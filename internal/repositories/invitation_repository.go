package repositories

import (
	"context"
	"errors"

	"github.com/Auxesia23/task_management/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type InvitationRepository interface {
	Create(ctx context.Context, projectId, userId, inviterId *uuid.UUID) error
	ReadByUser(ctx context.Context, userId *uuid.UUID, status *string) (*[]models.Invitation, error)
	Update(ctx context.Context, invitationId, userId *uuid.UUID, status *string) error
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
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return errors.New("invitation already exists")
			}
		}
		return errors.New("invitation creation failed")
	}
	return nil
}

func (r *invitationRepository) ReadByUser(ctx context.Context, userId *uuid.UUID, status *string) (*[]models.Invitation, error) {
	query := `
		SELECT 
			i.id AS id,
			p.name AS project_name,
			u.email AS user_email,
			inviter.email AS inviter_email,
			i.status AS status,
			i.created_at AS created_at
		FROM invitations i
		JOIN projects p ON p.id = i.project_id
		JOIN users u ON u.id = i.user_id
		JOIN users inviter ON inviter.id = i.inviter_id
		WHERE i.user_id = $1 AND status = $2;
		`

	var invitations []models.Invitation
	if err := r.db.SelectContext(ctx, &invitations, query, userId, status); err != nil {
		return nil, errors.New("error reading invitations")
	}

	return &invitations, nil
}

func (r *invitationRepository) Update(ctx context.Context, invitationId, userId *uuid.UUID, status *string) error {
	query := `
		UPDATE invitations 
		SET status = $1
		WHERE id = $2 AND user_id = $3 AND status = 'pending';
		`

	result, err := r.db.ExecContext(ctx, query, status, invitationId, userId)
	if err != nil {
		return errors.New("error updating invitation")
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return errors.New("error updating invitation")
	}
	if rows == 0 {
		return errors.New("invitation not found or has been accepted/rejected")
	}
	return nil
}
