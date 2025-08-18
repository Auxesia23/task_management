package handlers

import (
	"github.com/Auxesia23/task_management/internal/dto"
	"github.com/Auxesia23/task_management/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProjectHandler interface {
	CreateProjectHanlder(c *fiber.Ctx) error
	GetProjectsHanlder(c *fiber.Ctx) error
	ReadProjectByIdHanlder(c *fiber.Ctx) error
	UpdateProjectHanlder(c *fiber.Ctx) error
	DeleteProjectHanlder(c *fiber.Ctx) error
}

type projectHandler struct {
	projectServive services.ProjectService
}

func NewProjectHandler(projectServive services.ProjectService) ProjectHandler {
	return &projectHandler{projectServive}
}

func (h *projectHandler) CreateProjectHanlder(c *fiber.Ctx) error {
	user := c.Locals("user").(*dto.AccessTokenClaims)

	var input dto.ProjectRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Message: err.Error(),
			Status:  fiber.StatusBadRequest,
		})
	}

	createdProject, err := h.projectServive.CreateProject(c.Context(), &input, &user.UserID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Message: err.Error(),
			Status:  fiber.StatusBadRequest,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.Response{
		Message: "Project created successfully",
		Status:  fiber.StatusCreated,
		Data:    createdProject,
	})
}

func (h *projectHandler) GetProjectsHanlder(c *fiber.Ctx) error {
	projects, err := h.projectServive.GetAllProjects(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Message: err.Error(),
			Status:  fiber.StatusInternalServerError,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.Response{
		Message: "Projects retrieved successfully",
		Status:  fiber.StatusOK,
		Data:    projects,
	})
}

func (h *projectHandler) ReadProjectByIdHanlder(c *fiber.Ctx) error {
	id := c.Params("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Message: err.Error(),
			Status:  fiber.StatusBadRequest,
		})
	}

	project, err := h.projectServive.GetProjectByID(c.Context(), &parsedId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
			Message: err.Error(),
			Status:  fiber.StatusNotFound,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.Response{
		Message: "Project retrieved successfully",
		Status:  fiber.StatusOK,
		Data:    project,
	})
}

func (h *projectHandler) UpdateProjectHanlder(c *fiber.Ctx) error {
	user := c.Locals("user").(*dto.AccessTokenClaims)
	id := c.Params("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Message: err.Error(),
			Status:  fiber.StatusBadRequest,
		})
	}
	var input dto.ProjectRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Message: err.Error(),
			Status:  fiber.StatusBadRequest,
		})
	}

	updatedProject, err := h.projectServive.UpdateProject(c.Context(), &input, &parsedId, &user.UserID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
			Message: err.Error(),
			Status:  fiber.StatusNotFound,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.Response{
		Message: "Project updated successfully",
		Status:  fiber.StatusOK,
		Data:    updatedProject,
	})
}

func (h *projectHandler) DeleteProjectHanlder(c *fiber.Ctx) error {
	user := c.Locals("user").(*dto.AccessTokenClaims)
	id := c.Params("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Message: err.Error(),
			Status:  fiber.StatusBadRequest,
		})
	}

	err = h.projectServive.DeleteProject(c.Context(), &parsedId, &user.UserID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
			Message: err.Error(),
			Status:  fiber.StatusNotFound,
		})
	}
	return c.Status(fiber.StatusNoContent).JSON(dto.Response{
		Message: "Project deleted successfully",
		Status:  fiber.StatusNoContent,
	})
}
