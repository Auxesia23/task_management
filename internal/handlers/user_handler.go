package handlers

import (
	"github.com/Auxesia23/task_management/internal/dto"
	"github.com/Auxesia23/task_management/internal/services"
	"github.com/gofiber/fiber/v2"
)

type UserHandler interface {
	RegisterHandler(c *fiber.Ctx) error
	LoginHandler(c *fiber.Ctx) error
	RefreshHandler(c *fiber.Ctx) error
	SearchUserhandler(c *fiber.Ctx) error
}

type userHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) UserHandler {
	return &userHandler{
		userService,
	}
}

func (h *userHandler) RegisterHandler(c *fiber.Ctx) error {
	var input dto.UserRegister
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	createdUser, err := h.userService.UserRegister(c.Context(), &input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.Response{
		Status:  fiber.StatusCreated,
		Message: "User registered successfully",
		Data:    createdUser,
	})
}

func (h *userHandler) LoginHandler(c *fiber.Ctx) error {
	var input dto.UserLogin
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	token, err := h.userService.UserLogin(c.Context(), &input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.Response{
		Status:  fiber.StatusOK,
		Message: "User logged in successfully",
		Data:    token,
	})
}

func (h *userHandler) RefreshHandler(c *fiber.Ctx) error {
	var input dto.RefreshRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	token, err := h.userService.UserRefresh(c.Context(), input.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.Response{
		Status:  fiber.StatusOK,
		Message: "Token refreshed successfully",
		Data:    token,
	})
}

func (h *userHandler) SearchUserhandler(c *fiber.Ctx) error {
	q := c.Query("q")
	if q == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Query parameter is required",
		})
	}

	users, err := h.userService.UserSearchByUsername(c.Context(), &q)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
			Status:  fiber.StatusNotFound,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.Response{
		Status:  fiber.StatusOK,
		Message: "User search successfully",
		Data:    users,
	})
}
