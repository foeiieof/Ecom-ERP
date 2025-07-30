package shopee

import (
	"ecommerce/internal/adapter/dto"
	"time"
)

type ShopeeAuthEntity struct {
	PartnerID     string
  ShopID        string
  // Code          string
	AccessToken   string
	RefreshToken  string
	ExpiredAt     time.Time

	// CreatedAt     time.Time
	// CreatedBy     string

	// MoidifiedAt   time.Time
	// ModifiedBy    string
}


type ShopeeShopListEntity struct {
    ShopList []dto.IResAuthedShopList 
} 

type ShopeeOrderListEntity struct {
    OrderList []dto.IResOrderList
}
