package shopee

import (
	"ecommerce/internal/delivery/http/response"
  "ecommerce/internal/application/shopee/partner"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// IShopeeHandler: < IShopeeService
type IShopeeHandler interface {
	// GetShopeeAuthByShopId(c *fiber.Ctx) error
	GetWebHookAuthPartner(c *fiber.Ctx) error

	GetShopeeTokenAuthPartnerByShopId(c *fiber.Ctx) error

	PostShopAuthPartner(c *fiber.Ctx) error
	PostShopeeTokenAuthPartnerWithCode(c *fiber.Ctx) error

	PostShopeeDemoTemplate(c *fiber.Ctx) error

  //
  GetShopeeShopDetails(c *fiber.Ctx) error
	// Partner IShopeeHandler
	GetShopeeShopListByPartnerID(c *fiber.Ctx) error

	// Order
	GetShopeeOrderListByShopID(c *fiber.Ctx) error
	GetShopeeOrderDetailsByShopIDAndOrderSN(c *fiber.Ctx) error
}

type shopeeHandler struct {
	ShopeeService IShopeeService
  PartnerService partner.IShopeePartnerService
	Logger  *zap.Logger
	Valid   *validator.Validate
}

func NewShopeeHandler(service IShopeeService, partner partner.IShopeePartnerService,logger *zap.Logger, valid *validator.Validate) IShopeeHandler {
	return &shopeeHandler{
		ShopeeService: service,
    PartnerService: partner,
    Logger:  logger,
		Valid:   valid,
	}
}

// func (d *shopeeHandler) GetShopeeAuthByShopId(c *fiber.Ctx) error {
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
// return response.SuccessResponse(c, "demo router", "")
// }

type TPostShopAuthPartner struct {
	PartnerID   string `json:"partner_id"   validate:"required"`
	PartnerKey  string `json:"partner_key"  validate:"required"`
	PartnerName string `json:"partner_name"`
}

// type PartnerAuthRequest struct { }

func (d *shopeeHandler) PostShopAuthPartner(c *fiber.Ctx) error {
	// path := c.Path()
	var reqBody partner.IReqShopeePartnerDTO 
	var err error
	if err = c.BodyParser(&reqBody); err != nil {
		return response.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body : PostShopAuthPartner", err)
	}

	// Gen Link
	dataLink, err := d.ShopeeService.GenerateAuthLink(c.Context(),reqBody.PartnerName, reqBody.PartnerID, reqBody.SecretKey)
	if err != nil {
		d.Logger.Error("service.GenerateAuthLink :", zap.Error(err))
		return response.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body : PostShopAuthPartner", err)
	}

	// Save log request to DB
	_, err = d.ShopeeService.AddShopeeAuthRequest(c.Context(),reqBody.PartnerID, reqBody.SecretKey, reqBody.PartnerName, dataLink)
	if err != nil {
		d.Logger.Error("service.AddShopeeAuthRequest :", zap.Error(err))
		return response.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body : PostShopAuthPartner", err)
	}


  // partnerDTO := partner.ShopeePartnerDTO()
	_, err = d.PartnerService.AddShopeePartner(c.Context(),&reqBody)
	if err != nil {
		d.Logger.Error("service.AddShopeePartner :", zap.Error(err))
		return response.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body : PostShopAuthPartner", err)
	}

	params := map[string]string{"partner_id": reqBody.PartnerID, "partner_key": reqBody.SecretKey, "partner_name": reqBody.PartnerName, "link": dataLink}

	data := map[string]any{"Status": "POST", "param": params}

	return response.SuccessResponse(c, "PostShopAuthPartner", &data)
}

func (d *shopeeHandler) GetWebHookAuthPartner(c *fiber.Ctx) error {

	partnerId := c.Params("partnerId")
	code := c.Query("code")
	shopId := c.Query("shop_id")

  data, err := d.ShopeeService.WebhookAuthentication(c.Context() , partnerId, code, shopId)
  if err != nil { return response.ErrorResponse(c, fiber.StatusConflict, "shopee.handler.GetWebHookAuthPartner", err.Error())}

	return response.SuccessResponse(c, "GetWebHookAuthPartner", data)
}

type ReqShopeeTokenAuthPartner struct {
	PartnerID string `json:"partner_id" validate:"required"`
	Code      string `json:"code"       validate:"required"`
	ShopID    string `json:"shop_id"    validate:"required"`
}

func (d *shopeeHandler) PostShopeeTokenAuthPartnerWithCode(c *fiber.Ctx) error {

	var reqBody ReqShopeeTokenAuthPartner
	if err := c.BodyParser(&reqBody); err != nil {
		d.Logger.Error("handle.PostShopeeTokenAuthPartner : c.BodyParser(&reqBody) :", zap.Error(err))
		return response.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body : PostShopeeTokenAuthPartner", err)
	}

	if err := d.Valid.Struct(reqBody); err != nil {
		d.Logger.Error("handle.PostShopeeTokenAuthPartner : vilid.Struct(&reqBody) :", zap.Error(err))
		return response.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body : PostShopeeTokenAuthPartner", err)
	}

	// Generate sign
	dataGen, err := d.ShopeeService.CreateAccessAndRefreshTokenByCodeOnAdapter(c.Context(),reqBody.PartnerID, reqBody.ShopID, reqBody.Code)

	if err != nil {
		d.Logger.Error("handle.PostShopeeTokenAuthPartner : d.service.GetAccessAndRefreshToken :", zap.Error(err))
		return response.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body : PostShopeeTokenAuthPartner", err)
	}

	// ShopeeService
	return response.SuccessResponse(c, "PostShopeeTokenAuthPartner", dataGen)
}

func (d *shopeeHandler) GetShopeeTokenAuthPartnerByShopId(c *fiber.Ctx) error {
	// data,err := d.shopeeService.GetAccessToken("123")
	shopID := c.Params("shopeeShopID")
	if shopID == "" {
		d.Logger.Error("handle.GetShopeeTokenAuthPartnerByShopId:", zap.String("shopId", ""))
		return response.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body : GetShopeeTokenAuthPartnerByShopId", "shopId is required")
	}

	// d.logger.Debug("handle.GetShopeeTokenAuthPartnerByShopId", zap.String("shopId", shopID))

	data, err := d.ShopeeService.GetAccessTokenByShopID(c.Context(),shopID)
	if err != nil {
		d.Logger.Error("handle.GetShopeeTokenAuthPartnerByShopId : d.service.GetAccessToken :", zap.Error(err))
		return response.ErrorResponse(c, fiber.StatusNotFound, "ShopId no found", err.Error())
	}

	return response.SuccessResponse(c, "shopee router", data)
}

func (d *shopeeHandler) GetShopeeShopListByPartnerID(c *fiber.Ctx) error {
	partnerID := c.Params("partnerID")

	data, err := d.ShopeeService.GetShopeeShopListByPartnerID(c.Context(),partnerID)
	if err != nil {
		d.Logger.Error("handle.GetShopeeShopListByPartnerID : d.service.GetShopeeShopListByPartnerID :", zap.Error(err))
		return response.ErrorResponse(c, fiber.StatusNotFound, "ShopId no found", err.Error())
	}

  // 

	return response.SuccessResponse(c, "GetShopeeShopListByPartnerID", data)
}

type IReqQueryShopeeOrderListByShopID struct {
	// ShopID string `json:"shopeeShopID" validate:"required"`
	From   string `json:"from"    query:"from"   validate:"required"`
	To     string `json:"to"      query:"to"     validate:"required"`
	Page   string `json:"page"    query:"page"   `
	Size   string `json:"size"    query:"size"   validate:"required"`
	Status string `json:"status"  query:"status" ` // order_status
	Type   string `json:"type"    query:"type"   ` // creation_time, update_time
}

func (d *shopeeHandler) GetShopeeOrderListByShopID(c *fiber.Ctx) error {
	// Params
	shopID := c.Params("shopeeShopID")
	if shopID == "" {
		return response.ErrorResponse(c, fiber.StatusBadRequest, "shopeeHandle.GetShopeeOrderListByShopID", "shopId is required")
	}

	// Querys
	typeQuery := c.Query("type") // OrderType
	timeFromQuery := c.Query("from")
	timeToQuery := c.Query("to")
	statusQuery := c.Query("status")
	nextQuery := c.Query("page") // Cursor
	sizeQuery := c.Query("size") // PageSize

	// d.logger.Debug("time start day in unix", zap.String("timstamp", strconv.FormatInt( time.Now().Truncate(24*time.Hour).Unix(),10) ) )
	// d.logger.Debug("time end day in unix", zap.String("timstamp", strconv.FormatInt( time.Now().Truncate(24*time.Hour).Add(23 * time.Hour + 59* time.Minute + 59*time.Second).Unix(),10) ) )

	// check access and refresh
	// _, err := d.ShopeeService.GetAccessTokenByShopID(c.Context(),shopID)
	// if err != nil {
	// 	d.Logger.Error("handle.GetShopeeOrderListByShopID : d.service.GetAccessToken :", zap.Error(err))
	// 	return response.ErrorResponse(c, fiber.StatusNotFound, "ShopId no found", err.Error())
	// }

	// Valid section
	var queries IReqQueryShopeeOrderListByShopID
	if err := c.QueryParser(&queries); err != nil {
		d.Logger.Error("shopeeHandle.GetShopeeOrderListByShopID.c.QueryParser", zap.Error(err))
		return response.ErrorResponse(c, fiber.StatusBadRequest, "shopeeHandle.GetShopeeOrderListByShopID.c.QueryParser", "Invalid request body")
	}
	if err := d.Valid.Struct(queries); err != nil {
		d.Logger.Error("shopeeHandle.GetShopeeOrderListByShopID.queries", zap.Error(err))
		return response.ErrorResponse(c, fiber.StatusBadRequest, "shopeeHandle.GetShopeeOrderListByShopID.queries", "Invalid request body")
	}

	data, err := d.ShopeeService.GetShopeeOrderListByShopID(c.Context(),shopID, typeQuery, timeFromQuery, timeToQuery, statusQuery, nextQuery, sizeQuery)
	if err != nil {
		d.Logger.Error("handle.GetShopeeOrderListByShopID : d.service.GetShopeeOrderListByShopID :", zap.Error(err))
		return response.ErrorResponse(c, fiber.StatusNotFound, "usecase.GetShopeeOrderListByShopID :", err.Error())
	}
	// d.Logger.Debug("shopeeHandle.GetShopeeOrderListByShopID", zap.Any("data", data))

	// GetOrderListByShopID
	return response.SuccessResponse(c, "shopeeHandle.GetShopeeOrderListByShopID", data.OrderList)
}

func (d *shopeeHandler)GetShopeeShopDetails(c *fiber.Ctx) error {
  shopID := c.Params("shopeeShopID")
  userName,ok := c.Locals("username").(string)
  if !ok {
   return response.ErrorResponse(c,fiber.StatusUnauthorized, "handler.GetShopDetails", "unaurthorize")
  }

  // d.Logger.Debug("handler.GetShopeeShopDetailsByShopID", zap.String("userName", userName))

  res,err := d.ShopeeService.GetShopeeShopDetailsByShopID(c.Context(), userName ,shopID)
  if err != nil { return response.ErrorResponse(c, fiber.StatusConflict,"handler.GetShopeeShopDetails" ,err) }

  return response.SuccessResponse(c, "handle.GetShopeeShopDetails", res)
}

func (d *shopeeHandler) GetShopeeOrderDetailsByShopIDAndOrderSN(c *fiber.Ctx) error {

	shopIDParam := c.Params("shopeeShopID")
	orderSNParam := c.Params("orderSN")

	missing := []string{}

  if shopIDParam == "" { missing = append(missing, "shopID")    } 
	if orderSNParam == "" { missing = append(missing, "orderSN") }
  if len(missing) > 0 {
		return response.ErrorResponse(c, fiber.StatusBadRequest, "shopeeHandle.GetShopeeOrderListByShopSN", fmt.Sprintf("%s is required", missing[0]) )
  } 

	pendingQuery := c.Query("pending")
	optionQuery := c.Query("option")

	data, err := d.ShopeeService.GetShopeeOrderDetailByOrderSN(c.Context(),shopIDParam, orderSNParam, pendingQuery, optionQuery)
	if err != nil {
		d.Logger.Error("handle.GetShopeeOrderListByShopSN : d.service.GetShopeeOrderDetailByShopID :", zap.Error(err))
		return response.ErrorResponse(c, fiber.StatusNotFound, "usecase.GetShopeeOrderDetailByShopID :", err.Error())
	}

	// d.Logger.Debug("shopeeHandle.GetShopeeOrderListByShopSN", zap.Any("data", data.OrderList))

	// res := orderParams + pendingQuery + optionQuery

  // Test 
  // return response.SuccessResponse(c, "shopeeHandle.GetShopeeOrderListByShopSN", fmt.Sprintf("%s*-*%s", shopIDParam,orderSNParam ))
	return response.SuccessResponse(c, "shopeeHandle.GetShopeeOrderListByShopSN", data.OrderList)
}

// ------------------------------------------------- Template -------------------------------------------------------
// reqInterface  Template
type IReqShopeeDemoTemplate struct {
	PartnerID string `json:"partner_id" validate:"required"`
	Code      string `json:"code"       validate:"required"`
	ShopID    string `json:"shop_id"    validate:"required"`
}

// Template
func (d *shopeeHandler) PostShopeeDemoTemplate(c *fiber.Ctx) error {
	var reqBody IReqShopeeDemoTemplate
	if err := c.BodyParser(&reqBody); err != nil {
		return response.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body : PostShopeeDemoTemplate", err)
	}
	if err := d.Valid.Struct(reqBody); err != nil {
		return response.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body : PostShopeeDemoTemplate", err)
	}
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
