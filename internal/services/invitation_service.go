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
	GetInvitation(ctx context.Context, userId *uuid.UUID) (*[]dto.InvitationResponse, error)
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

func (s *invitationService) GetInvitation(ctx context.Context, userId *uuid.UUID) (*[]dto.InvitationResponse, error) {
	invitations, err := s.invitationRepo.ReadByUser(ctx, userId)
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
