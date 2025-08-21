package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Auxesia23/task_management/internal/dto"
	"github.com/Auxesia23/task_management/internal/models"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Create(ctx context.Context, conuser *dto.UserRegister) (*models.User, error)
	GetByEmail(ctx context.Context, email *string) (*models.User, error)
	GetByUsername(ctx context.Context, username *string) (*[]models.User, error)
}
type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *dto.UserRegister) (*models.User, error) {
	query := `
		INSERT INTO users (username, full_name, email, password_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING *;
	`

	var createdUser models.User
	if err := r.db.GetContext(ctx, &createdUser, query, user.Username, user.FullName, user.Email, user.Password); err != nil {
		return nil, errors.New("cannot create user")
	}

	return &createdUser, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email *string) (*models.User, error) {
	query := `
		SELECT * FROM users WHERE email = $1;
	`
	var user models.User
	if err := r.db.GetContext(ctx, &user, query, email); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("cannot get user")
	}

	return &user, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username *string) (*[]models.User, error) {
	query := `
			SELECT * FROM users WHERE username ILIKE $1;
		`

	searchTerm := fmt.Sprintf("%%%s%%", *username)

	var user []models.User
	if err := r.db.SelectContext(ctx, &user, query, searchTerm); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("cannot get user")
	}

	return &user, nil
}
