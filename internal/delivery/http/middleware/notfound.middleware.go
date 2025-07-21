package middleware

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type ErrorHandler struct {
	logger *zap.Logger
}

func NewErrorHandler(log *zap.Logger) *ErrorHandler {
	return &ErrorHandler{logger: log}
}

func (e *ErrorHandler) NotFoundHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// e.logger.Error("Route not found", zap.String("path", c.Path()), zap.String("method", c.Method()))
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":  "Route not found",
			"path":   c.OriginalURL(),
			"method": c.Method(),
			"ip":     c.IP(),
		})
	}
}
