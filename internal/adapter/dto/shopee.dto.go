package dto

import (
	"encoding/json"
)

// type IGenerateBodyQueryParams struct {
//   Query string 
//   Body *bytes.Buffer
// }

type IEnumShopeeTimeRange string
const (
  CREATE_TIME IEnumShopeeTimeRange = "create_time"
  UPDATE_TIME IEnumShopeeTimeRange = "update_time"
)

// UNPAID/READY_TO_SHIP/PROCESSED/SHIPPED/COMPLETED/IN_CANCEL/CANCELLED/INVOICE_PENDING
type IEnumShopeeOrderStatus string
const (
  UNPAID IEnumShopeeOrderStatus = "UNPAID"
  READY_TO_SHIP IEnumShopeeOrderStatus = "READY_TO_SHIP"
  PROCESSED IEnumShopeeOrderStatus = "PROCESSED"
  SHIPPED IEnumShopeeOrderStatus = "SHIPPED"
  COMPLETED IEnumShopeeOrderStatus = "COMPLETED"
  IN_CANCEL IEnumShopeeOrderStatus = "IN_CANCEL"
  CANCELLED IEnumShopeeOrderStatus = "CANCELLED"
  INVOICE_PENDING IEnumShopeeOrderStatus = "INVOICE_PENDING"
)

type IEnumShopeeOptionsFields string
const (
  OrderStatus IEnumShopeeOptionsFields = "order_status"
)

type IOptionShopeeQuery struct {
  TimeRange IEnumShopeeTimeRange
  TimeFrom    int64 //Unix
  TimeTo      int64 //Unix
  PageSize    int32 //Page 

  CursorPage  string //Base64
  OrderStatus IEnumShopeeOrderStatus
  ResponseOptionsField IEnumShopeeOptionsFields // Available value: order_status.
  RequestOrderStatus bool
  LogisticsChanelID  string //

}

type IResShopeeResponse struct {
	ReqquestID string `json:"request_id"`
	Error      string `json:"error"`
	Message    string `json:"message"`
  More       bool   `json:"more"`
}

type IResShopeeAuthRefreshResponse struct {
	IResShopeeResponse
  PartnerID    json.Number `json:"partner_id"`
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	ExpireIn     json.Number `json:"expire_in"`
	ShopID       json.Number `json:"shop_id"`
}

type IBGetRefreshToken struct {
	RefreshToken string `json:"refresh_token"`
	PartnerID    int32  `json:"partner_id"`
	ShopID       int32  `json:"shop_id"`
}

type IResSipAffiShopList struct {
	Region     string      `json:"region"`
	AffiShopID json.Number `json:"affi_shop_id"`
}

type IResAuthedShopList struct {
	ShopID          json.Number           `json:"shop_id"`
	ExpireTime      json.Number           `json:"expire_time"`
	AuthTime        json.Number           `json:"auth_time"`
	Region          string                `json:"region"`
	SipAffiShopList []IResSipAffiShopList `json:"sip_affi_shop_list"`
}

type IResGetShopByPartnerPublic struct {
  IResShopeeResponse
  AuthedShopList []IResAuthedShopList `json:"authed_shop_list"`
}

// old
// type IResOrderList struct {
// 	OrderSN     string `json:"order_sn"`
// 	OrderStatus string `json:"order_status"`
// 	BookingSN   string `json:"booking_sn"`
// }

// type IResGetOrderListByShopIDShop struct {
// 	IResShopeeResponse
//   OrderList []IResOrderList `json:"order_list"`
// }
type IResOrderList struct {
	OrderSN   string `json:"order_sn"`
	BookingSN string `json:"booking_sn"`
	// OrderStatus is optional and not present in current JSON
	OrderStatus string `json:"order_status,omitempty"`
}

type IResGetOrderListByShopIDShopWrapper struct {
	More       bool             `json:"more"`
	NextCursor string           `json:"next_cursor"`
	OrderList  []IResOrderList  `json:"order_list"`
}

type IResGetOrderListByShopIDShop struct {
	IResShopeeResponse         // embeds request_id, error, message
	Response IResGetOrderListByShopIDShopWrapper `json:"response"`
}

// Reminder Note :
// * Embeded all struct
