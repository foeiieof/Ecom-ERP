package partner

import (
	"ecommerce/internal/delivery/http/response"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type IReqShopeePartnerDTO struct {
	PartnerID   string `json:"partner_id"   validate:"required"`
	PartnerName string `json:"partner_name"`
	PartnerKey  string `json:"secret_key"  validate:"required"`
}

type IShopeePartnerHandler interface {
  CreateShopeePartner(c *fiber.Ctx) error
}

type shopeePartnerHandle struct {
  Logger *zap.Logger
  Validate *validator.Validate
  Service IShopeePartnerService
}

func NewShopeePartnerHandler(log *zap.Logger, valid *validator.Validate, srv IShopeePartnerService) IShopeePartnerHandler{
  return &shopeePartnerHandle{
    Logger: log,
    Validate: valid,
    Service: srv,
  }
}

func (d *shopeePartnerHandle)CreateShopeePartner(c *fiber.Ctx) error {
  return response.SuccessResponse(c, "handler.parter.CreateShopeePartner" , "")
}


