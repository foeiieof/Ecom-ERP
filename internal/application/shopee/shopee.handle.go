package shopee

import (
	"ecommerce/internal/delivery/http/response"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type IShopeeHandler interface {
  GetShopeeAuthByShopId(c *fiber.Ctx) error
  GetWebHookAuthPartner(c *fiber.Ctx) error
  PostShopAuthPartner(c *fiber.Ctx) error
} 

type shopeeHandler struct {
  service IShopeeService
  logger *zap.Logger
}

func NewShopeeHandler(service IShopeeService, logger *zap.Logger) IShopeeHandler {
  return &shopeeHandler{
    service: service,
    logger: logger,
  }
}

func(d *shopeeHandler) GetShopeeAuthByShopId(c *fiber.Ctx) error {
  // data,err := d.shopeeService.GetAccessToken("123")
  data := c.Locals("shopeeAccessToken")
  // if err != nil {
  //   code := fiber.StatusNotFound
  //   return response.ErrorResponse(c, code,"demo router", err)
  // }
  return response.SuccessResponse(c, "demo router", data)
}

type TPostShopAuthPartner struct {
	PartnerID   string `json:"partner_id"   validate:"required"`
	PartnerKey  string `json:"partner_key"  validate:"required"`
	PartnerName string `json:"partner_name"`

}


type PartnerAuthRequest struct {
}

func(d *shopeeHandler) PostShopAuthPartner(c *fiber.Ctx) error {

  // path := c.Path()
  var reqBody TPostShopAuthPartner
  if err := c.BodyParser(&reqBody); err != nil {
    return response.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body : PostShopAuthPartner", err)
  }


  // Gen Link
  dataLink, _ := d.service.GenerateAuthLink(reqBody.PartnerName,reqBody.PartnerID, reqBody.PartnerKey)

  // Save to DB


  params := map[string]string {"partner_id": reqBody.PartnerID, "partner_key": reqBody.PartnerKey, "partner_name": reqBody.PartnerName, "link" : dataLink}
  data := map[string]any {"Status": "POST", "param": params}

  return response.SuccessResponse(c, "PostShopAuthPartner", &data)  
}

func(d *shopeeHandler) GetWebHookAuthPartner(c *fiber.Ctx) error {

  partnerId := c.Params("partnerId")
  code := c.Query("code")
  shopId := c.Query("shop_id")

  d.logger.Info("Middleware:", zap.String("code", (code) ))
  d.logger.Info("Middleware:", zap.String("shop_id", (shopId) ))

  data, _ := d.service.WebhookAuthentication(partnerId,code, shopId)
  
  return response.SuccessResponse(c, "GetWebHookAuthPartner", data)
}
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
