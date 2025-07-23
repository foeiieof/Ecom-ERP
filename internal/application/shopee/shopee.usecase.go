package shopee

import (
	"crypto/hmac"
	"crypto/sha256"
	"ecommerce/internal/adapter"
	"ecommerce/internal/env"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"

	"go.uber.org/zap"
)

// Nameing Service
// [Method or Action]:Details

type IShopeeService interface {
	GetAccessToken(shopID string) (*ShopeeAuthModel, error)
	GetAccessAndRefreshToken(partnerID string, shopID string, code string) (*IResAccessAndRefreshToken, error)

	GenerateAuthLink(partnerName string, partnerId string, partnerKey string) (string, error)
	GenerateSignWithPathURL(state string, pathUrl string, partnerID string, partnerKey string, shopID string, code string, accessToken string) (*IResGenerateSignWithUri, error)

	WebhookAuthentication(partnerId string, code string, shopId string) (any, error)

	AddShopeeAuthRequest(partnerId string, partnerKey string, partnerName string, url string) (*ShopeeAuthRequestModel, error)
	AddShopeePartner(partnerId string, partnerKey string, partnerName string) (*ShopeePartnerModel, error)
}

type shopeeService struct {
	Config *env.Config
	Logger *zap.Logger

	ShopeeAdapter adapter.IShopeeService

	ShopeeAuthRepository        ShopeeAuthRepository
	ShopeeAuthRequestRepository ShopeeAuthRequestRepository
	ShopeePartnerRepository     ShopeePartnerRepository
}

func NewShopeeService(cfg *env.Config, logger *zap.Logger, adapter adapter.IShopeeService,
	auth ShopeeAuthRepository,
	authReq ShopeeAuthRequestRepository,
	shopeePartner ShopeePartnerRepository,
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

func (s *shopeeService) GetAccessToken(shopID string) (*ShopeeAuthModel, error) {

  data,err := s.ShopeeAuthRepository.GetShopeeAuthByShopId(shopID)
  if err != nil {
    // s.Logger.Error("GetAccessToken : s.ShopeeAuthRepository.GetShopeeAuthByShopId error", zap.Error(err))
   return nil , err 
  }
  return  data, nil
}

func (s *shopeeService) GenerateAuthLink(partnerName string, partnerId string, partnerKey string) (string, error) {

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
  redirectUrl := fmt.Sprintf("https://ecom-webhook.vercel.app/api/v1/webhook/auth_partner/%s",partnerId)
  // redirectUrl := "https://google.com"
	baseString := fmt.Sprintf("%s%s%s", partnerId, path, timest)
	h := hmac.New(sha256.New, []byte(partnerKey))
	h.Write([]byte(baseString))
	sign := hex.EncodeToString(h.Sum(nil))
	url := fmt.Sprintf("%s%s?partner_id=%s&timestamp=%s&sign=%s&redirect=%s", host, path, partnerId, timest, sign, redirectUrl)

	return url, nil
}

type IResGenerateSignWithUri struct {
	Method    string
	Path      string
	Sign      string
	Code      string
	TimeStamp time.Time
}

func (s *shopeeService) GenerateSignWithPathURL(state string, pathUrl string, partnerID string, partnerKey string, shopID string, code string, accessToken string) (*IResGenerateSignWithUri, error) {
	// var url string
	var method string
	// host := s.Config.Shopee.ShopeeApiBaseUrl
	timest := strconv.FormatInt(time.Now().Unix(), 10)
	path := fmt.Sprintf("%s%s", s.Config.Shopee.ShopeeApiBasePrefix, pathUrl)
	// s.Logger.Sugar().Debugf("adapter.shopee.GenerateSignWithPathURL: %s", path)

	var baseString string
	// baseString := fmt.Sprintf("%s%s%s", partnerID, path, timest)
	switch state {
	case "PUBLIC":
		// For Public APIs: partner_id, api path, timestamp
		baseString = fmt.Sprintf("%s%s%s", partnerID, path, timest)
	case "SHOP":
		// For Shop APIs: partner_id, api path, timestamp, access_token, shop_id
		baseString = fmt.Sprintf("%s%s%s%s%s", partnerID, path, timest, accessToken, shopID)

	case "MERCHANT":
		// Not available
		// For Merchant APIs: partner_id, api path, timestamp, access_token, merchant_id
		merchantID := ""
		baseString = fmt.Sprintf("%s%s%s%s%s", partnerID, path, timest, accessToken, merchantID)
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
		// url = fmt.Sprintf("%s%s?partner_id=%s&timestamp=%s&sign=%s", host, path, partnerID, timest, sign)

	default:
		s.Logger.Error("adapter.shopee.GenerateSignWithPathURL:invalid path")
		return nil, errors.New("adapter.shopee.GenerateSignWithPathURL: invalid path")
	}

	return &IResGenerateSignWithUri{
		Method:    method,
		Path:      path,
		Sign:      sign,
		Code:      code,
		TimeStamp: time.Now()}, nil
}

func (s *shopeeService) WebhookAuthentication(partnerId string, code string, shopId string) (any, error) {
	return map[string]string{"status": "ok", "partner_id": partnerId, "code": code, "shopId": shopId}, nil
}

func (s *shopeeService) AddShopeeAuthRequest(partnerId string, partnerKey string, partnerName string, url string) (*ShopeeAuthRequestModel, error) {
	data, err := s.ShopeeAuthRequestRepository.SaveShopeeAuthRequestWithName(partnerId, partnerKey, partnerName, url)
	if err != nil {
		return nil, errors.New("failed to insert shopee auth request")
	}
	return data, nil
}

func (s *shopeeService) AddShopeePartner(partnerId string, partnerKey string, partnerName string) (*ShopeePartnerModel, error) {
	data, err := s.ShopeePartnerRepository.CreateShopeePartner(partnerId, partnerKey, partnerName)
	if err != nil {
		return nil, errors.New("failed to insert shopee partner")
	}
	return data, nil
}

type IResAccessAndRefreshToken struct {
	AccessToken  string
	RefreshToken string
	ExpiredAt    time.Time
}

func (s *shopeeService) GetAccessAndRefreshToken(partnerID string, shopID string, code string) (*IResAccessAndRefreshToken, error) {


  // Get partner key
  partner, err := s.ShopeePartnerRepository.GetShopeePartnerByPartnerId(partnerID)
  if err != nil {
    return nil, errors.New(err.Error())
  }

	dataGen, err := s.GenerateSignWithPathURL("PUBLIC", "/auth/token/get", partner.PartnerID, partner.PartnerKey, shopID, code, "")
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

	resDB, error := s.ShopeeAuthRepository.CreateShopeeAuth(shopID, code,resApi.AccessToken, resApi.RefreshToken)
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
