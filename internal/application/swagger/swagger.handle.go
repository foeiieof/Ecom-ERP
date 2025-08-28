package swagger

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger" 
)

type SwaggerHandler interface {
  SwaggerIndex(c *fiber.Ctx) error 
}

type swaggerHandler struct{}

func NewSwaggerHandler() SwaggerHandler {
	return &swaggerHandler{}
}

func (s *swaggerHandler) SwaggerIndex(c *fiber.Ctx) error {
	return swagger.HandlerDefault(c) 
}
