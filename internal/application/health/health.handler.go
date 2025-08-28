package health

import (
	"ecommerce/internal/delivery/http/response"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type HealthHandler interface {
  HealthCheck(c *fiber.Ctx) error
}

type healthHandler struct {
  logger *zap.Logger
}

func NewHealthHandler(logger *zap.Logger) HealthHandler {
  return &healthHandler{
    logger: logger,
  }
}

// @Summary      Health check
// @Description  Check system status
// @Tags         Health
// @Success      200 {object} map[string]string
// @Router       /api/v1/health [get]
func (h *healthHandler) HealthCheck(c *fiber.Ctx) error {
	return  response.SuccessResponse(c,"ok!","")
  // c.JSON(fiber.Map{ "status":    "healthy", "timestamp": fmt.Sprintf("%d", c.Context().Time().Unix()), })
}
