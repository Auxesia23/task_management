package services

import (
	"context"
	"errors"

	"github.com/Auxesia23/task_management/internal/dto"
	"github.com/Auxesia23/task_management/internal/repositories"
	"github.com/google/uuid"
)

type InvitationService interface {
	CreateInvitation(ctx context.Context, projectId, userId, inviterId *uuid.UUID) error
	GetInvitation(ctx context.Context, userId *uuid.UUID, status *string) (*[]dto.InvitationResponse, error)
	UpdateInvitation(ctx context.Context, userId, invitationId *uuid.UUID, status *string) error
}

type invitationService struct {
	invitationRepo repositories.InvitationRepository
	projectRepo    repositories.ProjectRepository
}

func NewInvitationService(invitationRepo repositories.InvitationRepository, projectRepo repositories.ProjectRepository) InvitationService {
	return &invitationService{
		invitationRepo,
		projectRepo,
	}
}

func (s *invitationService) CreateInvitation(ctx context.Context, projectId, userId, inviterId *uuid.UUID) error {
	ok, err := s.projectRepo.OwnerCheck(ctx, projectId, inviterId)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("project doesn't exist")
	}

	if err = s.invitationRepo.Create(ctx, projectId, userId, inviterId); err != nil {
		return err
	}
	return nil
}

func (s *invitationService) GetInvitation(ctx context.Context, userId *uuid.UUID, status *string) (*[]dto.InvitationResponse, error) {
	invitations, err := s.invitationRepo.ReadByUser(ctx, userId, status)
	if err != nil {
		return nil, err
	}

	var response []dto.InvitationResponse
	for _, invitation := range *invitations {
		response = append(response, dto.InvitationResponse{
			ID:           invitation.ID,
			ProjectName:  invitation.ProjectName,
			UserEmail:    invitation.UserEmail,
			InviterEmail: invitation.InviterEmail,
			Status:       invitation.Status,
			CreatedAt:    invitation.CreatedAt,
		})
	}
	return &response, nil
}

func (s *invitationService) UpdateInvitation(ctx context.Context, userId, invitationId *uuid.UUID, status *string) error {
	if status == nil {
		return errors.New("status is required")
	}
	if *status != "accept" && *status != "reject" {
		return errors.New("status must be accept/reject")
	}

	var newStatus string
	switch *status {
	case "accept":
		newStatus = "accepted"
	case "reject":
		newStatus = "rejected"
	}
	if err := s.invitationRepo.Update(ctx, invitationId, userId, &newStatus); err != nil {
		return err
	}
	return nil
}
