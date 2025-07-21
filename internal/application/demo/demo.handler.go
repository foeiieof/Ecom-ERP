package demo

import (
	"ecommerce/internal/adapter/repository"
	"ecommerce/internal/delivery/http/response"

	"github.com/gofiber/fiber/v2"
)

type DemoHandler struct{
  db *repository.MongoCollectionRepository
}

func NewDemoHandler(repo *repository.MongoCollectionRepository) *DemoHandler {
	return &DemoHandler{db: repo}
}

type DemoResponse struct {
	Status    string `json:"status"`
}

// @Summary      Demo check
// @Description  Check system status
// @Tags         Health 
// @Success      200 {object} DemoResponse
// @Router       /demo [get]
func (d *DemoHandler) DemoCheck(c *fiber.Ctx) error {
  // data,err := d.db.ShopeeAuthCollection.GetShopeeAuth("123")
  // if err != nil {
  //   code := fiber.StatusNotFound
  //   return response.ErrorResponse(c, code,"demo router", err)
  // }
  // data := map[string]string{ "status":"healthy", }
  return response.SuccessResponse(c, "demo router", "")
}
