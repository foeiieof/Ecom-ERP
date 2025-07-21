package swagger

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger" 
)

type SwaggerHandler struct{}

func NewSwaggerHandler() *SwaggerHandler {
	return &SwaggerHandler{}
}

func (s *SwaggerHandler) SwaggerIndex(c *fiber.Ctx) error {
	return swagger.HandlerDefault(c) 
}
