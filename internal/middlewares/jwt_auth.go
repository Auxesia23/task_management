package middlewares

import (
	"strings"

	"github.com/Auxesia23/task_management/internal/auth"
	"github.com/Auxesia23/task_management/internal/dto"
	"github.com/gofiber/fiber/v2"
)

func JWTAuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Message: "Unauthorized",
			Status:  fiber.StatusUnauthorized,
		})
	}
	if ok := strings.HasPrefix(authHeader, "Bearer"); !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Message: "Unauthorized",
			Status:  fiber.StatusUnauthorized,
		})
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Message: "Unauthorized",
			Status:  fiber.StatusUnauthorized,
		})
	}
	claims, err := auth.ValidateAccessToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Message: "Unauthorized",
			Status:  fiber.StatusUnauthorized,
		})
	}
	c.Locals("user", claims)
	return c.Next()
}
