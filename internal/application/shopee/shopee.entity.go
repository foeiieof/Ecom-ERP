package shopee

import (
	"ecommerce/internal/adapter/dto"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ShopeeAuthEntity struct {
	PartnerID string
	ShopID    string // Code          string
	AccessToken  string
	RefreshToken string
	ExpiredAt    time.Time

	// CreatedAt     time.Time
	// CreatedBy     string

	// MoidifiedAt   time.Time
	// ModifiedBy    string
}

type ShopeeShopListEntity struct {
	ShopList []dto.IResAuthedShopList
}

type ShopeeOrderListEntity struct {
  OrderList []ShopeeOrderEntity `json:"order_list"`
}

type ShopeeOrderListWithDetailEntity struct {
  OrderList []dto.IResOrderListWithDetails
}


// type SipAffiShopsEntity struct {
//   AffiShopID string  
//   Region string
// }
// type LinkedDirectShopListEntity struct {
//   DirectShopID      string
//   DirectShopRegion  string 
// }
// type OutletShopInfoListEntity struct {
//   OuletShopID string 
// }
type ShopeeShopDetailsEntityDTO struct {
  ID            string 
  ShopID        string 
  ShopLogo      string 
  Description   string
  ShopName      string 
  InvoiceIssuer string


  Region string  
  Status ShopeeShopStatusEnum 
  SipAffiShops []SipAffiShops_Struct
  IsCB  bool
  IsSip bool 
  IsUpgradedCBSC bool 
  MerchantID string 
  ShopFullFilmentFlag ShopFullFilmentFlagEnum
  IsMainShop bool 
  IsDirectShop bool
  LinkedMainShopID string 
  LinkedDirectShopList []LinkedDirectShopList_Struct
  IsOneAwb bool
  IsMartShop bool 
  IsOutletShop bool 
  MartShopID string 
  OutletShopInfoList []OutletShopInfoList_Struct
  CreatedAt time.Time
  CreatedBy string 
  UpdatedAt time.Time 
  UpdatedBy string 
}


func ShopeeShopEntityToModel(enti ShopeeShopDetailsEntityDTO) *ShopeeShopDetailsModel{
  Oid := bson.NewObjectID()
  if enti.ID != "" { if ok,err := bson.ObjectIDFromHex(enti.ID); err == nil { Oid = ok } }

  return &ShopeeShopDetailsModel{
    ID: Oid,
    ShopID: enti.ShopID,
    ShopLogo: enti.ShopLogo,
    Description: enti.Description,
    ShopName: enti.ShopName,

    Region: enti.Region,
    Status: enti.Status,
    SipAffiShops: enti.SipAffiShops ,
    IsCB: enti.IsCB,
    IsSip: enti.IsSip,
    IsUpgradedCBSC: enti.IsUpgradedCBSC,
    MerchantID: enti.MerchantID,
    ShopFullFilmentFlag: enti.ShopFullFilmentFlag,
    IsMainShop: enti.IsMainShop,
    IsDirectShop: enti.IsDirectShop,
    LinkedMainShopID: enti.LinkedMainShopID,
    LinkedDirectShopList: enti.LinkedDirectShopList,
    IsOneAwb: enti.IsOneAwb,
    IsMartShop: enti.IsMartShop,
    IsOutletShop: enti.IsOutletShop,
    MartShopID: enti.MartShopID,
    OutletShopInfoList: enti.OutletShopInfoList,

    CreatedAt: enti.CreatedAt,
    CreatedBy: enti.CreatedBy,
    UpdatedAt: enti.CreatedAt,
    UpdatedBy: enti.UpdatedBy,
  }
}


func ShopeeShopModelToEntity(model *ShopeeShopDetailsModel) *ShopeeShopDetailsEntityDTO {
  return &ShopeeShopDetailsEntityDTO{
    ID: model.ID.Hex(),
    ShopID: model.ShopID,
    ShopLogo: model.ShopLogo,
    Description: model.Description,
    ShopName: model.ShopName,
    InvoiceIssuer: model.InvoiceIssuer,

    Region: model.Region,
    Status: model.Status,
    SipAffiShops: model.SipAffiShops,
    IsCB: model.IsCB,
    IsSip: model.IsSip,
    IsUpgradedCBSC: model.IsUpgradedCBSC,
    MerchantID: model.MerchantID,
    ShopFullFilmentFlag: model.ShopFullFilmentFlag,
    IsMainShop: model.IsMainShop,
    IsDirectShop: model.IsDirectShop,
    LinkedMainShopID: model.LinkedMainShopID,
    LinkedDirectShopList: model.LinkedDirectShopList,
    IsOneAwb: model.IsOneAwb,
    IsMartShop: model.IsMartShop,
    IsOutletShop: model.IsOutletShop,
    MartShopID: model.MartShopID,
    OutletShopInfoList: model.OutletShopInfoList,

    CreatedAt: model.CreatedAt,
    CreatedBy: model.CreatedBy,
    UpdatedAt: model.UpdatedAt,
    UpdatedBy: model.UpdatedBy,
  }
}

// ----------------- [Entity] - Start.Collection("shopee_oder") ----------------
// [Concept.Start]  
// shopee_order_list    -----\
//                           |--> ShopeeOrderDetails
// shopee_order_details -----\
// [Concept.End]

// [Compnent Struct.Start]
type ShopeeOrderStatusEnum string
const (
  UNPAID      ShopeeOrderStatusEnum = "UNPAID"
  READYTOSHIP ShopeeOrderStatusEnum = "READY_TO_SHIP"
  PROCESSED   ShopeeOrderStatusEnum = "PROCESSED"
  RETRYSHIP   ShopeeOrderStatusEnum = "RETRY_SHIP"
  SHIPPED     ShopeeOrderStatusEnum = "SHIPPED"
  TOCONFIRMRECEIVE ShopeeOrderStatusEnum = "TO_CONFIRM_RECEIVE"
  INCANCEL    ShopeeOrderStatusEnum = "IN_CANCEL"
  CANCELLED   ShopeeOrderStatusEnum = "CANCELLED"
  TORETURN    ShopeeOrderStatusEnum = "TO_TO_RETURN"
  COMPLETED   ShopeeOrderStatusEnum = "COMPLETED"
)
// type ShopeeOrderTermEnum string
// const (
//   SYSTEMPENDING ShopeeOrderTermEnum = "SYSTEM_PENDING"
//   KYCPENDING    ShopeeOrderTermEnum = "KYC_PENDING"
// )
type ShopeeRecipientAddressEntity struct {
  Name     string
  Phone    string
  Town     string
  District string
  City     string
  State    string
  Region   string
  ZipCode  string
  FullAddress string
}
type ShopeePromotionTypeEnum string
const (
  PRODUCTPROMOTION ShopeePromotionTypeEnum = "product_promotion"
  FLASHSALE     ShopeePromotionTypeEnum = "flash_sale"
  BUNDLEDEAL    ShopeePromotionTypeEnum = "bundle_deal"
  ADDONDEALMAIN ShopeePromotionTypeEnum = "add_on_deal_main"
  ADDONDEALSUB  ShopeePromotionTypeEnum = "add_on_deal_sub"
)
type ShopeeImageInfoEntity struct {
  ImageURL string
}

// []Object
type ShopeeItemListInPackageListEntity struct {
  ItemID  string
  ModelID string
  ModelQuantity     int
  OrderItemID       string
  PromotionGroupID  string
  ProductLocationID string
}
// []Object
type ShopeePackageListEntity struct {
  PackageNumber       string
  LogisticsStatus     string
  LogisticsChannelID  string
  ShippingCarrier     string
  AllowSelfDesignAWB  bool 
  ItemList []ShopeeItemListInPackageListEntity 
  ParcelChargeableWeight int
  GroupShipmentID string
}
// Object 
type ShopeeInvoiceDataEntity struct {
  Number       string
  SeriesNumber string
  AccessKey    string
  IssueDate    time.Time
  TotalValue   float64
  ProductTotalValue float64
  TaxCode      string
}
type ShopeePrescriptionCheckStatusEnum int 
const(
  NONE ShopeePrescriptionCheckStatusEnum = iota
  PASSED
  FAILED
)
type ShopeeFulfillmentFlagEnum string
const(
  FULFILBYSHOPEE = "fulfilled_by_shopee"
  FULFILBYCB = "fulfilled_by_cb_seller"
  FULFILBYLOCAL = "fulfilled_by_local_seller"
)
type ShopeeItemListEntity struct {
  ItemID    string
  ItemName  string
  ItemSKU   string
  ModelID   string
  ModelName string
  ModelSKU  string
  ModelQualityPurchased int
  ModelOriginPrice      float64
  ModelDiscountedPrice  float64
  WholeSale bool
  Weight    float64 
  AddOnDeal bool
  MainItem  bool
  AddOnDealID   string
  PromotionType ShopeePromotionTypeEnum
  PromotionID string
  OrderItemID string
  PromotionGroupID string
  ImageInfo ShopeeImageInfoEntity
  ProductLocationID []string
  IsPrescriptionItem bool
  IsB2COwnedItem bool
}
// Component Struct
type ShopeeOrderDetailsEntity struct {
  // ordersn 
  Region       string `json:"region"`
  Currency     string
  Cod          bool
  TotalAmount  float64
  PendingTerms []string
  // order_status 
  ShippingCarrier string
  PaymentMethod   string
  EstimatedShippingFee float64
  MessageToSeller string
  CreateTime    time.Time
  UpdateTime    time.Time 
  DaysToShip    int
  ShipByDate    int // ??
  BuyerUserId   string
  BuyerUsername string
  RecipientAddress ShopeeRecipientAddressEntity
  ActualShippingFee float64
  GoodsToDeclare bool 
  Note          string
  NoteUpdateTime time.Time
  ItemList       []ShopeeItemListEntity

  PayTime time.Time
  DropShipper string
  DropShipperPhone string
  SplitUp   bool
  BuyerCancelReason string
  CancelBy  string
  CancelReason string
  ActualShippingFeeConfirmed bool
  BuyerCPFID string
  FulFillmentFlag ShopeeFulfillmentFlagEnum
  PickupDoneTime time.Time

  PackageList []ShopeePackageListEntity
  InvoiceData ShopeeInvoiceDataEntity

  CheckoutShippingCarrier string
  ReverseShippingFee float64
  OrderChargeableWeightGram int 
  PrescriptionImages []string
  PrescriptionCheckStatus ShopeePrescriptionCheckStatusEnum  
  BookingSN string
  AdvancePackage bool
  ReturnRequestDueDate time.Time
  // payment_info : []object [only for BR]
  CreatedAt   time.Time
  CreatedBy   string
  UpdatedAt time.Time
  UpdatedBy  string

}
// [Component Struct.End]

// [Core Struct.Start]
type ShopeeOrderEntity struct {
  ID          string `json:"id"`
  ShopID      string `json:"shop_id"`
  OrderSN     string `json:"order_sn"`
  OrderStatus ShopeeOrderStatusEnum `json:"order_status"`
  BookingSN   string `json:"booking_sn"`
  ShopeeOrderDetailsEntity 
}
// [Core Struct.Emd]
// [Method Start]
func ShopeeOrderEntityToModel(enti *ShopeeOrderEntity) *ShopeeOrderModel{
  oID := bson.NewObjectID()

  if enti.ID != "" {
    if ok,err := bson.ObjectIDFromHex(enti.ID); err == nil {
      oID = ok
    } 
  }


  items := []ShopeeItemListModel{}
  if len (enti.ItemList) > 0 {
    for _,i := range enti.ItemList {
      items = append(items, ShopeeItemListModel{
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
        ImageInfo: i.ImageInfo,
        ProductLocationID: i.ProductLocationID,
        IsPrescriptionItem: i.IsPrescriptionItem,
        IsB2COwnedItem: i.IsB2COwnedItem,
      })
    } 
  } 

  packagesL := []ShopeePackageListModel{}
  if len(enti.PackageList) > 0 {

    for _,p := range enti.PackageList {
      items := []ShopeeItemListInPackageListModel{}
      if len(p.ItemList) > 0 {
        for _, i := range p.ItemList  {
          items = append(items, ShopeeItemListInPackageListModel{
            ItemID: i.ItemID,
            ModelID: i.ModelID,
            ModelQuantity: i.ModelQuantity,
            OrderItemID: i.OrderItemID,
            PromotionGroupID: i.PromotionGroupID,
            ProductLocationID: i.ProductLocationID,
          })
        }
      }

      packagesL = append(packagesL, ShopeePackageListModel{
        PackageNumber: p.PackageNumber,
        LogisticsStatus: p.LogisticsStatus,
        LogisticsChannelID: p.LogisticsChannelID,
        ShippingCarrier: p.ShippingCarrier,
        AllowSelfDesignAWB: p.AllowSelfDesignAWB,
        ItemList: items ,
        ParcelChargeableWeight: p.ParcelChargeableWeight,
        GroupShipmentID: p.GroupShipmentID,
      })
    }
  }

  return &ShopeeOrderModel{
    ID: oID,
    ShopID: enti.ShopID,
    
    OrderSN: enti.OrderSN,
    BookingSN: enti.BookingSN,
    OrderStatus: enti.OrderStatus,
    Region: enti.Region,
    Currency: enti.Currency,
    Cod: enti.Cod,
    TotalAmount: enti.TotalAmount,
    PendingTerms: enti.PendingTerms,
    ShippingCarrier: enti.ShippingCarrier,
    PaymentMethod: enti.PaymentMethod,
    EstimatedShippingFee: enti.EstimatedShippingFee,
    MessageToSeller: enti.MessageToSeller,
    CreateTime: enti.CreateTime,
    UpdateTime: enti.UpdateTime,
    DaysToShip: enti.DaysToShip,
    ShipByDate: enti.ShipByDate,
    BuyerUserId: enti.BuyerUserId,
    BuyerUsername: enti.BuyerUsername,
    
    RecipientAddress: ShopeeRecipientAddressModel(enti.RecipientAddress),

    ActualShippingFee: enti.ActualShippingFee,
    GoodsToDeclare: enti.GoodsToDeclare,
    Note: enti.Note,
    NoteUpdateTime: enti.NoteUpdateTime,
    
    ItemList: items ,
    PayTime: enti.PayTime,
    DropShipper: enti.DropShipper,
    DropShipperPhone: enti.DropShipperPhone,
    SplitUp: enti.SplitUp,
    BuyerCancelReason: enti.BuyerCancelReason,
    CancelBy: enti.CancelBy,
    CancelReason: enti.CancelReason,
    ActualShippingFeeConfirmed: enti.ActualShippingFeeConfirmed,
    BuyerCPFID: enti.BuyerCPFID,
    FulFillmentFlag: enti.FulFillmentFlag,
    PickupDoneTime: enti.PickupDoneTime,

    PackageList: packagesL,

    InvoiceData: ShopeeInvoiceDataModel{
      Number: enti.InvoiceData.Number,
      SeriesNumber: enti.InvoiceData.SeriesNumber,
      AccessKey: enti.InvoiceData.AccessKey,
      IssueDate: enti.InvoiceData.IssueDate,
      TotalValue: enti.InvoiceData.TotalValue,
      ProductTotalValue: enti.InvoiceData.ProductTotalValue,
      TaxCode: enti.InvoiceData.TaxCode,
    },

    CheckoutShippingCarrier: enti.CheckoutShippingCarrier,
    ReverseShippingFee: enti.ReverseShippingFee,
    OrderChargeableWeightGram: enti.OrderChargeableWeightGram,
    PrescriptionImages: enti.PrescriptionImages,
    PrescriptionCheckStatus: enti.PrescriptionCheckStatus,
    AdvancePackage: enti.AdvancePackage,
    ReturnRequestDueDate: enti.ReturnRequestDueDate ,

    CreatedAt: enti.CreatedAt,
    CreatedBy: enti.CreatedBy,
    UpdatedAt: enti.UpdatedAt,
    UpdatedBy: enti.UpdatedBy,
  }
}
// [Method End]
// ----------------- [Entity] - End.Collection("shopee_order")   ----------------

// Template
// ----------------- [DTO/Entity/Model] - Start.Collection("Shop_rder") ----------------
// [Concept] : xxxx
// [Compnent Struct.Start]
// [Component Struct.End]
// [Core Struct.Start]
// [Core Struct.Emd]
// ----------------- [DTO/Entity/Model] - End.Collection("Shop_rder")   ----------------
