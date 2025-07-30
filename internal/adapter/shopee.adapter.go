package adapter

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	"ecommerce/internal/adapter/dto"
	"ecommerce/internal/env"
)


type IShopeeService interface {

	// waiting reface replace body abd query gen
	// GenerateBodyQueryParams()(,error)

	GenerateSignWithPathURL(state string, pathUrl string, partnerID string, partnerKey string, shopID string, code string, accessToken string) (*IResGenerateSignWithUri, error)

	GetAccessToken(partnerID string, shopID string, code string, signCode string) (*IResShopeeAuthResponse, error)
	GetRefreshToken(partnerID string, shopID string, refreshToken string, signCode string) (*dto.IResShopeeAuthRefreshResponse, error)
	// ExchangeToken(ctx context.Context, code string, redirectURI string, partnerID string) (*ShopeeAuthResponse, error)
	GetShopByPartnerPublic(partnerID string, signCode string) (*dto.IResGetShopByPartnerPublic, error)
	GetOrderListByShopID(partnerID string, accessToken string, shopID string, signCode string, optsShopee *dto.IOptionShopeeQuery) (*dto.IResGetOrderListByShopIDShop, error)
  GetOrderDetailByOrderSN(partnerID string, partnerKey string,accessToken string, shopID string, orderList []string, pending bool, option bool) (*dto.IResOrderDetailByOrderSN, error)
}

type shopeeApi struct {
  Config *env.Config

	BaseURL    string
	PrefixURL  string
	Logger     *zap.Logger
	HttpClient *http.Client
}

type IResShopeeAuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpireIn     int    `json:"expire_in"`
	ShopID       int64  `json:"shop_id"`
	Error        string `json:"error"`
	Message      string `json:"message"`
}

func NewShopeeAPI(config *env.Config, baseURL string, prefix string, log *zap.Logger) IShopeeService {
	return &shopeeApi{
    Config: config,
		BaseURL:    baseURL,
		PrefixURL:  prefix,
		Logger:     log,
		HttpClient: &http.Client{Timeout: 10 * time.Second},
	}
}


type IResGenerateSignWithUri struct {
	Method    string
	Path      string
	Sign      string
	Code      string
	TimeStamp time.Time
}

func (s *shopeeApi) GenerateSignWithPathURL(state string, pathUrl string, partnerID string, partnerKey string, shopID string, code string, accessToken string) (*IResGenerateSignWithUri, error) {
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

  case "/api/v2/order/get_order_detail":
      method = "GET"

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

// InterfaceBody : IBXXX
type IBGetAccessToken struct {
	Code      string `json:"code"`
	PartnerID int64  `json:"partner_id"`
	ShopID    int64  `json:"shop_id"`
}

func (s *shopeeApi) GetAccessToken(partnerID string, shopID string, code string, signCode string) (*IResShopeeAuthResponse, error) {

	// seperate component
	timeStp := strconv.FormatInt(time.Now().Unix(), 10)
	url := fmt.Sprintf("%s%s/auth/token/get?partner_id=%s&timestamp=%s&sign=%s", s.BaseURL, s.PrefixURL, partnerID, timeStp, signCode)

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
		Code:      code,
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

	resp, err := s.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	s.Logger.Debug("adapter.GetAccessToken: Outgoing Request After",
		zap.String("url", req.URL.String()),
		zap.String("method", req.Method),
		zap.Any("headers", req.Header),
		zap.ByteString("body", bodyBytes),
	)
	// s.logger.Sugar().Debugf("Shopee Raw Response Body: %s", zap.ByteString("response",bodyBytes))
	// s.logger.Sugar().Debugf("adapter.GetAccessToken : %s", resp.Body)

	var authResp IResShopeeAuthResponse
	// for debug
	if err := json.Unmarshal(bodyBytes, &authResp); err != nil {
		s.Logger.Error("adapter.GetAccessToken : json.Unmarshal(bodyBytes ,&authResp) :", zap.Error(err))
		return nil, err
	}

	// for before debug
	// if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
	// 	return nil, err
	// }
	// if authResp.Error != "" {
	// 	return nil, errors.New(authResp.Message)
	// }

	return &authResp, nil
}

func (s *shopeeApi) GetRefreshToken(partnerID string, shopID string, refreshToken string, signCode string) (*dto.IResShopeeAuthRefreshResponse, error) {

	// Public Api
	timeStp := strconv.FormatInt(time.Now().Unix(), 10)
	url := fmt.Sprintf("%s%s/auth/access_token/get?partner_id=%s&timestamp=%s&sign=%s", s.BaseURL, s.PrefixURL, partnerID, timeStp, signCode)

	partnerIDInt, err := func() (int32, error) {
		i64, err := strconv.ParseInt(partnerID, 10, 64)
		return int32(i64), err
	}()

	shopIDInt, err := func() (int32, error) {
		i64, err := strconv.ParseInt(shopID, 10, 64)
		return int32(i64), err
	}()

	payload := &dto.IBGetRefreshToken{
		RefreshToken: refreshToken,
		PartnerID:    int32(partnerIDInt),
		ShopID:       int32(shopIDInt),
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(context.Background(), "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	s.Logger.Debug("adapter.GetAccessToken: Outgoing Request Before",
		zap.String("url", req.URL.String()),
		zap.String("method", req.Method),
		zap.Any("headers", req.Header),
		zap.ByteString("body", body),
	)

	resp, err := s.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	s.Logger.Debug("adapter.GetRefreshToken.bodyBytes", zap.Any("bodyBytes", bodyBytes))

	var resShopeeAuthRefreshToken dto.IResShopeeAuthRefreshResponse
	if err := json.Unmarshal(bodyBytes, &resShopeeAuthRefreshToken); err != nil {
		s.Logger.Error("adapter.GetAccessToken : json.Unmarshal(bodyBytes ,&authRes) :", zap.Error(err))
		return nil, err
	}

	return &resShopeeAuthRefreshToken, nil
}

func (s *shopeeApi) GetShopByPartnerPublic(partnerID string, signCode string) (*dto.IResGetShopByPartnerPublic, error) {
	// https://partner.test-stable.shopeemobile.com/api/v2/public/get_shops_by_partner
	// query: partner_id, timestamp, sign

	// Public Api
	timeStp := strconv.FormatInt(time.Now().Unix(), 10)

	url := fmt.Sprintf("%s%s/public/get_shops_by_partner?partner_id=%s&timestamp=%s&sign=%s", s.BaseURL, s.PrefixURL, partnerID, timeStp, signCode)

	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	s.Logger.Debug("bodyBytes", zap.Any("bodyBytes", bodyBytes))

	var resGetShopByPartnerPublic dto.IResGetShopByPartnerPublic
	if err := json.Unmarshal(bodyBytes, &resGetShopByPartnerPublic); err != nil {
		s.Logger.Error("adapter.GetAccessToken : json.Unmarshal(bodyBytes ,&authRes) :", zap.Error(err))
		return nil, err
	}

	if resGetShopByPartnerPublic.Error != "" {
		return nil, errors.New(resGetShopByPartnerPublic.Message)
	}

	s.Logger.Debug("resGetShopByPartnerPublic", zap.Any("resGetShopByPartnerPublic", resGetShopByPartnerPublic))

	return &resGetShopByPartnerPublic, nil
}

func (s *shopeeApi) GetOrderListByShopID(partnerID string, accessToken string, shopID string, signCode string, optsShopee *dto.IOptionShopeeQuery) (*dto.IResGetOrderListByShopIDShop, error) {

	// var optsTimeRan dto.IEnumShopeeTimeRange = dto.CREATE_TIME
	// var OrderStatus dto.IEnumShopeeOrderStatus = dto.PROCESSED

	timeStp := strconv.FormatInt(time.Now().Unix(), 10)

	q := url.Values{}
	paramsQuery := map[string]string{
		// common
		"partner_id":   partnerID,
		"timestamp":    timeStp,
		"access_token": accessToken,
		"shop_id":      shopID,
		"sign":         signCode,
		// required
		// "time_range_field": string(optsShopee.TimeRange),
		// "time_from":        string(strconv.FormatInt(optsShopee.TimeFrom, 10)),
		// "time_to":          string(strconv.FormatInt(int64(optsShopee.TimeTo), 10)),
		// "page_size":        string(strconv.FormatInt(int64(optsShopee.PageSize), 10)),
		"time_range_field": "",
		"time_from":        "",
		"time_to":          "",
		"page_size":        "",
		// opts
		"cursor":                       string(optsShopee.CursorPage),
		"order_status":                 string(optsShopee.OrderStatus),
		"response_optional_fields":     string(optsShopee.ResponseOptionsField),
		"request_order_status_pending": string(strconv.FormatBool(optsShopee.RequestOrderStatus)),
		"logistics_channel_id":         string(optsShopee.LogisticsChanelID),
	}

	// Set Params Reuired
	if string(optsShopee.TimeRange) == "" {
		paramsQuery["time_range_field"] = "create_time"
	} else {
		paramsQuery["time_range_field"] = string(optsShopee.TimeRange)
	}

	if string(strconv.FormatInt(optsShopee.TimeFrom, 10)) == "" {
		paramsQuery["time_from"] = time.Now().Truncate(24 * time.Hour).Format(time.RFC3339)
	} else {
		paramsQuery["time_from"] = string(strconv.FormatInt(int64(optsShopee.TimeFrom), 10))
	}

	if string(strconv.FormatInt(int64(optsShopee.TimeTo), 10)) == "" {
		paramsQuery["time_to"] = time.Now().Truncate(24 * time.Hour).Add(24*time.Hour + 59*time.Minute + 59*time.Second).Format(time.RFC3339)
	} else {
		paramsQuery["time_to"] = string(strconv.FormatInt(int64(optsShopee.TimeTo), 10))
	}

	if string(strconv.FormatInt(int64(optsShopee.PageSize), 10)) == "" {
		paramsQuery["page_size"] = "20"
	} else {
		paramsQuery["page_size"] = string(strconv.FormatInt(int64(optsShopee.PageSize), 10))
	}
	// End Set Params Reuired

	for k, v := range paramsQuery {
		if v != "" {
			q.Set(k, v)
		}
	}

	queryString := q.Encode()
	url := fmt.Sprintf("%s%s/order/get_order_list?%s", s.BaseURL, s.PrefixURL, queryString)
	// s.logger.Debug("GetOrderListByShopID", zap.String("url", url))

	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	// s.logger.Debug("GetOrderListByShopID", zap.String("url", url))

	resp, err := s.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	s.Logger.Debug("adapter.GetOrderListByShopID.bodyBytes", zap.Any("bodyBytes", bodyBytes))

	var resGetShopByPartnerPublic dto.IResGetOrderListByShopIDShop
	if err := json.Unmarshal(bodyBytes, &resGetShopByPartnerPublic); err != nil {
		s.Logger.Error("adapter.GetAccessToken : json.Unmarshal(bodyBytes ,&authRes) :", zap.Error(err))
		return nil, err
	}
	if resGetShopByPartnerPublic.Error != "" {
		s.Logger.Error("adapter.GetAccessToken : json.Unmarshal(bodyBytes ,&authRes) :",
			zap.String("error", resGetShopByPartnerPublic.Error),
			zap.String("message", resGetShopByPartnerPublic.Message))
		return nil, errors.New(resGetShopByPartnerPublic.Error)
	}
	// if resGetShopByPartnerPublic.OrderList == nil {
	//   resGetShopByPartnerPublic.OrderList = []dto.IResOrderList{}
	// }

	s.Logger.Debug("resGetShopByPartnerPublic", zap.Any("resGetShopByPartnerPublic", resGetShopByPartnerPublic))

	return &resGetShopByPartnerPublic, nil
}

func (s *shopeeApi) GetOrderDetailByOrderSN(partnerID string, partnerKey string,accessToken string, shopID string, orderList []string, pending bool, option bool) (*dto.IResOrderDetailByOrderSN, error) {

  genData,err := s.GenerateSignWithPathURL("SHOP", "/order/get_order_detail", partnerID, partnerKey, shopID, "", accessToken)
  if err != nil {
    s.Logger.Error("adapter.GetOrderDetailByOrderSN : s.GenerateSignWithPathURL error", zap.Error(err))
    return nil, err
  }

	timeStp := strconv.FormatInt(time.Now().Unix(), 10)

	q := url.Values{}
	paramsQuery := map[string]string{
		// common
		"partner_id":   partnerID,
		"timestamp":    timeStp,
		"access_token": accessToken,
		"shop_id":      shopID,
		"sign":         genData.Sign,
		// required
		"order_sn_list": "",
		// opts
		"request_order_status_pending": "",
		"response_optional_fields":     "",
	}

	if len(orderList) > 0 {
		paramsQuery["order_sn_list"] = strings.Join(orderList, ",")
	}
	if pending {
		paramsQuery["request_order_status_pending"] = "true"
	}
	if option {
		paramsQuery["response_optional_fields"] = "total_amount"
	}
	for k, v := range paramsQuery {
		if v != "" {
			q.Set(k, v)
		}
	}
	queryString := q.Encode()

	url := fmt.Sprintf("%s%s/order/get_order_detail?%s", s.BaseURL, s.PrefixURL, queryString)

	s.Logger.Debug("GetOrderListByShopID", zap.String("url", url))

	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	s.Logger.Debug("GetOrderListByShopID", zap.String("url", url))

	resp, err := s.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	s.Logger.Debug("adapter.GetOrderListByShopID.bodyBytes", zap.Any("bodyBytes", bodyBytes))

	var resGetOrderDetailByOrderSN dto.IResOrderDetailByOrderSN
	if err := json.Unmarshal(bodyBytes, &resGetOrderDetailByOrderSN); err != nil {
		s.Logger.Error("adapter.GetAccessToken : json.Unmarshal(bodyBytes ,&authRes) :", zap.Error(err))
		return nil, err
	}
	if resGetOrderDetailByOrderSN.Error != "" {
		s.Logger.Error("adapter.GetAccessToken : json.Unmarshal(bodyBytes ,&resGetOrderDetailByOrderSN) :",
			zap.String("error", resGetOrderDetailByOrderSN.Error),
			zap.String("message", resGetOrderDetailByOrderSN.Message))
		return nil, errors.New(resGetOrderDetailByOrderSN.Error)
	}

	s.Logger.Debug("resGetShopByPartnerPublic", zap.Any("resGetShopByPartnerPublic", resGetOrderDetailByOrderSN))

	return &resGetOrderDetailByOrderSN, nil
}

// ------------------------------------ Demo template -----------------------------------
// Method:Service:Type
// func (s *shopeeAapi) GetShopByPartnerPublic() (string, error) { return "", nil }
// ------------------------------------ End Demo template -------------------------------
