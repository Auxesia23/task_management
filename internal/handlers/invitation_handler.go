package handlers

import (
	"github.com/Auxesia23/task_management/internal/dto"
	"github.com/Auxesia23/task_management/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Invitationhandler interface {
	CreateInvitationHandler(ctx *fiber.Ctx) error
}

type InvitationHandler struct {
	invitationService services.InvitationService
}

func NewInvitationHandler(invitationService services.InvitationService) InvitationHandler {
	return InvitationHandler{
		invitationService,
	}
}

func (h *InvitationHandler) CreateInvitationHandler(c *fiber.Ctx) error {
	inviter := c.Locals("user").(*dto.AccessTokenClaims)
	projectId := c.Params("id")
	userId := c.Params("user_id")

	parsedProjectId, err := uuid.Parse(projectId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	parsedUserId, err := uuid.Parse(userId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if err := h.invitationService.CreateInvitation(c.Context(), &parsedProjectId, &parsedUserId, &inviter.UserID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.Response{
		Status:  fiber.StatusCreated,
		Message: "Invitation created successfully",
	})
}
