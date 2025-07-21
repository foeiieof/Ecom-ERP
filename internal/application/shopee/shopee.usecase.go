package shopee

import (
	"crypto/hmac"
	"crypto/sha256"
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
	GetAccessToken(shopID string) (string, error)
	GenerateAuthLink(partnerName string, partnerId string, partnerKey string) (string, error)
	WebhookAuthentication(partnerId string, code string, shopId string) (any, error)
}

type shopeeService struct {
	Logger               *zap.Logger
	ShopeeAuthRepository ShopeeAuthRepository
	Config               *env.Config
}

func NewShopeeService(repo ShopeeAuthRepository, logger *zap.Logger, cfg *env.Config) IShopeeService {
	return &shopeeService{
		Logger:               logger,
		ShopeeAuthRepository: repo,
		Config:               cfg,
	}
}

func (s *shopeeService) GetAccessToken(shopID string) (string, error) {
	return s.ShopeeAuthRepository.GetShopeeAuthByShopId(shopID)
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
	path := fmt.Sprintf("%s/shop/auth_partner",s.Config.Shopee.ShopeeApiBasePrefix )
  redirectUrl := fmt.Sprintf("http://%s%s%s/shopee/webhook/auth_partner/%s", s.Config.Server.Host, s.Config.Server.Port, s.Config.Server.Prefix,partnerId  )

  baseString := fmt.Sprintf("%s%s%s", partnerId, path, timest)
	h := hmac.New(sha256.New, []byte(partnerKey))
	h.Write([]byte(baseString))
	sign := hex.EncodeToString(h.Sum(nil))
	url := fmt.Sprintf("%s%s?partner_id=%s&timestamp=%s&sign=%s&redirect=%s",host,path, partnerId, timest, sign, redirectUrl)

	return url, nil
}

func (s *shopeeService) WebhookAuthentication(partnerId string, code string, shopId string) (any, error) {
	return map[string]string{"status": "ok", "partner_id": partnerId, "code": code, "shopId": shopId}, nil
}
