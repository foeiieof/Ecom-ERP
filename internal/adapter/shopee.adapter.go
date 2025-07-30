package adapter

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"go.uber.org/zap"

	"ecommerce/internal/adapter/dto"
)

type IShopeeService interface {

	// waiting reface replace body abd query gen
	// GenerateBodyQueryParams()(,error)

	GetAccessToken(partnerID string, shopID string, code string, signCode string) (*IResShopeeAuthResponse, error)
	GetRefreshToken(partnerID string, shopID string, refreshToken string, signCode string) (*dto.IResShopeeAuthRefreshResponse, error)
	// ExchangeToken(ctx context.Context, code string, redirectURI string, partnerID string) (*ShopeeAuthResponse, error)
	GetShopByPartnerPublic(partnerID string, signCode string) (*dto.IResGetShopByPartnerPublic, error)
	GetOrderListByShopID(partnerID string, accessToken string, shopID string, signCode string, optsShopee *dto.IOptionShopeeQuery) (*dto.IResGetOrderListByShopIDShop, error)
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

func NewShopeeAPI(baseURL string, prefix string, log *zap.Logger) IShopeeService {
	return &shopeeApi{
		baseURL:    baseURL,
		prefixURL:  prefix,
		logger:     log,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
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

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}


	s.logger.Debug("adapter.GetAccessToken: Outgoing Request After",
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
		s.logger.Error("adapter.GetAccessToken : json.Unmarshal(bodyBytes ,&authResp) :", zap.Error(err))
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
	url := fmt.Sprintf("%s%s/auth/access_token/get?partner_id=%s&timestamp=%s&sign=%s", s.baseURL, s.prefixURL, partnerID, timeStp, signCode)

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

	s.logger.Debug("adapter.GetAccessToken: Outgoing Request Before",
		zap.String("url", req.URL.String()),
		zap.String("method", req.Method),
		zap.Any("headers", req.Header),
		zap.ByteString("body", body),
	)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	s.logger.Debug("adapter.GetRefreshToken.bodyBytes", zap.Any("bodyBytes", bodyBytes))

	var resShopeeAuthRefreshToken dto.IResShopeeAuthRefreshResponse
	if err := json.Unmarshal(bodyBytes, &resShopeeAuthRefreshToken); err != nil {
		s.logger.Error("adapter.GetAccessToken : json.Unmarshal(bodyBytes ,&authRes) :", zap.Error(err))
		return nil, err
	}

	return &resShopeeAuthRefreshToken, nil
}

func (s *shopeeApi) GetShopByPartnerPublic(partnerID string, signCode string) (*dto.IResGetShopByPartnerPublic, error) {
	// https://partner.test-stable.shopeemobile.com/api/v2/public/get_shops_by_partner
	// query: partner_id, timestamp, sign

	// Public Api
	timeStp := strconv.FormatInt(time.Now().Unix(), 10)

	url := fmt.Sprintf("%s%s/public/get_shops_by_partner?partner_id=%s&timestamp=%s&sign=%s", s.baseURL, s.prefixURL, partnerID, timeStp, signCode)

	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	s.logger.Debug("bodyBytes", zap.Any("bodyBytes", bodyBytes))

	var resGetShopByPartnerPublic dto.IResGetShopByPartnerPublic
	if err := json.Unmarshal(bodyBytes, &resGetShopByPartnerPublic); err != nil {
		s.logger.Error("adapter.GetAccessToken : json.Unmarshal(bodyBytes ,&authRes) :", zap.Error(err))
		return nil, err
	}

  if resGetShopByPartnerPublic.Error != "" {
    return nil, errors.New(resGetShopByPartnerPublic.Message)
  }

  s.logger.Debug("resGetShopByPartnerPublic", zap.Any("resGetShopByPartnerPublic", resGetShopByPartnerPublic))

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
	} else { paramsQuery["time_range_field"] = string(optsShopee.TimeRange) }

	if string(strconv.FormatInt(optsShopee.TimeFrom, 10)) == "" {
		paramsQuery["time_from"] = time.Now().Truncate(24 * time.Hour).Format(time.RFC3339)
	} else { paramsQuery["time_from"] = string(strconv.FormatInt(int64(optsShopee.TimeFrom), 10)) }

	if string(strconv.FormatInt(int64(optsShopee.TimeTo), 10)) == "" {
		paramsQuery["time_to"] = time.Now().Truncate(24 * time.Hour).Add(24*time.Hour + 59*time.Minute + 59*time.Second).Format(time.RFC3339)
	} else { paramsQuery["time_to"] = string(strconv.FormatInt(int64(optsShopee.TimeTo), 10)) }

	if string(strconv.FormatInt(int64(optsShopee.PageSize), 10)) == "" {
		paramsQuery["page_size"] = "20"
	} else { paramsQuery["page_size"] = string(strconv.FormatInt(int64(optsShopee.PageSize), 10)) }
  // End Set Params Reuired

	for k, v := range paramsQuery {
		if v != "" { q.Set(k, v)
		}
	}

	queryString := q.Encode()
	url := fmt.Sprintf("%s%s/order/get_order_list?%s", s.baseURL, s.prefixURL, queryString)
	// s.logger.Debug("GetOrderListByShopID", zap.String("url", url))

	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

  // s.logger.Debug("GetOrderListByShopID", zap.String("url", url))

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err }
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err }

	s.logger.Debug("adapter.GetOrderListByShopID.bodyBytes", zap.Any("bodyBytes", bodyBytes))

	var resGetShopByPartnerPublic dto.IResGetOrderListByShopIDShop
	if err := json.Unmarshal(bodyBytes, &resGetShopByPartnerPublic); err != nil {
		s.logger.Error("adapter.GetAccessToken : json.Unmarshal(bodyBytes ,&authRes) :", zap.Error(err))
		return nil, err
	}
  if resGetShopByPartnerPublic.Error != "" {
    s.logger.Error("adapter.GetAccessToken : json.Unmarshal(bodyBytes ,&authRes) :", 
      zap.String("error", resGetShopByPartnerPublic.Error), 
      zap.String("message", resGetShopByPartnerPublic.Message))
    return nil, errors.New(resGetShopByPartnerPublic.Error)
  }
  // if resGetShopByPartnerPublic.OrderList == nil {
  //   resGetShopByPartnerPublic.OrderList = []dto.IResOrderList{} 
  // }

  s.logger.Debug("resGetShopByPartnerPublic", zap.Any("resGetShopByPartnerPublic", resGetShopByPartnerPublic))

	return &resGetShopByPartnerPublic, nil
}

// ------------------------------------ Demo template -----------------------------------
// Method:Service:Type
// func (s *shopeeAapi) GetShopByPartnerPublic() (string, error) { return "", nil }
// ------------------------------------ End Demo template -------------------------------
