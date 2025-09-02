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
  GetShopeeShopDetailsByShopID(ctx context.Context, shopID string) ( *ShopeeShopDetailsDTO,error)
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
}

func NewShopeeService(cfg *env.Config, logger *zap.Logger, adapter adapter.IShopeeService,
	auth ShopeeAuthRepository,
	authReq ShopeeAuthRequestRepository,
	shopeePartner partner.ShopeePartnerRepository,
) IShopeeService {
	return &shopeeService{
		Config:                      cfg,
		Logger:                      logger,
		ShopeeAdapter:               adapter,
		ShopeeAuthRepository:        auth,
		ShopeeAuthRequestRepository: authReq,
		ShopeePartnerRepository:     shopeePartner,
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

type ShopeeShopDetailsDTO struct {

} 

func (s *shopeeService)GetShopeeShopDetailsByShopID(ctx context.Context, shopID string) ( *ShopeeShopDetailsDTO,error) {

  // 0. check in db

  shop,err := s.ShopeeAuthRepository.GetShopeeShopAuthByShopId(shopID)
  if err != nil { return nil ,err}

  partner,err := s.ShopeePartnerRepository.GetShopeePartnerByID(ctx,shop.PartnerID)

  // 1. get from adapter 
  params := &adapter.IReqShopeeAdapter{ PartnerID: partner.PartnerID, AccessToken: shop.AccessToken, ShopID: shopID, SecretKey: partner.SecretKey, }
  _,err = s.ShopeeAdapter.GetShopProfile(ctx, params)

  // 2. save to DB


  return nil,nil
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

  // accessToken
  shopData,err := s.GetAccessTokenByShopID(ctx,shopID)
  if err != nil {
    s.Logger.Error("usecase.GetShopeeOrderListByShopID : s.GetAccessTokenByShopID error", zap.Error(err))
    return nil, err
  }


	partnerData, err := s.ShopeePartnerRepository.GetShopeePartnerByID(ctx,shopData.PartnerID)
	if err != nil {
		s.Logger.Error("usecase.GetShopeeOrderListByShopID : s.ShopeePartnerRepository.GetShopeePartnerByPartnerId error", zap.Error(err))
		return nil, err
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

	orderList := &ShopeeOrderListEntity{OrderList: orderData.Response.OrderList}

	// return orderData, nil
	return orderList, nil
}



// *** waiting for test
func (s *shopeeService) GetShopeeOrderDetailByOrderSN(ctx context.Context,shopID string,orderSN string, pending string, option string) (*ShopeeOrderListWithDetailEntity, error) {

  // accessToken
  // check shopID through GetAccessTokenByShopID 
  shopData,err := s.GetAccessTokenByShopID(ctx,shopID)
  if err != nil {
    s.Logger.Error("usecase.GetShopeeOrderListByShopID : s.GetAccessTokenByShopID error", zap.Error(err))
    return nil, err
  }

	partnerData, err := s.ShopeePartnerRepository.GetShopeePartnerByID(ctx,shopData.PartnerID)
	if err != nil {
		s.Logger.Error("usecase.GetShopeeOrderListByShopID : s.ShopeePartnerRepository.GetShopeePartnerByPartnerId error", zap.Error(err))
		return nil, err
	}

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

  s.Logger.Debug("orderDetailData", zap.Any("orderDetailData", orderDetailData))

  orderListWithDetail := &ShopeeOrderListWithDetailEntity{ OrderList: orderDetailData.Response.OrderList, }

  return orderListWithDetail, nil
}
