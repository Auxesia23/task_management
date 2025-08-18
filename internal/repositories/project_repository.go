package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Auxesia23/task_management/internal/dto"
	"github.com/Auxesia23/task_management/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ProjectRepository interface {
	Create(ctx context.Context, project *dto.ProjectRequest, ownerId *uuid.UUID) (*models.Project, error)
	ReadAll(ctx context.Context) (*[]models.Project, error)
	ReadById(ctx context.Context, projectId *uuid.UUID) (*models.Project, error)
	Update(ctx context.Context, project *dto.ProjectRequest, projectId *uuid.UUID, ownerId *uuid.UUID) (*models.Project, error)
	Delete(ctx context.Context, projectId *uuid.UUID, userId *uuid.UUID) error
}

type projectRepository struct {
	db *sqlx.DB
}

func NewProjectRepository(db *sqlx.DB) ProjectRepository {
	return &projectRepository{
		db: db,
	}
}

func (r *projectRepository) Create(ctx context.Context, project *dto.ProjectRequest, ownerId *uuid.UUID) (*models.Project, error) {
	query := `
		INSERT INTO projects (name,description,owner_id) 
		VALUES ($1,$2,$3) 
		RETURNING *;
	`
	var createdProject models.Project
	if err := r.db.GetContext(ctx, &createdProject, query, project.Name, project.Description, ownerId); err != nil {
		return nil, errors.New("project creation failed")
	}
	return &createdProject, nil
}

func (r *projectRepository) ReadAll(ctx context.Context) (*[]models.Project, error) {
	query := `
			SELECT * FROM projects;
		`
	var projects []models.Project
	if err := r.db.SelectContext(ctx, &projects, query); err != nil {
		return nil, errors.New("project read failed")
	}
	return &projects, nil
}

func (r *projectRepository) ReadById(ctx context.Context, projectId *uuid.UUID) (*models.Project, error) {
	query := `
			SELECT * FROM projects WHERE id = $1;
		`
	var project models.Project
	if err := r.db.GetContext(ctx, &project, query, projectId); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("project not found")
		}
		return nil, errors.New("project read failed")
	}
	return &project, nil
}

func (r *projectRepository) Update(ctx context.Context, project *dto.ProjectRequest, projectId *uuid.UUID, ownerId *uuid.UUID) (*models.Project, error) {
	query := `
			UPDATE projects
			SET name = $1,
				description = $2
			WHERE id = $3 AND owner_id = $4
			RETURNING *;
		`
	var updatedProject models.Project
	err := r.db.GetContext(ctx, &updatedProject, query, project.Name, project.Description, projectId, ownerId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("project not found")
		}
		return nil, err
	}
	return &updatedProject, nil
}

func (r *projectRepository) Delete(ctx context.Context, projectId *uuid.UUID, userId *uuid.UUID) error {
	query := `
			DELETE FROM projects WHERE id = $1 AND owner_id = $2;
		`
	result, err := r.db.ExecContext(ctx, query, projectId, userId)
	if err != nil {
		return errors.New("project delete failed")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.New("project delete failed")
	}
	if rowsAffected == 0 {
		return errors.New("project not found")
	}
	return nil
}
