package demo

import (
	"ecommerce/internal/adapter/repository"
	"ecommerce/internal/delivery/http/response"

	"github.com/gofiber/fiber/v2"
)

type DemoHandler interface {
  DemoCheck(c *fiber.Ctx) error 
}

type demoHandler struct{
  db repository.IMongoCollectionRepository
}

func NewDemoHandler(repo repository.IMongoCollectionRepository) DemoHandler {
	return &demoHandler{db: repo}
}

type DemoResponse struct {
	Status    string `json:"status"`
}

func (d *demoHandler) DemoCheck(c *fiber.Ctx) error {
  // if err != nil {
  //   code := fiber.StatusNotFound
  //   return response.ErrorResponse(c, code,"demo router", err)
  // }
  // data := map[string]string{ "status":"healthy", }
  return response.SuccessResponse(c, "demo router", "")
}
