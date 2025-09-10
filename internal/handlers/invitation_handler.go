package handlers

import (
	"github.com/Auxesia23/task_management/internal/dto"
	"github.com/Auxesia23/task_management/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Invitationhandler interface {
	CreateInvitationHandler(ctx *fiber.Ctx) error
	GetInvitationHandler(ctx *fiber.Ctx) error
	UpdateInvitationHandler(ctx *fiber.Ctx) error
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

	var input dto.InvitationRequest
	err := c.BodyParser(&input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}
	parsedProjectId, err := uuid.Parse(projectId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if err := h.invitationService.CreateInvitation(c.Context(), &parsedProjectId, &input.UserID, &inviter.UserID); err != nil {
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

func (h *InvitationHandler) GetInvitationHandler(c *fiber.Ctx) error {
	user := c.Locals("user").(*dto.AccessTokenClaims)
	status := c.Query("status", "pending")

	invitations, err := h.invitationService.GetInvitation(c.Context(), &user.UserID, &status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.Response{
		Status:  fiber.StatusOK,
		Message: "Invitations get successfully",
		Data:    invitations,
	})
}

func (h *InvitationHandler) UpdateInvitationHandler(c *fiber.Ctx) error {
	user := c.Locals("user").(*dto.AccessTokenClaims)
	status := c.Query("status")
	invitationId := c.Params("id")
	parsedInvitationId, err := uuid.Parse(invitationId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if err := h.invitationService.UpdateInvitation(c.Context(), &user.UserID, &parsedInvitationId, &status); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.Response{
		Status:  fiber.StatusOK,
		Message: "Invitation accepted successfully",
	})
}
