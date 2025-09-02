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
type ShopeeShopProfileModel struct {
  ID          bson.ObjectID `bson:"_id"`
  ShopID      string `bson:"shop_id"`
  ShopLogo    string `bson:"shop_logo"`
  Description string `bson:"description"`
  ShopName    string `bson:"invoice_issuer"` 
}
// -- ShopeeShopDetails Collections
type SipAffiShopsModel struct {
  ID   bson.ObjectID `bson:"_id"`
  AffiShopID string   `bson:"affi_shop_id"`
  Region string   `bson:"region"`
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

type LinkedDirectShopListModel struct {
  DirectShopID string `bson:"direct_shop_id"`
  DirectShopRegion string `bson:"direct_shop_region"`
} 

type OutletShopInfoListModel struct {
  OutletShopID string `bson:"outlet_shop_id"`
}
// Core profile 
type ShopeeShopDetailsModel struct {
  // ID    bson.ObjectID  `bson:"_id"`
  // ShopName string   `bson:"shop_name"`
  ShopeeShopProfileModel
  Region string   `bson:"region"`
  Status string   `bson:"status"`
  SipAffiShops []SipAffiShopsModel `bson:"sip_affi_shops"`
  IsCB  bool `bson:"is_cb"`
  IsSip bool `bson:"is_sip"`
  ISUpgradedCBSC bool `bson:"is_upgraded_cbsc"`
  MerchantID string `bson:"merchant_id"`
  ShopFullFilmentFlag ShopFullFilmentFlagEnum `bson:"shop_fulfillment_flag"`
  IsMainShop bool `bson:"is_main_shop"`
  IsDirectShop bool `bson:"is_direct_shop"`
  LinkedMainShopID string `bson:"linked_main_shop_id"`
  LinkedDirectShopList []LinkedDirectShopListModel `bson:"linked_direct_shop_list"`
  IsOneAwb bool `bson:"is_one_awb"`
  IsMartShop bool `bson:"is_mart_shop"`
  IsOutletShop bool `bson:"is_outlet_shop"`
  MartShopID string `bson:"mart_shop_id"`
  OutletShopInfoList OutletShopInfoListModel `bson:"outlet_shop_info_list"`
}

// ---------------------------------------------- Demo Template ------------------------------------
// type DemoModel Struct {
//   ... 
//   CreatedAt time.Time   `bson:"created_at"`
// 	CreatedBy string      `bson:"created_by"`
// 	MoidifiedAt time.Time `bson:"modified_at"`
// 	ModifiedBy  string    `bson:"modified_by"`
// }
// ---------------------------------------------- End Demo Template --------------------------------
