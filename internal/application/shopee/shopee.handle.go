package shopee

import (
	"ecommerce/internal/delivery/http/response" 
  "github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)
// IShopeeHandler: < IShopeeService
type IShopeeHandler interface {
	GetShopeeAuthByShopId(c *fiber.Ctx) error
	GetWebHookAuthPartner(c *fiber.Ctx) error

  GetShopeeTokenAuthPartnerByShopId(c *fiber.Ctx) error
	
  PostShopAuthPartner(c *fiber.Ctx) error
  PostShopeeTokenAuthPartner(c *fiber.Ctx) error

  PostShopeeDemoTemplate(c *fiber.Ctx) error
}

type shopeeHandler struct {
	service IShopeeService
	logger  *zap.Logger
  valid   *validator.Validate
}

func NewShopeeHandler(service IShopeeService, logger *zap.Logger, valid *validator.Validate) IShopeeHandler {
	return &shopeeHandler{
		service: service,
		logger:  logger,
    valid:   valid,
	}
}

func (d *shopeeHandler) GetShopeeAuthByShopId(c *fiber.Ctx) error {
	// data,err := d.shopeeService.GetAccessToken("123")
	// shopID := c.Params("shopeeShopID")
 //  data,err := d.service.GetAccessToken(shopID)
 //  if err != nil {
 //    return response.ErrorResponse(c, fiber.StatusBadRequest, "", err)
 //  }
	// // if err != nil {
	// //   code := fiber.StatusNotFound return response.ErrorResponse(c, code,"demo router", err)
	// // }
	// return response.SuccessResponse(c, "demo router", data)
  return response.SuccessResponse(c, "demo router", "")
}

type TPostShopAuthPartner struct {
	PartnerID   string `json:"partner_id"   validate:"required"`
	PartnerKey  string `json:"partner_key"  validate:"required"`
	PartnerName string `json:"partner_name"`
}

// type PartnerAuthRequest struct { }

func (d *shopeeHandler) PostShopAuthPartner(c *fiber.Ctx) error {
	// path := c.Path()
	var reqBody TPostShopAuthPartner
  var err error
	if err = c.BodyParser(&reqBody); err != nil {
		return response.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body : PostShopAuthPartner", err)
	}

	// Gen Link
	dataLink, err := d.service.GenerateAuthLink(reqBody.PartnerName, reqBody.PartnerID, reqBody.PartnerKey)
  if err != nil {
    d.logger.Error("service.GenerateAuthLink :", zap.Error(err))
    return response.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body : PostShopAuthPartner", err)
  }

	// Save log request to DB
	_, err = d.service.AddShopeeAuthRequest(reqBody.PartnerID, reqBody.PartnerKey, reqBody.PartnerName, dataLink)
	if err != nil {
    d.logger.Error("service.AddShopeeAuthRequest :", zap.Error(err))
		return response.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body : PostShopAuthPartner", err)
	}

  _, err = d.service.AddShopeePartner(reqBody.PartnerID, reqBody.PartnerKey, reqBody.PartnerName) 
  if err != nil {
    d.logger.Error("service.AddShopeePartner :", zap.Error(err))
    return response.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body : PostShopAuthPartner", err)
  }

  params := map[string]string{"partner_id": reqBody.PartnerID, "partner_key": reqBody.PartnerKey, "partner_name": reqBody.PartnerName, "link": dataLink}

	data := map[string]any{"Status": "POST", "param": params}

	return response.SuccessResponse(c, "PostShopAuthPartner", &data)
}

func (d *shopeeHandler) GetWebHookAuthPartner(c *fiber.Ctx) error {

	partnerId := c.Params("partnerId")
	code := c.Query("code")
	shopId := c.Query("shop_id")

	d.logger.Info("Middleware:", zap.String("code", (code)))
	d.logger.Info("Middleware:", zap.String("shop_id", (shopId)))

	data, _ := d.service.WebhookAuthentication(partnerId, code, shopId)

	return response.SuccessResponse(c, "GetWebHookAuthPartner", data)
}


type ReqShopeeTokenAuthPartner struct {
  PartnerID  string `json:"partner_id" validate:"required"`
  Code       string `json:"code"       validate:"required"`
  ShopID     string `json:"shop_id"    validate:"required"`
}

func (d *shopeeHandler) PostShopeeTokenAuthPartner(c *fiber.Ctx) error {
  
  var reqBody ReqShopeeTokenAuthPartner
  if err := c.BodyParser(&reqBody); err != nil {
    d.logger.Error("handle.PostShopeeTokenAuthPartner : c.BodyParser(&reqBody) :", zap.Error(err))
    return response.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body : PostShopeeTokenAuthPartner", err)
  }
  
  if err := d.valid.Struct(reqBody); err != nil {
    d.logger.Error("handle.PostShopeeTokenAuthPartner : vilid.Struct(&reqBody) :", zap.Error(err))
    return response.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body : PostShopeeTokenAuthPartner", err)
  }

  // Generate sign
  dataGen, err := d.service.GetAccessAndRefreshToken(reqBody.PartnerID, reqBody.ShopID, reqBody.Code)

  if err != nil {
    d.logger.Error("handle.PostShopeeTokenAuthPartner : d.service.GetAccessAndRefreshToken :", zap.Error(err))
    return response.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body : PostShopeeTokenAuthPartner", err)
  }

  // ShopeeService 
  return response.SuccessResponse(c, "PostShopeeTokenAuthPartner", dataGen)
}

func (d *shopeeHandler) GetShopeeTokenAuthPartnerByShopId(c *fiber.Ctx) error {
  // data,err := d.shopeeService.GetAccessToken("123")
  shopID := c.Params("shopeeShopID")
  data,err := d.service.GetAccessToken(shopID)
  if err != nil {
    d.logger.Error("handle.GetShopeeTokenAuthPartnerByShopId : d.service.GetAccessToken :", zap.Error(err))
    return response.ErrorResponse(c, fiber.StatusNotFound, "ShopId no found", err.Error()) 
  }

  return response.SuccessResponse(c, "shopee router", data)
}


// ------------------------------------------------- Template -------------------------------------------------------
// reqInterface  Template
type ReqShopeeDemoTemplate struct {
  PartnerID  string `json:"partner_id" validate:"required"`
  Code       string `json:"code"       validate:"required"`
  ShopID     string `json:"shop_id"    validate:"required"`
} 
// Template
func (d *shopeeHandler) PostShopeeDemoTemplate(c *fiber.Ctx) error {
  var reqBody ReqShopeeDemoTemplate
  if err := c.BodyParser(&reqBody); err != nil {
    return response.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body : PostShopeeDemoTemplate", err) }
  if err := d.valid.Struct(reqBody); err != nil {
    return response.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body : PostShopeeDemoTemplate", err) }
  return response.SuccessResponse(c, "PostShopeeTokenAuthPartner", reqBody)
}
// ------------------------------------------------- End - Template --------------------------------------------------


// Prototype
// func(d *shopeeHandler) PostShopAuthPartner(c *fiber.Ctx) error {
//   // path := c.Path()
//   var reqBody TPostShopAuthPartner
//   if err := c.BodyParser(&reqBody); err != nil {
//     return response.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body : PostShopAuthPartner", err)
//   }

//   data := map[string]string{"Status": "POST", "param": reqBody.ShopID}
//   return response.SuccessResponse(c, "PostShopAuthPartner", &data)
// }
