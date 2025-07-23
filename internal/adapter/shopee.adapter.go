package adapter

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"
)

type IShopeeService interface {
	GetAccessToken(partnerID string, shopID string, code string, signCode string) (*IResShopeeAuthResponse, error)
	// ExchangeToken(ctx context.Context, code string, redirectURI string, partnerID string) (*ShopeeAuthResponse, error)
}

type shopeeApi struct {
	baseURL    string
	prefixURL  string
	logger     *zap.Logger
	httpClient *http.Client
}

type IResShopeeAuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpireIn     int    `json:"expire_in"`
	ShopID       int64  `json:"shop_id"`
	Error        string `json:"error"`
	Message      string `json:"message"`
}

// NewShopeeAPI initializes adapter
func NewShopeeAPI(baseURL string, prefix string, log *zap.Logger) IShopeeService {
	return &shopeeApi{
		baseURL:    baseURL,
		prefixURL:  prefix,
		logger:     log,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

type IBGetAccessToken struct {
	Code      string `json:"code"`
	PartnerID int64  `json:"partner_id"`
	ShopID    int64  `json:"shop_id"`
}

func (s *shopeeApi) GetAccessToken(partnerID string, shopID string, code string, signCode string) (*IResShopeeAuthResponse, error) {

	// seperate component
	timeStp := strconv.FormatInt(time.Now().Unix(), 10)
	url := fmt.Sprintf("%s%s/auth/token/get?partner_id=%s&timestamp=%s&sign=%s", s.baseURL, s.prefixURL, partnerID, timeStp, signCode)

	// body : code , partner_id, shop_id || main_account_id

	partnerIDInt, err := strconv.ParseInt(partnerID, 10, 64)
	if err != nil {
		return nil, errors.New("failed to convert partnerID to int64")
	}
	shopIDInt, err := strconv.ParseInt(shopID, 10, 64)
	if err != nil {
		return nil, errors.New("failed to convert shopID to int64")
	}

	payload := &IBGetAccessToken{
		Code:       code,
		PartnerID: partnerIDInt,
    ShopID:    shopIDInt}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(context.Background(), "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// s.logger.Debug("adapter.GetAccessToken: Outgoing Request Before",
	// 	zap.String("url", req.URL.String()),
	// 	zap.String("method", req.Method),
	// 	zap.Any("headers", req.Header),
	// 	zap.ByteString("body", body),
	// )

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// แสดง raw body ทั้งหมด

	// s.logger.Debug("adapter.GetAccessToken: Outgoing Request After",
	// 	zap.String("url", req.URL.String()),
	// 	zap.String("method", req.Method),
	// 	zap.Any("headers", req.Header),
	// 	zap.ByteString("body", bodyBytes),
	// )
	// s.logger.Sugar().Debugf("Shopee Raw Response Body: %s", zap.ByteString("response",bodyBytes))
	// s.logger.Sugar().Debugf("adapter.GetAccessToken : %s", resp.Body)

	var authResp IResShopeeAuthResponse
	// for debug
	if err := json.Unmarshal(bodyBytes, &authResp); err != nil {
		s.logger.Error("adapter.GetAccessToken : json.Unmarshal(bodyBytes ,&authResp) :", zap.Error(err))
		return nil, err
	}

	// for before debug
	// if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
	// 	return nil, err
	// }
	if authResp.Error != "" {
		return nil, errors.New(authResp.Message)
	}

	return &authResp, nil
}
