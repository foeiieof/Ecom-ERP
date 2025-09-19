package shopee

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"ecommerce/internal/adapter"
	"ecommerce/internal/adapter/dto"
	"ecommerce/internal/application/shopee/partner"
	"ecommerce/internal/env"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

// Nameing Service
// [Method or Action]:Details

type IShopeeService interface {
	GetAccessTokenByShopID(ctx context.Context,shopID string) (*ShopeeAuthEntity, error)
	GetRefreshTokenOnAdapter(ctx context.Context,partnerID string, shopID string, refreshToken string) (*ShopeeAuthEntity, error)
	CreateAccessAndRefreshTokenByCodeOnAdapter(ctx context.Context,partnerID string, shopID string, code string) (*IResAccessAndRefreshToken, error)

	GenerateAuthLink(ctx context.Context,partnerName string, partnerId string, partnerKey string) (string, error)
	// GenerateSignWithPathURL(ctx context.Context,state string, pathUrl string, partnerID string, partnerKey string, shopID string, code string, accessToken string) (*IGenerateSignWithUri, error)

	WebhookAuthentication(ctx context.Context,partnerId string, code string, shopId string) (any, error)

	AddShopeeAuthRequest(ctx context.Context,partnerId string, partnerKey string, partnerName string, url string) (*ShopeeAuthRequestModel, error)
	// AddShopeePartner(ctx context.Context,partnerId string, partnerKey string, partnerName string) (*ShopeePartnerModel, error)
  GetShopeeShopDetailsByShopID(ctx context.Context, user string ,shopID string) ( *ShopeeShopDetailsEntityDTO,error)
	GetShopeeShopListByPartnerID(ctx context.Context,partnerID string) (*[]IResShopeeShopList, error)

	// order
	GetShopeeOrderListByShopID(ctx context.Context,shopID string, timeType string, timeFrom string, timeTo string, status string, page string, size string) (*ShopeeOrderListEntity, error)
  GetShopeeOrderDetailByOrderSN(ctx context.Context,shopID string,orderSN string, pending string, option string) (*ShopeeOrderListWithDetailEntity, error)

}

type shopeeService struct {
	Config *env.Config
	Logger *zap.Logger

	ShopeeAdapter adapter.IShopeeService

	ShopeeAuthRepository        ShopeeAuthRepository // for Collect shopee shop may contains (access token , refresh token , ...other)
	ShopeeAuthRequestRepository ShopeeAuthRequestRepository
	ShopeePartnerRepository     partner.ShopeePartnerRepository
  ShopeeShopDetailsRepository ShopeeShopDetailsRepository
  ShopeeOrderRepository       ShopeeOrderRepository
}

func NewShopeeService(cfg *env.Config, logger *zap.Logger, adapter adapter.IShopeeService,
	auth    ShopeeAuthRepository,
	authReq ShopeeAuthRequestRepository,
	shopeePartner partner.ShopeePartnerRepository,
  shopeeShop  ShopeeShopDetailsRepository,
  shopeeOrder ShopeeOrderRepository,
) IShopeeService {
	return &shopeeService{
		Config:                      cfg,
		Logger:                      logger,
		ShopeeAdapter:               adapter,
		ShopeeAuthRepository:        auth,
		ShopeeAuthRequestRepository: authReq,
		ShopeePartnerRepository:     shopeePartner,
    ShopeeShopDetailsRepository: shopeeShop,
    ShopeeOrderRepository:       shopeeOrder,
	}
}

func (s *shopeeService) GetAccessTokenByShopID(ctx context.Context,shopID string) (*ShopeeAuthEntity, error) {
	data, err := s.ShopeeAuthRepository.GetShopeeShopAuthByShopId(shopID)
	if err != nil {

		return nil, err
	}
  // s.Logger.Debug("usecase.GetAccessTokenByShopID", zap.Any("data", data))
	// s.Logger.Debug("GetAccessTokenByShopID", zap.Any("data", data))

	// use GetRefreshTokenOnAdapter when ExpiredAt is over  time

	expTime := data.ExpiredAt
	currentTime := time.Now()
	diffTime := expTime.Sub(currentTime)

  // s.Logger.Debug("diffTime", zap.Any("diffTime", diffTime.Minutes()))

	if diffTime.Minutes() < 2 {
		// GetNew Accessstoken with adapter
    // s.Logger.Debug("diffTime", zap.Any("diffTime", diffTime.Minutes()))
		accessToken, err := s.GetRefreshTokenOnAdapter(ctx,data.PartnerID, data.ShopID, data.RefreshToken)
		if err != nil {
			return nil, err
		}
    s.Logger.Debug("accessToken", zap.Any("accessToken", accessToken))
		// s.Logger.Debug("diffTime", zap.Any("diffTime", diffTime.Minutes()))
		return accessToken, nil
	}

	return &ShopeeAuthEntity{
		PartnerID:    data.PartnerID,
		ShopID:       data.ShopID,
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
		ExpiredAt:    data.ExpiredAt,
	}, nil
}

func (s *shopeeService) GetRefreshTokenOnAdapter(ctx context.Context,partnerID string, shopID string, refreshToken string) (*ShopeeAuthEntity, error) {
	// gen log from req

	partnerData, err := s.ShopeePartnerRepository.GetShopeePartnerByID(ctx,partnerID)
	if err != nil {
		return nil, errors.New("usecase.GetRefreshTokenOnAdapter : Partner_ID not found")
	}

	dataGen, err := s.GenerateSignWithPathURL(ctx,"PUBLIC", "/auth/access_token/get", partnerData.PartnerID, partnerData.SecretKey, shopID, "", "")
	if err != nil {
		s.Logger.Error("usecase.GetAccessAndRefreshToken : s.GenerateSignWithPathURL error", zap.Error(err))
		return nil, errors.New(err.Error())
	}
  // s.Logger.Debug("dataGen", zap.Any("dataGen", dataGen))

	res, err := s.ShopeeAdapter.GetRefreshToken(partnerID, shopID, refreshToken, dataGen.Sign)
	if err != nil {
		return nil, err
	}
  s.Logger.Debug("res", zap.Any("res", res))

	// Create log_refresh_token

	updated, err := s.ShopeeAuthRepository.UpdateShopeeShopAuth(partnerID, "",shopID, res.AccessToken, res.RefreshToken)
	if err != nil {
		return nil, err
	}
	// update refresh token
	// req to access token
// update access token
	// updateShopeeShopAuth.AccessToken, nil
	return &ShopeeAuthEntity{
		PartnerID:    updated.PartnerID,
		ShopID:       updated.ShopID,
		AccessToken:  updated.AccessToken,
		RefreshToken: updated.RefreshToken,
		ExpiredAt:    updated.ExpiredAt,
	}, nil
}

func (s *shopeeService) GenerateAuthLink(ctx context.Context,partnerName string, partnerId string, partnerKey string) (string, error) {

	if partnerName == "" || partnerId == "" || partnerKey == "" {
		return "", errors.New("partnerName or partnerId or partnerKey is required")
	}
	timest := strconv.FormatInt(time.Now().Unix(), 10)

	// host := "https://partner.test.shopeemobile.com"
	// path := "/api/v2/shop/auth_partner"
	// redirectUrl := "https://www.baidu.com/"

	host := s.Config.Shopee.ShopeeApiBaseUrl
	path := fmt.Sprintf("%s/shop/auth_partner", s.Config.Shopee.ShopeeApiBasePrefix)
	// redirectUrl := fmt.Sprintf("https://%s%s%s/shopee/webhook/auth_partner/%s", s.Config.Server.Host, s.Config.Server.Port, s.Config.Server.Prefix, partnerId)
	redirectUrl := fmt.Sprintf("https://ecom-webhook.vercel.app/api/v1/webhook/auth_partner/%s", partnerId)
	// redirectUrl := "https://google.com"
	baseString := fmt.Sprintf("%s%s%s", partnerId, path, timest)
	h := hmac.New(sha256.New, []byte(partnerKey))
	h.Write([]byte(baseString))
	sign := hex.EncodeToString(h.Sum(nil))
	url := fmt.Sprintf("%s%s?partner_id=%s&timestamp=%s&sign=%s&redirect=%s", host, path, partnerId, timest, sign, redirectUrl)

	return url, nil
}

type IGenerateSignWithUri struct {
	Method    string
	Path      string
	Sign      string
	Code      string
	TimeStamp time.Time
}

func (s *shopeeService) GenerateSignWithPathURL(ctx context.Context,state string, pathUrl string, partnerID string, partnerKey string, shopID string, code string, accessToken string) (*IGenerateSignWithUri, error) {
	// var url string
	var method string
	// host := s.Config.Shopee.ShopeeApiBaseUrl
	timest := strconv.FormatInt(time.Now().Unix(), 10)
	path := fmt.Sprintf("%s%s", s.Config.Shopee.ShopeeApiBasePrefix, pathUrl)
	// s.Logger.Sugar().Debugf("adapter.shopee.GenerateSignWithPathURL: %s", path)
var baseString string // baseString := fmt.Sprintf("%s%s%s", partnerID, path, timest)
	switch state {
	case "PUBLIC":
		// For Public APIs: partner_id, api path, timestamp
		baseString = fmt.Sprintf("%s%s%s", partnerID, path, timest)
		// break  // redundant
	case "SHOP":
		// For Shop APIs: partner_id, api path, timestamp, access_token, shop_id
		baseString = fmt.Sprintf("%s%s%s%s%s", partnerID, path, timest, accessToken, shopID)
		// break

	case "MERCHANT":
		// Not available
		// For Merchant APIs: partner_id, api path, timestamp, access_token, merchant_id
		merchantID := ""
		baseString = fmt.Sprintf("%s%s%s%s%s", partnerID, path, timest, accessToken, merchantID)
		// break
	default:
		s.Logger.Error("adapter.shopee.GenerateSignWithPathURL: invalid state")
		return nil, errors.New("adapter.shopee.GenerateSignWithPathURL: invalid state")
	}

	// s.Logger.Sugar().Debugf("adapter.shopee.GenerateSignWithPathURL: baseString - %s", baseString)

	h := hmac.New(sha256.New, []byte(partnerKey))
	h.Write([]byte(baseString))
	sign := hex.EncodeToString(h.Sum(nil))

	switch path {

	case "/api/v2/auth/token/get":
		method = "GET"
		// break

	case "/api/v2/auth/access_token/get":
		method = "POST"
		// break

	case "/api/v2/public/get_shops_by_partner":
		method = "GET"
	// break
	// url = fmt.Sprintf("%s%s?partner_id=%s&timestamp=%s&sign=%s", host, path, partnerID, timest, sign)
	case "/api/v2/order/get_order_list":
		method = "GET"

	default:
		s.Logger.Error("adapter.shopee.GenerateSignWithPathURL:invalid path")
		return nil, errors.New("adapter.shopee.GenerateSignWithPathURL: invalid path")
	}

	return &IGenerateSignWithUri{
		Method:    method,
		Path:      path,
		Sign:      sign,
		Code:      code,
		TimeStamp: time.Now()}, nil
}

func (s *shopeeService) WebhookAuthentication(ctx  context.Context,partnerId string, code string, shopId string) (any, error) {


  // 0. get partner details
  partner, err := s.ShopeePartnerRepository.GetShopeePartnerByID(ctx, partnerId)
  if err != nil { return nil , err}

  // 1. get sign code
  // "/api/v2/auth/token/get"
  sign,err := s.ShopeeAdapter.GenerateSignWithPathURL("PUBLIC", "/auth/token/get", partner.PartnerID ,partner.SecretKey,shopId,code, "" )
  if err != nil { return nil, err}
  s.Logger.Info("usecase.WebhookAuthentication", zap.String("val" , partner.PartnerID))

  // 2. get access token 
  adapter, err := s.ShopeeAdapter.GetAccessToken(partner.PartnerID,shopId, code, sign.Sign)
  if err != nil { return nil ,err }
  // s.Logger.Info("shopee.usecase.WebhookAuthentication", zap.String("val","xxxxxxxxxxxxxxxxxxx" ))


  // 3. save to db --> ShopeeShopAuthRepositoryo
  _,err = s.ShopeeAuthRepository.UpdateShopeeShopAuth(partner.PartnerID, code,shopId,adapter.AccessToken,adapter.RefreshToken  ) 
  if err != nil { return nil, err } 

  

  return map[string]string{"status": "ok", "partner_id": partnerId, "code": code, "shopId": shopId, "access_token": adapter.AccessToken, "refresh_token": adapter.RefreshToken}, nil
}

func (s *shopeeService) AddShopeeAuthRequest(ctx context.Context,partnerId string, partnerKey string, partnerName string, url string) (*ShopeeAuthRequestModel, error) {
	data, err := s.ShopeeAuthRequestRepository.SaveShopeeAuthRequestWithName(partnerId, partnerKey, partnerName, url)
	if err != nil {
		return nil, errors.New("failed to insert shopee auth request")
	}
	return data, nil
}

// func (s *shopeeService) AddShopeePartner(ctx context.Context,partnerId string, partnerKey string, partnerName string) (*ShopeePartnerModel, error) {
// 	data, err := s.ShopeePartnerRepository.CreateShopeePartner(ctx,partnerId, )
// 	if err != nil {
// 		return nil, errors.New("failed to insert shopee partner")
// 	}
// 	return data, nil
// }

type IResAccessAndRefreshToken struct {
	AccessToken  string
	RefreshToken string
	ExpiredAt    time.Time
}

func (s *shopeeService) CreateAccessAndRefreshTokenByCodeOnAdapter(ctx context.Context,partnerID string, shopID string, code string) (*IResAccessAndRefreshToken, error) {

	// Get partner key
	partner, err := s.ShopeePartnerRepository.GetShopeePartnerByID(ctx ,partnerID)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	dataGen, err := s.GenerateSignWithPathURL(ctx,"PUBLIC", "/auth/token/get", partner.PartnerID, partner.SecretKey, shopID, code, "")
	if err != nil {
		s.Logger.Error("usecase.GetAccessAndRefreshToken : s.GenerateSignWithPathURL error", zap.Error(err))
		return nil, errors.New(err.Error())
	}

	// s.Logger.Debug("usecase.GetAccessAndRefreshToken : dataGen", zap.Any("dataGen", dataGen))

	resApi, err := s.ShopeeAdapter.GetAccessToken(partnerID, shopID, dataGen.Code, dataGen.Sign)
	if err != nil {
		s.Logger.Error("usecase.GetAccessAndRefreshToken : s.ShopeeAdapter.GetAccessToken error", zap.Error(err))
		return nil, errors.New(err.Error())
	}

	resDB, error := s.ShopeeAuthRepository.CreateShopeeAuth(partnerID, shopID, code, resApi.AccessToken, resApi.RefreshToken)
	if error != nil {
		s.Logger.Error("usecase.GetAccessAndRefreshToken : s.ShopeeAuthRepository.CreateShopeeAuth error", zap.Error(error))
		return nil, errors.New(error.Error())
	}

	return &IResAccessAndRefreshToken{
		AccessToken:  resDB.AccessToken,
		RefreshToken: resDB.RefreshToken,
		ExpiredAt:    resDB.ExpiredAt}, nil

	// // --layer ---
	// // : apllication
	// // : external : shopee api

	// if partnerName == "" || partnerId == "" || partnerKey == "" {
	// 	return "", errors.New("partnerName or partnerId or partnerKey is required")
	// }

	// timest := strconv.FormatInt(time.Now().Unix(), 10)

	// // host := "https://partner.test.shopeemobile.com"
	// // path := "/api/v2/shop/auth_partner"
	// // redirectUrl := "https://www.baidu.com/"

	// host := s.Config.Shopee.ShopeeApiBaseUrl
	// path := fmt.Sprintf("%s/shop/auth_partner", s.Config.Shopee.ShopeeApiBasePrefix)
	// redirectUrl := fmt.Sprintf("https://%s%s%s/shopee/webhook/auth_partner/%s", s.Config.Server.Host, s.Config.Server.Port, s.Config.Server.Prefix, partnerId)

	// baseString := fmt.Sprintf("%s%s%s", partnerId, path, timest)
	// h := hmac.New(sha256.New, []byte(partnerKey))
	// h.Write([]byte(baseString))
	// sign := hex.EncodeToString(h.Sum(nil))
	// url := fmt.Sprintf("%s%s?partner_id=%s&timestamp=%s&sign=%s&redirect=%s", host, path, partnerId, timest, sign, redirectUrl)
}


func (s *shopeeService)GetShopeeShopDetailsByShopID(ctx context.Context, user string,shopID string) ( *ShopeeShopDetailsEntityDTO,error) {
  // 0. check in db
  shop,err := s.ShopeeAuthRepository.GetShopeeShopAuthByShopId(shopID)
  if err != nil { return nil ,err}
  partner,err := s.ShopeePartnerRepository.GetShopeePartnerByID(ctx,shop.PartnerID)
  if err != nil { return nil,err}

  // 1. get from adapter 
  params := &adapter.IReqShopeeAdapter{ PartnerID: partner.PartnerID, AccessToken: shop.AccessToken, ShopID: shopID, SecretKey: partner.SecretKey, }
  
  dto, err := s.ShopeeShopDetailsRepository.GetShopeeShopDetailsByShopID(ctx, shop.ShopID)
  if err != nil {
    s.Logger.Info("usecase.GetShopeeShopDetailsByShopID", zap.String("val","fetch on adapter"))
    // then dto -> Entity
    // refactoring : go routine 
    obj,err := s.ShopeeAdapter.GetShopProfile(ctx, params)
    if err != nil { return nil, err}

    objOpts, err := s.ShopeeAdapter.GetShopInfo(ctx, params)
    if err != nil { return nil, err}
    // 2. Create obj
    // not found then create new one
    // 4. save to DB

    var sifAffs []SipAffiShops_Struct
    var linkShops []LinkedDirectShopList_Struct
    var outletShops []OutletShopInfoList_Struct

    if len(objOpts.SipAffiShops) > 0 {
      for _,o := range objOpts.SipAffiShops {
        sifAffs = append(sifAffs, SipAffiShops_Struct{
          AffiShopID: strconv.FormatInt(int64(o.AffiShopID),10),
          Region: o.Region, 
        })
      }
    }

    if len(objOpts.LinkedDirectShopList) > 0 {
      for _,o := range objOpts.LinkedDirectShopList {
        linkShops = append(linkShops, LinkedDirectShopList_Struct{
          DirectShopID: strconv.FormatInt(int64(o.DirectShopID),10),
          DirectShopRegion: o.DirectShopRegion,
        })
      }
    }

    if len(objOpts.OutletShopInfoList) > 0 {
      for _,o := range objOpts.OutletShopInfoList {
        outletShops = append(outletShops, OutletShopInfoList_Struct{
          OutletShopID: strconv.FormatInt(int64(o.OutletShopID),10),
        })
      }
    }

    objShopeeShop := ShopeeShopDetailsEntityDTO{
      ShopID: shop.ShopID,
      ShopLogo: obj.ShopLogo,
      Description: obj.Description,
      ShopName: obj.ShopName,
      InvoiceIssuer: obj.InvoiceIssuer,

      Region: objOpts.Region,
      Status: ShopeeShopStatusEnum(objOpts.Status),
      SipAffiShops: sifAffs,
      IsCB: objOpts.IsCB,
      IsSip: objOpts.IsSip,
      IsUpgradedCBSC: objOpts.ISUpgradedCBSC,
      MerchantID: strconv.FormatInt(int64(objOpts.MerchantID),10),
      ShopFullFilmentFlag: ShopFullFilmentFlagEnum(objOpts.ShopFullFilmentFlag),
      IsMainShop: objOpts.IsMainShop,
      IsDirectShop: objOpts.IsDirectShop,
      LinkedMainShopID: strconv.FormatInt(int64(objOpts.LinkedMainShopID),10),
      LinkedDirectShopList: linkShops,
      IsOneAwb: objOpts.IsOneAwb,
      IsMartShop: objOpts.IsMartShop,
      IsOutletShop: objOpts.IsOutletShop,
      MartShopID: strconv.FormatInt(int64(objOpts.MartShopID), 10),
      OutletShopInfoList: outletShops,

      CreatedAt: time.Now(),
      CreatedBy: user,
      UpdatedAt: time.Now(),
      UpdatedBy: user,
    }
    dto, err := s.ShopeeShopDetailsRepository.CreateShopeeShopDetails(ctx, &objShopeeShop)
    if err != nil { return nil, err}
    return dto, nil
  }

  // 5. onvert to DTO
  // bc: both sane  Entity<->DTO 
  return dto,nil
} 

type IResShopeeShopList struct {
	ShopID   string
	ExpireAt time.Time
	Region   string
	// SipAffiShopList []IResSipAffiShopList `json:"sip_affi_shop_list"`
}

func (s *shopeeService) GetShopeeShopListByPartnerID(ctx context.Context,partnerID string) (*[]IResShopeeShopList, error) {
	// Refac
	// Get partner
	partnerData, err := s.ShopeePartnerRepository.GetShopeePartnerByID(ctx,partnerID)
	if err != nil {
		s.Logger.Error("usecase.GetShopeeShopListByPartnerID : s.ShopeePartnerRepository.GetShopeePartnerByPartnerId error", zap.Error(err))
		return nil, err
	}

	genData, err := s.GenerateSignWithPathURL(ctx,"PUBLIC", "/public/get_shops_by_partner", partnerData.PartnerID, partnerData.SecretKey, "", "", "")

	if err != nil {
		s.Logger.Error("usecase.GetShopeeShopListByPartnerID : s.GenerateSignWithPathURL error", zap.Error(err))
		return nil, err
	}

	shopListData, err := s.ShopeeAdapter.GetShopByPartnerPublic(partnerData.PartnerID, genData.Sign)
	if err != nil {
		s.Logger.Error("usecase.GetShopeeShopListByPartnerID : s.ShopeeAdapter.GetShopByPartnerPublic error", zap.Error(err))
		return nil, err
	}

	s.Logger.Debug("usecase.GetShopeeShopListByPartnerID : shopListData", zap.Any("shopListData", shopListData.AuthedShopList))

	// DTO -> Entity
	// Get sign
	// Get data from adapter

	var data []IResShopeeShopList

	for _, v := range shopListData.AuthedShopList {
		expireInt, err := v.ExpireTime.Int64()
		if err != nil {
			return nil, err
		}

		data = append(data, IResShopeeShopList{
			ShopID:   v.ShopID.String(),
			ExpireAt: time.Unix(expireInt, 0),
			Region:   v.Region,
		})
	
  // add to --> DB (stored)
    _,err = s.ShopeeAuthRepository.CreateShopeeAuth(partnerID, string(v.ShopID), "", "", "")
    if err != nil { 
      s.Logger.Info("usecase.GetShopeeShopListByPartnerID", zap.String("info", "failed create ShopeeShopAuth "))
    }

  }

	return &data, nil
}

func (s *shopeeService) GetShopeeOrderListByShopID(ctx context.Context,shopID string, timeType string, timeFrom string, timeTo string, status string, page string, size string) (*ShopeeOrderListEntity, error) {
	// shopDataRepo, err := s.ShopeeAuthRepository.GetShopeeShopAuthByShopId(shopID)
	// if err != nil {
	// 	s.Logger.Error("usecase.GetShopeeOrderListByShopID : s.ShopeeAuthRepository.GetShopeeShopAuthByShopId error", zap.Error(err))
	// 	return nil, err
	// }


  // accessToken if expired then send refresh_token to update access token
  // !!dont delete marktime : 4/09/2025,10:24 !!
  // shopData,err := s.GetAccessTokenByShopID(ctx,shopID)
  // if err != nil {
  //   s.Logger.Error("usecase.GetShopeeOrderListByShopID : s.GetAccessTokenByShopID error", zap.Error(err))
  //   return nil, err
  // }

  // 0. check in db
  shopData,err := s.ShopeeAuthRepository.GetShopeeShopAuthByShopId(shopID)
  if err != nil { return nil ,err}
  partnerData,err := s.ShopeePartnerRepository.GetShopeePartnerByID(ctx,shopData.PartnerID)
  if err != nil { 
		s.Logger.Error("usecase.GetShopeeOrderListByShopID : s.ShopeePartnerRepository.GetShopeePartnerByPartnerId error", zap.Error(err))
    return nil,err
  }


	//  ----------- set concurrent
	genData, err := s.GenerateSignWithPathURL(ctx,"SHOP", "/order/get_order_list", partnerData.PartnerID, partnerData.SecretKey, shopData.ShopID, "", shopData.AccessToken)
	if err != nil {
		s.Logger.Error("usecase.GetShopeeOrderListByShopID : s.GenerateSignWithPathURL error", zap.Error(err))
		return nil, err
	}

	// Paesr to string
	var optsQuery dto.IOptionShopeeQuery
	if timeType == string(dto.UPDATE_TIME) {
		optsQuery.TimeRange = dto.UPDATE_TIME
	} else {
		optsQuery.TimeRange = dto.CREATE_TIME
	}

	// handle(time) : 12345667

	optsQuery.TimeFrom, err = strconv.ParseInt(timeFrom, 10, 64)
	if err != nil {
		s.Logger.Error("usecase.GetShopeeOrderListByShopID : optsQuery.TimeFrom error", zap.Error(err))
		return nil, err
	}

	optsQuery.TimeTo, err = strconv.ParseInt(timeTo, 10, 64)
	if err != nil {
		s.Logger.Error("usecase.GetShopeeOrderListByShopID : optsQuery.TimeTo error", zap.Error(err))
		return nil, err
	}

	sizeParam, err := strconv.ParseInt(size, 10, 32)
	optsQuery.PageSize = int32(sizeParam)
	if err != nil {
		s.Logger.Error("usecase.GetShopeeOrderListByShopID : optsQuery.PageSize error", zap.Error(err))
		return nil, err
	}

	orderData, err := s.ShopeeAdapter.GetOrderListByShopID(shopData.PartnerID, shopData.AccessToken, shopData.ShopID, genData.Sign, &optsQuery)
	if err != nil { return nil, err }
  s.Logger.Debug("orderData", zap.Any("orderData", orderData))

  // test section available to delete
  // orderSN := [2]string{"string", "string"}

  orderSN := []string{} // slices
  if len(orderData.Response.OrderList) > 0 {
    for _,o := range orderData.Response.OrderList {
      orderSN = append(orderSN, o.OrderSN)
    }
  }

  params := &adapter.IReqShopeeAdapter{
    PartnerID: partnerData.PartnerID,
    AccessToken: shopData.AccessToken,
    ShopID: shopData.ShopID,
    SecretKey: partnerData.SecretKey,
    OrderSN:orderSN,
  }

  // GetOrderDetails
  orderDetails,err := s.ShopeeAdapter.GetOrderDetailListByOrderSN(ctx, params)
  if err != nil { return nil,err}
  // s.Logger.Debug("usecase.GetShopOrderListByShopID", zap.Any("orderDetails", orderDetails))


  // create map for lookup in orderDetails\
  orderDetailsMap := make(map[string]dto.IResOrderListWithDetails, len(orderDetails))
  for _,d := range orderDetails {
    orderDetailsMap[d.OrderSN] = d
  }

  // Create new OrderDetails  
  var orderComps []ShopeeOrderEntity
  if len(orderSN) > 0 {
    for _,o := range orderSN {
      // declared 
      details := orderDetailsMap[o]

      // RecipientAddress 
      recipient := ShopeeRecipientAddressEntity{
            Name: details.RecipientAddress.Name,
            Phone: details.RecipientAddress.Phone,
            Town: details.RecipientAddress.Town,
            District: details.RecipientAddress.District,
            City: details.RecipientAddress.City,
            State: details.RecipientAddress.State,
            Region: details.RecipientAddress.Region,
            ZipCode: details.RecipientAddress.Zipcode,
            FullAddress: details.RecipientAddress.FullAddress,
      }
      // check item list
  
      var items []ShopeeItemListEntity
      for _,i := range details.ItemList {
        items = append(items,ShopeeItemListEntity{
          ItemID: strconv.FormatInt(i.ItemID, 10),
          ItemName: i.ItemName,
          ItemSKU: i.ItemSKU,
          ModelID: strconv.FormatInt(i.ModelID, 10),
          ModelName: i.ModelName,
          ModelSKU: i.ModelSKU,
          ModelQualityPurchased: i.ModelQtyPurchased,
          ModelOriginPrice: i.ModelOriginalPrice,
          ModelDiscountedPrice: i.ModelDiscountedPrice,
          WholeSale: i.Wholesale,
          Weight: i.Weight,
          AddOnDeal: i.AddOnDeal,
          PromotionType: ShopeePromotionTypeEnum(i.PromotionType),
          PromotionID: strconv.FormatInt(i.PromotionID,10),
          OrderItemID: strconv.FormatInt(i.OrderItemID,10),
          PromotionGroupID: strconv.FormatInt(int64(i.PromotionGroupID), 10),
          ImageInfo: ShopeeImageInfoEntity{ImageURL: i.ImageInfo.ImageURL},
          ProductLocationID: i.ProductLocationID,
          IsPrescriptionItem: i.IsPrescriptionItem,
          IsB2COwnedItem: i.IsB2COwnedItem,
        } )
      } 

      var packages []ShopeePackageListEntity
      for _,p :=  range details.PackageList {

        var itemInPackage []ShopeeItemListInPackageListEntity 
        for _, iip := range p.ItemList {
          itemInPackage = append(itemInPackage, ShopeeItemListInPackageListEntity{
            ItemID: strconv.FormatInt(iip.ItemID,10),
            ModelID:strconv.FormatInt(iip.ModelID,10) ,
            ModelQuantity: iip.ModelQuantity,
            OrderItemID: strconv.FormatInt(iip.OrderItemID,10),
            PromotionGroupID: strconv.FormatInt(int64(iip.PromotionGroupID),10) ,
            ProductLocationID: iip.ProductLocationID,
          })
        }

        packages = append(packages, ShopeePackageListEntity{
          PackageNumber: p.PackageNumber,
          LogisticsStatus: p.LogisticsStatus,
          LogisticsChannelID: strconv.FormatInt(p.LogisticsChannelID, 10),
          ShippingCarrier: p.ShippingCarrier,
          AllowSelfDesignAWB: p.AllowSelfDesignAWB,
          ItemList: itemInPackage,
          ParcelChargeableWeight: p.ParcelChargeableWeight,
          GroupShipmentID: p.SortingGroup,
        })
      }


      // check before  create entity
      order := ShopeeOrderEntity{
        OrderSN: o,
        OrderStatus: ShopeeOrderStatusEnum(details.OrderStatus),
        BookingSN: details.BookingSN,
        ShopeeOrderDetailsEntity: ShopeeOrderDetailsEntity{
          Region: details.Region,
          Currency: details.Currency,
          Cod: details.COD,
          TotalAmount: details.TotalAmount,
          PendingTerms: details.PendingTerms,

          ShippingCarrier: details.ShippingCarrier,
          PaymentMethod: details.PaymentMethod,
          EstimatedShippingFee: details.EstimatedShippingFee,
          MessageToSeller: details.MessageToSeller,
          CreateTime: time.Unix(details.CreateTime, 0),
          UpdateTime: time.Unix(details.UpdateTime, 0),
          DaysToShip: details.DaysToShip,
          ShipByDate: int(details.ShipByDate),
          BuyerUserId: strconv.FormatInt(int64(details.BuyerUserID),10),
          BuyerUsername: details.BuyerUsername,
          RecipientAddress: recipient,
          ActualShippingFee: details.ActualShippingFee,
          GoodsToDeclare: details.GoodsToDeclare,
          Note: details.Note,
          NoteUpdateTime: time.Unix(details.NoteUpdateTime,0),

          ItemList: items,

          PayTime: time.Unix(details.PayTime,0),
          DropShipper: details.Dropshipper,
          DropShipperPhone: details.DropshipperPhone,
          SplitUp: details.SplitUp,
          BuyerCancelReason: details.BuyerCancelReason,
          CancelBy: details.CancelBy,
          CancelReason: details.CancelReason,
          ActualShippingFeeConfirmed: details.ActualShippingFeeConfirmed,

          BuyerCPFID: details.BuyerCPFID,
          FulFillmentFlag: ShopeeFulfillmentFlagEnum(details.FulfillmentFlag),
          PickupDoneTime: time.Unix(details.PickupDoneTime,0),
          PackageList:packages ,
          InvoiceData: ShopeeInvoiceDataEntity{
            Number: details.InvoiceData.Number,
            SeriesNumber: details.InvoiceData.SeriesNumber,
            AccessKey: details.InvoiceData.AccessKey,
            IssueDate: time.Unix(details.InvoiceData.IssueDate,0),
            TotalValue: details.InvoiceData.TotalValue,
            ProductTotalValue: details.InvoiceData.ProductsTotalValue,
            TaxCode: details.InvoiceData.TaxCode,
          },

          CheckoutShippingCarrier: details.CheckoutShippingCarrier,
          ReverseShippingFee: details.ReverseShippingFee,
          OrderChargeableWeightGram: details.OrderChargeableWeight,
          PrescriptionImages: details.PrescriptionImages,
          PrescriptionCheckStatus: ShopeePrescriptionCheckStatusEnum(details.PrescriptionStatus),
          BookingSN: details.BookingSN,
          AdvancePackage: details.AdvancePackage,
          ReturnRequestDueDate: time.Unix (details.ReturnRequestDueDate, 0),
        },

      }
      orderComps = append(orderComps, order)
    }
  }

  
  // s.Logger.Debug("usecase.GetShopeeOrderListByShopID", zap.String("orderComps", strconv.FormatInt(int64(len(orderComps)), 10) ))

  if len(orderComps) > 0 {
    var newOrders []ShopeeOrderEntity

    for _, order  := range orderComps {
      if shopID != "" {
      order.ShopID = shopID
      }
      // check before save to db
      res,err := s.ShopeeOrderRepository.GetShopeeOrderByOrderSN(ctx, order.OrderSN)
      if err != nil {
      // save to DB with loop
        res,err = s.ShopeeOrderRepository.CrateShopeeOrderWithDetails(ctx, &order) 
        if err != nil {
          s.Logger.Info("usecase.GetShopeeOrderListByShopID", zap.String("saveOrder", err.Error() )) 
          continue
        }
      }
      newOrders = append(newOrders, *res )
    }

    if len(newOrders) > 0 {
      orderComps = newOrders
    }
  }


  // s.Logger.Debug("usecase.GetShopeeOrderListByShopID", zap.Any("orderComp", orderComps))

	orderList := &ShopeeOrderListEntity{OrderList: orderComps}

	// return orderData, nil
	return orderList , nil
}



// *** waiting for test
func (s *shopeeService) GetShopeeOrderDetailByOrderSN(ctx context.Context,shopID string,orderSN string, pending string, option string) (*ShopeeOrderListWithDetailEntity, error) {

  // accessToken
  // check shopID through GetAccessTokenByShopID 
 //  shopData,err := s.GetAccessTokenByShopID(ctx,shopID)
 //  if err != nil {
 //    s.Logger.Error("usecase.GetShopeeOrderDetailByOrderSN : s.GetAccessTokenByShopID error", zap.Error(err))
 //    return nil, err
 //  }

	// partnerData, err := s.ShopeePartnerRepository.GetShopeePartnerByID(ctx,shopData.PartnerID)
	// if err != nil {
	// 	s.Logger.Error("usecase.GetShopeeOrderDetailByOrderSN : s.ShopeePartnerRepository.GetShopeePartnerByPartnerId error", zap.Error(err))
	// 	return nil, err
	// }


  // 0. check in db
  shopData,err := s.ShopeeAuthRepository.GetShopeeShopAuthByShopId(shopID)
  if err != nil { return nil ,err}
  partnerData,err := s.ShopeePartnerRepository.GetShopeePartnerByID(ctx,shopData.PartnerID)
  if err != nil { 
		s.Logger.Error("usecase.GetShopeeOrderListByShopID : s.ShopeePartnerRepository.GetShopeePartnerByPartnerId error", zap.Error(err))
    return nil,err
  }
  // ------- 
  orderSNList := strings.Split(orderSN, ",")
  if len(orderSNList) < 1 {
    return nil, errors.New("orderSN is required")
  }

  var pendingOpts bool
  pendingParse,err := strconv.ParseBool(pending)
  if err != nil { pendingParse = false }
  pendingOpts = pendingParse

  var optionOpts bool
  optionParse,err := strconv.ParseBool(option)
  if err != nil { optionParse = false }
  optionOpts = optionParse

  orderDetailData, err := s.ShopeeAdapter.GetOrderDetailByOrderSN(
    partnerData.PartnerID, partnerData.SecretKey, shopData.AccessToken , shopData.ShopID, orderSNList, pendingOpts, optionOpts)
  if err != nil { return nil, err }

  // s.Logger.Debug("orderDetailData", zap.Any("orderDetailData", orderDetailData))

  orderListWithDetail := &ShopeeOrderListWithDetailEntity{ OrderList: orderDetailData.Response.OrderList, }

  return orderListWithDetail, nil
}
