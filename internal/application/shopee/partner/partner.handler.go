package partner

import (
	"ecommerce/internal/delivery/http/response"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type IReqShopeePartnerDTO struct {
	PartnerName string `json:"partner_name"`
	PartnerID   string `json:"partner_id" validate:"required"`
	SecretKey   string `json:"secret_key" validate:"required"`
  Username    *string 
}
// type IReqUpdateShopeePartnerDTO struct { }

type IShopeePartnerHandler interface {
  CreateShopeePartner(c *fiber.Ctx) error
  GetShopeePartnerByID(c *fiber.Ctx) error
  GetAllShopeePartner(c *fiber.Ctx) error
  UpdateShopeePartnerByID(c *fiber.Ctx) error
  DeleteShopeePartnerByID(c *fiber.Ctx) error
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
  // Check param parse

  username := c.Locals("username").(string)
  var reqBody IReqShopeePartnerDTO
  reqBody.Username = &username

  if err := c.BodyParser(&reqBody); err != nil {
    return response.ErrorResponse(c,fiber.StatusBadRequest, "handler.CreateShopeePartner", "invalid body")
  } 

  if err := d.Validate.Struct(&reqBody); err != nil {
    return response.ErrorResponse(c,fiber.StatusBadRequest, "handler.CreateShopeePartner", "invalid body")
  }

  res,err := d.Service.AddShopeePartner(c.Context(),&reqBody)
  if err != nil {
    return response.ErrorResponse(c, fiber.StatusInternalServerError, "handler.CreateShopeePartner",err.Error())
  }

  // d.Logger.Info("handler.CreateShopeePartner" , zap.Any("reqBody", reqBody))
  // Send Service 

  return response.SuccessResponse(c, "handler.parter.CreateShopeePartner" , res)
}

func (d *shopeePartnerHandle)GetShopeePartnerByID(c *fiber.Ctx) error {
  partnerID := c.Params("partnerID")
  if partnerID == "" {
    return response.ErrorResponse(c,fiber.StatusBadRequest, "handler.GetShopeePartner", "partnerID is required")
  }

  res,err := d.Service.GetShopeePartnerByID(c.Context(),partnerID)
  if err != nil {
    return response.ErrorResponse(c,fiber.StatusBadRequest, "handler.GetShopeePartnerByID", err.Error())
  }
  return response.SuccessResponse(c,"handler.GetShopeePartnerByID",res)
}

func (d *shopeePartnerHandle)GetAllShopeePartner(c *fiber.Ctx) error {

  res,err := d.Service.GetAllShopeePartner(c.Context())
  if err != nil {
    return response.ErrorResponse(c,fiber.StatusBadGateway, "handler.GetAllShopeePartner", err.Error())
  }

  return response.SuccessResponse(c,"handler.GetAllShopeePartner", res)
}

func (d *shopeePartnerHandle)UpdateShopeePartnerByID(c *fiber.Ctx) error {
  var update IReqShopeePartnerDTO

  username := c.Locals("username").(string)
  if username != "" { update.Username = &username }

  err := c.BodyParser(&update);if err != nil {
    return response.ErrorResponse(c, fiber.StatusBadRequest , "handler.UpdateShopeePartnerByID", "some params is required")
  }

  if c.Params("partnerID") != update.PartnerID {
    return response.ErrorResponse(c, fiber.StatusBadRequest , "handler.UpdateShopeePartnerByID", "invalidate params")
  }

  err = d.Validate.Struct(&update); if err != nil {
    return response.ErrorResponse(c, fiber.StatusBadRequest, "handler.UpdateShopeePartnerByID", "invalidate params")
  }


  res, err := d.Service.UpdateShopeePartner(c.Context(), &update)
  if err != nil {
    return response.ErrorResponse(c, fiber.StatusBadRequest, "handler.UpdateShopeePartnerByID",err.Error())
  }

  return response.SuccessResponse(c,"handler.UpdateShopeePartnerByID",res)
}

func(d *shopeePartnerHandle)DeleteShopeePartnerByID(c *fiber.Ctx) error {
  partnerID := c.Params("partnerID"); if partnerID == "" {
    return response.ErrorResponse(c,fiber.StatusBadRequest ,"handler.DeleteShopeePartnerByID", "partnerID is required")
  }

  res, err  := d.Service.DeleteShopeePartnerByID(c.Context(), partnerID )
  if err != nil {
    return response.ErrorResponse(c, fiber.StatusBadRequest, "handler.DeleteShopeePartnerByID", err.Error())
  }  

  return response.SuccessResponse(c, "handler.DeleteShopeePartnerByID", res) 
}


