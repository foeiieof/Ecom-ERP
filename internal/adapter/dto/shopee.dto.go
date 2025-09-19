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
  TotalAmount IEnumShopeeOptionsFields = "total_amount"
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
  More       *bool   `json:"more"`
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

type IResOrderDetailReceiptAddress struct {
	Name         string `json:"name"`
	Phone        string `json:"phone"`
	Town         string `json:"town"`
	District     string `json:"district"`
	City         string `json:"city"`
	State        string `json:"state"`
	Region       string `json:"region"`
	Zipcode      string `json:"zipcode"`
	FullAddress  string `json:"full_address"`
	VirtualPhone string `json:"virtual_contact_number"`
}

type IResOrderDetailImage struct {
	ImageURL string `json:"image_url"`
}

type IResOrderDetailItem struct {
	ItemID               int64   `json:"item_id"`
	ItemName             string  `json:"item_name"`
	ItemSKU              string  `json:"item_sku"`
	ModelID              int64   `json:"model_id"`
	ModelName            string  `json:"model_name"`
	ModelSKU             string  `json:"model_sku"`
	ModelQtyPurchased    int     `json:"model_quantity_purchased"`
	ModelOriginalPrice   float64 `json:"model_original_price"`
	ModelDiscountedPrice float64 `json:"model_discounted_price"`
	Wholesale            bool    `json:"wholesale"`
	Weight               float64 `json:"weight"`
	AddOnDeal            bool    `json:"add_on_deal"`
	MainItem             bool    `json:"main_item"`
	AddOnDealID          int64   `json:"add_on_deal_id"`
	PromotionType        string  `json:"promotion_type"`
	PromotionID          int64   `json:"promotion_id"`
	OrderItemID          int64   `json:"order_item_id"`
	PromotionGroupID     int     `json:"promotion_group_id"`
	ImageInfo            IResOrderDetailImage   `json:"image_info"`
	ProductLocationID    []string  `json:"product_location_id"`
	IsPrescriptionItem   bool    `json:"is_prescription_item"`
	IsB2COwnedItem       bool    `json:"is_b2c_owned_item"`
}

type IResOrderDetailPackItem struct {
	ItemID            int64  `json:"item_id"`
	ModelID           int64  `json:"model_id"`
	ModelQuantity     int    `json:"model_quantity"`
	OrderItemID       int64  `json:"order_item_id"`
	PromotionGroupID  int    `json:"promotion_group_id"`
	ProductLocationID string `json:"product_location_id"`
}

type IResOrderDetailPackage struct {
	PackageNumber         string      `json:"package_number"`
	LogisticsStatus       string      `json:"logistics_status"`
	LogisticsChannelID    int64       `json:"logistics_channel_id"`
	ShippingCarrier       string      `json:"shipping_carrier"`
	AllowSelfDesignAWB    bool        `json:"allow_self_design_awb"`
	ItemList              []IResOrderDetailPackItem  `json:"item_list"`
	GroupShipmentID       int64       `json:"group_shipment_id"`
	ParcelChargeableWeight int        `json:"parcel_chargeable_weight"`
	PackageQueryNumber    string      `json:"package_query_number"`
	SortingGroup          string      `json:"sorting_group"`
}

type IResOrderDetailInvoice struct {
	Number             string  `json:"number"`
	SeriesNumber       string  `json:"series_number"`
	AccessKey          string  `json:"access_key"`
	IssueDate          int64   `json:"issue_date"`
	TotalValue         float64 `json:"total_value"`
	ProductsTotalValue float64 `json:"products_total_value"`
	TaxCode            string  `json:"tax_code"`
}

type IResOrderListWithDetails struct {
	OrderSN                 string           `json:"order_sn"`
	Region                  string           `json:"region"`
	Currency                string           `json:"currency"`
	COD                     bool             `json:"cod"`
	TotalAmount             float64          `json:"total_amount"`
	PendingTerms            []string         `json:"pending_terms"`
	OrderStatus             string           `json:"order_status"`
	ShippingCarrier         string           `json:"shipping_carrier"`
	PaymentMethod           string           `json:"payment_method"`
	EstimatedShippingFee    float64          `json:"estimated_shipping_fee"`
	MessageToSeller         string           `json:"message_to_seller"`
	CreateTime              int64            `json:"create_time"`
	UpdateTime              int64            `json:"update_time"`
	DaysToShip              int              `json:"days_to_ship"`
	ShipByDate              int64            `json:"ship_by_date"`
	BuyerUserID             int              `json:"buyer_user_id"`
	BuyerUsername           string           `json:"buyer_username"`
	RecipientAddress        IResOrderDetailReceiptAddress `json:"recipient_address"`
	ActualShippingFee       float64          `json:"actual_shipping_fee"`
	GoodsToDeclare          bool             `json:"goods_to_declare"`
	Note                    string           `json:"note"`
	NoteUpdateTime          int64            `json:"note_update_time"`
	ItemList                []IResOrderDetailItem `json:"item_list"`
	PayTime                 int64            `json:"pay_time"`
	Dropshipper             string           `json:"dropshipper"`
	DropshipperPhone        string           `json:"dropshipper_phone"`
	SplitUp                 bool             `json:"split_up"`
	BuyerCancelReason       string           `json:"buyer_cancel_reason"`
	CancelBy                string           `json:"cancel_by"`
	CancelReason            string           `json:"cancel_reason"`
	ActualShippingFeeConfirmed bool             `json:"actual_shipping_fee_confirmed"`
	BuyerCPFID              string           `json:"buyer_cpf_id"`
	FulfillmentFlag         string           `json:"fulfillment_flag"`
	PickupDoneTime          int64            `json:"pickup_done_time"`
	PackageList             []IResOrderDetailPackage      `json:"package_list"`
	InvoiceData             IResOrderDetailInvoice        `json:"invoice_data"`
	CheckoutShippingCarrier string           `json:"checkout_shipping_carrier"`
	ReverseShippingFee      float64          `json:"reverse_shipping_fee"`
	OrderChargeableWeight   int              `json:"order_chargeable_weight_gram"`
	PrescriptionImages      []string         `json:"prescription_images"`
	PrescriptionStatus      int              `json:"prescription_check_status"`
	EDTFrom                 int64            `json:"edt_from"`
	EDTTo                   int64            `json:"edt_to"`
	BookingSN               string           `json:"booking_sn"`
	AdvancePackage          bool             `json:"advance_package"`
	ReturnRequestDueDate    int64            `json:"return_request_due_date"`
}

type IResOrderDetailByOrderSNShopWrapper struct {
	OrderList  []IResOrderListWithDetails  `json:"order_list"` // parse from response
}

type IResOrderDetailByOrderSN struct {
  IResShopeeResponse 
  Warning  []string `json:"warning,omitempty"`
  Response IResOrderDetailByOrderSNShopWrapper `json:"response"`

}

// Reminder Note :
// * Embeded all struct

// path : /api/v2/shop/get_shop_info
type IResSipAffiShopsDTO struct {
  AffiShopID int    `json:"affi_shop_id"`
  Region     string `json:"region"`
}

type IResLinkedDirectShopListDTO struct {
  DirectShopID int `json:"direct_shop_id"`
  DirectShopRegion string `json:"direct_shop_region"`
}

type IResOutletShopInfoList struct {
  OutletShopID int `json:"outlet_shop_id"`
}

type IResShopGetShopInfoDTO struct {
  IResShopeeResponse
  ShopName      string   `json:"shop_name"`
  Region        string   `json:"region"`
  Status        string   `json:"status"`
  SipAffiShops  []IResSipAffiShopsDTO `json:"sip_affi_shops"`
  IsCB          bool   `json:"is_cb"`
  IsSip          bool  `json:"is_sip"`
  ISUpgradedCBSC bool  `json:"is_upgraded_cbsc"`
  MerchantID     int   `json:"merchant_id"`
  ShopFullFilmentFlag   string `json:"shop_fulfillment_flag"`
  IsMainShop     bool  `json:"is_main_shop"`
  IsDirectShop   bool  `json:"is_direct_shop"`
  LinkedMainShopID int `json:"linked_main_shop_id"`
  LinkedDirectShopList []IResLinkedDirectShopListDTO `json:"linked_direct_shop_list"`
  IsOneAwb     bool    `json:"is_one_awb"`
  IsMartShop   bool    `json:"is_mart_shop"`
  IsOutletShop bool    `json:"is_outlet_shop"`
  MartShopID   int     `json:"mart_shop_id"`
  OutletShopInfoList   []IResOutletShopInfoList `bson:"outlet_shop_info_list"`
} 

// path : /api/v2/shop/get_profile
type IResShopGetProfile_ResponseDTO struct {
  ShopLogo      string `json:"shop_logo"`
  Description   string `json:"description"`
  ShopName      string `json:"shop_name"`
  InvoiceIssuer string `json:"invoice_issuer"`

}

type IResShopGetProfileDTO struct {
  IResShopeeResponse 
  Resoponse IResShopGetProfile_ResponseDTO `json:"response"`
}

// [Rules] 
// path : ***/***
// type [IRes|IReq]:Table:Method:DTO
// 
