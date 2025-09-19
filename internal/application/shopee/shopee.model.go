package shopee

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ShopeeAuthModel struct {
  PartnerID    string    `bson:"partner_id"`
	ShopID       string    `bson:"shop_id"`
  Code         string    `bson:"code"`
	AccessToken  string    `bson:"access_token"`
	RefreshToken string    `bson:"refresh_token"`
	ExpiredAt    time.Time `bson:"expired_at"`

	CreatedAt   time.Time `bson:"created_at"`
	CreatedBy   string    `bson:"created_by"`

	MoidifiedAt time.Time `bson:"modified_at"`
	ModifiedBy  string    `bson:"modified_by"`
}

type ShopeeAuthRequestModel struct {
	PartnerID   string    `bson:"partner_id"`
	PartnerKey  string    `bson:"partner_key"`
	PartnerName string    `bson:"partner_name"`
  GeneratedUrl string   `bson:"generated_url"`

	CreatedAt time.Time   `bson:"created_at"`
	CreatedBy string      `bson:"created_by"`

	MoidifiedAt time.Time `bson:"modified_at"`
	ModifiedBy  string    `bson:"modified_by"`
}

type ShopeePartnerModel struct {  
  PartnerID   string    `bson:"partner_id"`
	PartnerKey  string    `bson:"partner_key"`
	PartnerName string    `bson:"partner_name"`

  CreatedAt time.Time   `bson:"created_at"`
	CreatedBy string      `bson:"created_by"`

	MoidifiedAt time.Time `bson:"modified_at"`
	ModifiedBy  string    `bson:"modified_by"`
}

// in DB Collection
// -- ShopeeShopProfle
// -- ShopeeShopDetails Collections

type SipAffiShops_Struct struct {
  AffiShopID  string   `bson:"affi_shop_id" json:"affi_shop_id"`
  Region      string   `bson:"region" json:"region"`
}

type ShopFullFilmentFlagEnum string 

const (
  PURE_FBS_SHOP ShopFullFilmentFlagEnum = "Pure-FBS Shop"
  PURE_3PF_SHOP ShopFullFilmentFlagEnum = "Pure-3PF Shop"
  PFF_FBS_SHOP  ShopFullFilmentFlagEnum = "PFF-FBS Shop"
  PFF_3PF_SHOP  ShopFullFilmentFlagEnum = "PFF-3PF Shop"
  LFF_HYBRID_SHOP ShopFullFilmentFlagEnum = "LFF Hybrid Shop"
  OTHERS   ShopFullFilmentFlagEnum = "Others"
  UNKNOWN  ShopFullFilmentFlagEnum = "Unknown"
)

type LinkedDirectShopList_Struct struct {
  DirectShopID string `bson:"direct_shop_id" json:"direct_shop_id"`
  DirectShopRegion string `bson:"direct_shop_region" json:"direct_shop_region"`
} 

type OutletShopInfoList_Struct struct {
  OutletShopID string `bson:"outlet_shop_id" json:"outlet_shop_id"`
}

type ShopeeShopStatusEnum string

const (
  NORMAL ShopeeShopStatusEnum = "NORMAL"
  BANNED ShopeeShopStatusEnum = "BANNED"
  FROZEN ShopeeShopStatusEnum = "FROZEN"
)
// Core profile 
type ShopeeShopDetailsModel struct {
  ID            bson.ObjectID `bson:"_id"`
  ShopID        string `bson:"shop_id"`

  ShopLogo      string `bson:"shop_logo"`
  Description   string `bson:"description"`
  InvoiceIssuer string `bson:"invoice_issuer "`
  ShopName      string `bson:"shop_name"`

  Region      string   `bson:"region"`
  Status      ShopeeShopStatusEnum  `bson:"status"`
  SipAffiShops []SipAffiShops_Struct `bson:"sip_affi_shops"`
  IsCB  bool `bson:"is_cb"`
  IsSip bool `bson:"is_sip"`
  IsUpgradedCBSC bool `bson:"is_upgraded_cbsc"`
  MerchantID string `bson:"merchant_id"`
  ShopFullFilmentFlag ShopFullFilmentFlagEnum `bson:"shop_fulfillment_flag"`
  IsMainShop bool `bson:"is_main_shop"`
  IsDirectShop bool `bson:"is_direct_shop"`
  LinkedMainShopID string `bson:"linked_main_shop_id"`
  LinkedDirectShopList []LinkedDirectShopList_Struct `bson:"linked_direct_shop_list"`
  IsOneAwb bool `bson:"is_one_awb"`
  IsMartShop bool `bson:"is_mart_shop"`
  IsOutletShop bool `bson:"is_outlet_shop"`
  MartShopID string `bson:"mart_shop_id"`
  OutletShopInfoList []OutletShopInfoList_Struct `bson:"outlet_shop_info_list"`

  CreatedAt time.Time `bson:"created_at"`
  CreatedBy string `bson:"created_by"`
  UpdatedAt time.Time `bson:"updated_at"`
  UpdatedBy string `bson:"updated_by"`
}

// ----------------- [Model] - Start.Collection("shop_order") ----------------
// [Concept] : 
// [Compnent Struct.Start]
// [Component Struct.End]
// [Core Struct.Start]
type ShopeeRecipientAddressModel struct {
  Name     string `bson:"name"`
  Phone    string `bson:"phone"`
  Town     string `bson:"town"`
  District string `bson:"district"`
  City     string `bson:"city"`
  State    string `bson:"state"`
  Region   string `bson:"region"`
  ZipCode  string `bson:"zip_code"`
  FullAddress string `bson:"full_address"`
}

type ShopeeItemListModel struct {
  ItemID    string `bson:"item_id"`
  ItemName  string `bson:"item_name"`
  ItemSKU   string `bson:"item_sku"`
  ModelID   string `bson:"model_id"`
  ModelName string `bson:"model_name"`
  ModelSKU  string `bson:"model_sku"`
  ModelQualityPurchased int `bson:"model_quality_purchased"`
  ModelOriginPrice      float64 `bson:"model_origin_price"`
  ModelDiscountedPrice  float64 `bson:"model_discount_price"`
  WholeSale bool `bson:"whole_sale"`
  Weight    float64 `bson:"weight"`
  AddOnDeal bool `bson:"add_on_deal"`
  MainItem  bool `bson:"main_item"`
  AddOnDealID   string      `bson:"add_on_deal_id"`
  PromotionType ShopeePromotionTypeEnum `bson:"promotion_type"`
  PromotionID string `bson:"promotion_id"`
  OrderItemID string `bson:"order_item_id"`
  PromotionGroupID string   `bson:"promotion_group_id"`
  ImageInfo ShopeeImageInfoEntity `bson:"image_info"`
  ProductLocationID  []string  `bson:"production_location_id"`
  IsPrescriptionItem bool   `bson:"is_prescription_item"`
  IsB2COwnedItem bool       `bson:"is_b2c_owned_item"`
}

type ShopeeItemListInPackageListModel struct {
  ItemID  string `bson:"item_id"`
  ModelID string `bson:"model_id"`
  ModelQuantity       int     `bson:"model_quantity"`
  OrderItemID         string  `bson:"order_item_id"`
  PromotionGroupID    string  `bson:"promotion_group_id"`
  ProductLocationID   string  `bson:"product_location_id"`
}

type ShopeePackageListModel struct {
    PackageNumber       string
    LogisticsStatus     string
    LogisticsChannelID  string
    ShippingCarrier     string
    AllowSelfDesignAWB  bool 
    ItemList []ShopeeItemListInPackageListModel
    ParcelChargeableWeight int
    GroupShipmentID string
}

type ShopeeInvoiceDataModel struct {
	Number            string
	SeriesNumber      string
	AccessKey         string
	IssueDate         time.Time
	TotalValue        float64
	ProductTotalValue float64
	TaxCode           string
}

type ShopeeOrderModel struct {
  ID bson.ObjectID    `bson:"_id"`
  ShopID    string    `bson:"shop_id"` 
  
  OrderSN   string    `bson:"order_sn"`
  BookingSN string    `bson:"booking_sn"`
  OrderStatus ShopeeOrderStatusEnum `bson:"order_status"`
  Region       string `bson:"region"` 
  Currency     string `bson:"currency"`
  Cod          bool   `bson:"cod"` 
  TotalAmount  float64   `bson:"total_amount"`
  PendingTerms  []string `bson:"pending_terms"` 
  ShippingCarrier string `bson:"shipping_carrier"`
  PaymentMethod   string `bson:"payment_method"`
  EstimatedShippingFee float64 `bson:"estimated_shipping_fee"`
  MessageToSeller string  `bson:"message_to_seller"`
  CreateTime    time.Time `bson:"create_time"`
  UpdateTime    time.Time `bson:"update_time"`
  DaysToShip    int               `bson:"days_to_ship"`
  ShipByDate    int               `bson:"ship_by_date"` // ??
  BuyerUserId   string            `bson:"buyer_user_id"`
  BuyerUsername string            `bson:"buyer_username"`
  RecipientAddress ShopeeRecipientAddressModel `bson:"recipient_address"`
  ActualShippingFee float64       `bson:"actual_shipping_fee"`
  GoodsToDeclare bool             `bson:"goods_to_declare"`
  Note          string            `bson:"note"`
  NoteUpdateTime time.Time        `bson:"note_update_time"`
  ItemList       []ShopeeItemListModel `bson:"item_list"`

  PayTime time.Time               `bson:"pay_time"`
  DropShipper string              `bson:"dropshipper"`
  DropShipperPhone string         `bson:"dropshipper_phone"`
  SplitUp   bool                  `bson:"split_up"`
  BuyerCancelReason string        `bson:"buyer_cancel_reason"`
  CancelBy  string                `bson:"cancel_by"`
  CancelReason string             `bson:"cancel_reason"`
  ActualShippingFeeConfirmed bool `bson:"actual_shipping_fee_confirmed"`
  BuyerCPFID string               `bson:"buyer_cpf_id"`
  FulFillmentFlag ShopeeFulfillmentFlagEnum `bson:"fulfillment_flag"`
  PickupDoneTime time.Time        `bson:"pickup_done_time"`
  PackageList []ShopeePackageListModel  `bson:"package_list"` 
  InvoiceData ShopeeInvoiceDataModel    `bson:"invoice_data"`
  CheckoutShippingCarrier string  `bson:"checkout_shipping_carrier"`
  ReverseShippingFee float64      `bson:"reverse_shipping_fee"`
  OrderChargeableWeightGram int   `bson:"order_chargeable_weight_gram"`
  PrescriptionImages []string     `bson:"prescription_images"`
  PrescriptionCheckStatus ShopeePrescriptionCheckStatusEnum `bson:"prescription_check_status"` 
  AdvancePackage bool             `bson:"advance_package"`
  ReturnRequestDueDate     time.Time  `bson:"return_request_due_date"`

  CreatedAt   time.Time   `bson:"created_at"`
  CreatedBy   string      `bson:"created_by"`
  UpdatedAt   time.Time   `bson:"updated_at"`
  UpdatedBy   string      `bson:"updated_by"`
}
// [Core Struct.Emd]
// [Method Start]
func ShopeeOrderModelToEntity(model *ShopeeOrderModel) *ShopeeOrderEntity {

  recipient := ShopeeRecipientAddressEntity{
    Name:   model.RecipientAddress.Name,
    Phone:  model.RecipientAddress.Phone,
    Town:   model.RecipientAddress.Town,
    District: model.RecipientAddress.District,
    City:   model.RecipientAddress.City,
    State:  model.RecipientAddress.State,
    Region: model.RecipientAddress.Region,
    ZipCode: model.RecipientAddress.ZipCode,
    FullAddress: model.RecipientAddress.FullAddress,
  }

  items := []ShopeeItemListEntity{}
  if len(model.ItemList) > 0 {
    for _,i := range model.ItemList {
      items = append(items, ShopeeItemListEntity{
        ItemID: i.ItemID,
        ItemName: i.ItemName,
        ItemSKU: i.ItemSKU,
        ModelID: i.ModelID,
        ModelName: i.ModelName,
        ModelSKU: i.ModelSKU,
        ModelQualityPurchased: i.ModelQualityPurchased,
        ModelOriginPrice: i.ModelOriginPrice,
        ModelDiscountedPrice: i.ModelDiscountedPrice,
        WholeSale: i.WholeSale,
        Weight: i.Weight,
        AddOnDeal: i.AddOnDeal,
        MainItem: i.MainItem,
        AddOnDealID: i.AddOnDealID,
        PromotionType: i.PromotionType,
        PromotionID: i.PromotionID,
        OrderItemID: i.OrderItemID,
        PromotionGroupID: i.PromotionGroupID,
        ImageInfo: ShopeeImageInfoEntity{ ImageURL: i.ImageInfo.ImageURL, },
        ProductLocationID: i.ProductLocationID,
        IsPrescriptionItem: i.IsPrescriptionItem,
        IsB2COwnedItem: i.IsB2COwnedItem, })
    }
  }

  packagesPack := []ShopeePackageListEntity{} 

  if len(model.PackageList) > 0 {
    for _,pk := range model.PackageList {

      pkItems := []ShopeeItemListInPackageListEntity{} 
      if len(pk.ItemList) > 0 {
        for _,i := range pk.ItemList {
          pkItems = append(pkItems, ShopeeItemListInPackageListEntity{
            ItemID: i.ItemID,
            ModelID: i.ModelID,
            ModelQuantity: i.ModelQuantity,
            OrderItemID: i.OrderItemID,
            PromotionGroupID: i.PromotionGroupID,
            ProductLocationID: i.ProductLocationID,
          })
        }
      }


      packagesPack = append(packagesPack, ShopeePackageListEntity{
        PackageNumber: pk.PackageNumber,
        LogisticsStatus: pk.LogisticsStatus,
        LogisticsChannelID: pk.LogisticsChannelID,
        ShippingCarrier: pk.ShippingCarrier,
        AllowSelfDesignAWB: pk.AllowSelfDesignAWB,
        ItemList: pkItems,
        ParcelChargeableWeight: pk.ParcelChargeableWeight,
        GroupShipmentID: pk.GroupShipmentID,
      })
    }
  } 
  

  
  details := ShopeeOrderDetailsEntity{
    Region: model.Region,
    Currency: model.Currency,
    Cod: model.Cod,
    TotalAmount: model.TotalAmount,
    PendingTerms: model.PendingTerms,

    ShippingCarrier: model.ShippingCarrier,
    PaymentMethod: model.PaymentMethod,
    EstimatedShippingFee:  model.EstimatedShippingFee,
    MessageToSeller: model.MessageToSeller,
    CreateTime: model.CreateTime,
    UpdateTime: model.UpdateTime,

    DaysToShip: model.DaysToShip,
    ShipByDate: model.ShipByDate,
    BuyerUserId: model.BuyerUserId,
    BuyerUsername: model.BuyerUsername,
    RecipientAddress: recipient,
    ActualShippingFee: model.ActualShippingFee,
    GoodsToDeclare: model.GoodsToDeclare,
    Note: model.Note,
    NoteUpdateTime: model.NoteUpdateTime,
    ItemList: items,
    PayTime: model.PayTime,
    DropShipper: model.DropShipper,
    DropShipperPhone: model.DropShipperPhone,
    SplitUp: model.SplitUp,
    BuyerCancelReason: model.BuyerCancelReason,
    CancelBy: model.CancelBy,
    CancelReason: model.CancelReason,
    ActualShippingFeeConfirmed: model.ActualShippingFeeConfirmed,
    BuyerCPFID: model.BuyerCPFID,
    FulFillmentFlag: model.FulFillmentFlag,
    PickupDoneTime: model.PickupDoneTime,
    PackageList: packagesPack,
      
    InvoiceData: ShopeeInvoiceDataEntity{
      Number: model.InvoiceData.Number,
      SeriesNumber: model.InvoiceData.SeriesNumber,
      AccessKey: model.InvoiceData.AccessKey,
      IssueDate: model.InvoiceData.IssueDate,
      TotalValue: model.InvoiceData.TotalValue,
      ProductTotalValue: model.InvoiceData.ProductTotalValue,
      TaxCode: model.InvoiceData.TaxCode,
    },
    CheckoutShippingCarrier: model.CheckoutShippingCarrier,
    ReverseShippingFee: model.ReverseShippingFee,
    OrderChargeableWeightGram: model.OrderChargeableWeightGram,
    PrescriptionImages: model.PrescriptionImages,
    PrescriptionCheckStatus: model.PrescriptionCheckStatus,
    BookingSN: model.BookingSN,
    AdvancePackage: model.AdvancePackage,
    ReturnRequestDueDate: model.ReturnRequestDueDate,
  }

  return &ShopeeOrderEntity{
    ID: model.ID.Hex(),
    ShopID: model.ShopID,
    OrderSN: model.OrderSN,
    OrderStatus: model.OrderStatus,
    BookingSN: model.BookingSN,
    ShopeeOrderDetailsEntity: details,
  }
}

// [Method End]
// ----------------- [Model] - End.Collection("shop_order")   ----------------
// ---------------------------------------------- Demo Struct Template ------------------------------------
// tyoe xxxx_StructEntity Struct {} : special case  for Universal qStruct (DTO , Model , Entity)
// type DemoModel Struct {
//   ... 
//  CreatedAt time.Time   `bson:"created_at"`
// 	CreatedBy string      `bson:"created_by"`
// 	MoidifiedAt time.Time `bson:"modified_at"`
// 	ModifiedBy  string    `bson:"modified_by"`
// }
// ---------------------------------------------- End Demo Struct Template --------------------------------
// Template
// ----------------- [DTO/Entity/Model] - Start.Collection("Shop_rder") ----------------
// [Concept] : xxxx
// [Compnent Struct.Start]
// [Component Struct.End]
// [Core Struct.Start]
// [Core Struct.Emd]
// [Method Start]
// [Method End]
// ----------------- [DTO/Entity/Model] - End.Collection("Shop_rder")   ----------------
