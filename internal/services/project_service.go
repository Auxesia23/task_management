package services

import (
	"context"

	"github.com/Auxesia23/task_management/internal/dto"
	"github.com/Auxesia23/task_management/internal/repositories"
	"github.com/google/uuid"
)

type ProjectService interface {
	CreateProject(ctx context.Context, project *dto.ProjectRequest, ownerId *uuid.UUID) (*dto.ProjectResponse, error)
	GetAllProjects(ctx context.Context) (*[]dto.ProjectResponse, error)
	GetProjectByID(ctx context.Context, projectId *uuid.UUID) (*dto.ProjectResponse, error)
	UpdateProject(ctx context.Context, project *dto.ProjectRequest, projectId *uuid.UUID, ownerId *uuid.UUID) (*dto.ProjectResponse, error)
	DeleteProject(ctx context.Context, projectId *uuid.UUID, ownerId *uuid.UUID) error
}

type projectService struct {
	projectRepo repositories.ProjectRepository
}

func NewProjectService(projectRepo repositories.ProjectRepository) ProjectService {
	return &projectService{
		projectRepo,
	}
}

func (s *projectService) CreateProject(ctx context.Context, project *dto.ProjectRequest, ownerId *uuid.UUID) (*dto.ProjectResponse, error) {
	createdProject, err := s.projectRepo.Create(ctx, project, ownerId)
	if err != nil {
		return nil, err
	}

	response := &dto.ProjectResponse{
		ID:          createdProject.ID,
		Name:        createdProject.Name,
		Description: createdProject.Description,
		OwnerId:     createdProject.OwnerID,
		CreatedAt:   createdProject.CreatedAt,
	}
	return response, nil
}

func (s *projectService) GetAllProjects(ctx context.Context) (*[]dto.ProjectResponse, error) {
	projects, err := s.projectRepo.ReadAll(ctx)
	if err != nil {
		return nil, err
	}
	var response []dto.ProjectResponse
	for _, project := range *projects {
		response = append(response, dto.ProjectResponse{
			ID:          project.ID,
			Name:        project.Name,
			Description: project.Description,
			OwnerId:     project.OwnerID,
			CreatedAt:   project.CreatedAt,
		})
	}

	return &response, nil
}

func (s *projectService) GetProjectByID(ctx context.Context, projectId *uuid.UUID) (*dto.ProjectResponse, error) {
	project, err := s.projectRepo.ReadById(ctx, projectId)
	if err != nil {
		return nil, err
	}

	response := &dto.ProjectResponse{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		OwnerId:     project.OwnerID,
		CreatedAt:   project.CreatedAt,
	}

	return response, nil
}

func (s *projectService) UpdateProject(ctx context.Context, project *dto.ProjectRequest, projectId *uuid.UUID, ownerId *uuid.UUID) (*dto.ProjectResponse, error) {
	updatedProject, err := s.projectRepo.Update(ctx, project, projectId, ownerId)
	if err != nil {
		return nil, err
	}

	response := &dto.ProjectResponse{
		ID:          updatedProject.ID,
		Name:        updatedProject.Name,
		Description: updatedProject.Description,
		OwnerId:     updatedProject.OwnerID,
		CreatedAt:   updatedProject.CreatedAt,
	}

	return response, nil
}

func (s *projectService) DeleteProject(ctx context.Context, projectId *uuid.UUID, ownerId *uuid.UUID) error {
	err := s.projectRepo.Delete(ctx, projectId, ownerId)
	if err != nil {
		return err
	}
	return nil
}
