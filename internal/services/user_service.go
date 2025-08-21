package services

import (
	"context"
	"errors"

	"github.com/Auxesia23/task_management/internal/auth"
	"github.com/Auxesia23/task_management/internal/dto"
	"github.com/Auxesia23/task_management/internal/repositories"
)

type UserService interface {
	UserRegister(ctx context.Context, user *dto.UserRegister) (*dto.UserResponse, error)
	UserLogin(ctx context.Context, user *dto.UserLogin) (*dto.TokenResponse, error)
	UserRefresh(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error)
	UserSearchByUsername(ctx context.Context, username *string) (*[]dto.UserResponse, error)
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo,
	}
}

func (s *userService) UserRegister(ctx context.Context, user *dto.UserRegister) (*dto.UserResponse, error) {
	if ok := auth.ValidateEmail(user.Email); !ok {
		return nil, errors.New("Invalid email")
	}

	hashed_password, err := auth.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashed_password

	createdUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	response := &dto.UserResponse{
		ID:       createdUser.ID,
		Username: createdUser.Username,
		FullName: createdUser.FullName,
		Email:    createdUser.Email,
	}
	return response, nil
}

func (s *userService) UserLogin(ctx context.Context, userLogin *dto.UserLogin) (*dto.TokenResponse, error) {
	if ok := auth.ValidateEmail(userLogin.Email); !ok {
		return nil, errors.New("Invalid email")
	}

	user, err := s.userRepo.GetByEmail(ctx, &userLogin.Email)
	if err != nil {
		return nil, err
	}

	err = auth.ComparePassword(user.Password, userLogin.Password)
	if err != nil {
		return nil, err
	}

	accesToken, err := auth.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}
	refreshToken, err := auth.GenerateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	tokenResponse := &dto.TokenResponse{
		AccessToken:  accesToken,
		RefreshToken: refreshToken,
	}
	return tokenResponse, nil
}

func (s *userService) UserRefresh(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error) {
	claims, err := auth.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByEmail(ctx, &claims.Email)
	if err != nil {
		return nil, err
	}

	accesToken, err := auth.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}
	tokenResponse := &dto.RefreshResponse{
		AccessToken: accesToken,
	}
	return tokenResponse, nil
}

func (s *userService) UserSearchByUsername(ctx context.Context, username *string) (*[]dto.UserResponse, error) {
	users, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	var response []dto.UserResponse
	for _, user := range *users {
		response = append(response, dto.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			FullName: user.FullName,
			Email:    user.Email,
		})
	}

	return &response, nil
}
