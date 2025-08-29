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

type Account struct {
  demo string
}
// is available 
//  @Summary      List accounts
//  @Description  get accounts
//  @Tags         accounts
//  @Accept       json
//  @Produce      json
//  @Param        q    query     string  false  "name search by q"  Format(email)
//  @Success      200  {array}   response.APIResponse[Account]
//  @Router       /health [get]
func (h *healthHandler) HealthCheck(c *fiber.Ctx) error {
	return  response.SuccessResponse(c,"ok!","")
  // c.JSON(fiber.Map{ "status":    "healthy", "timestamp": fmt.Sprintf("%d", c.Context().Time().Unix()), })
}
