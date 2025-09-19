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

type ShopeeTypeAPIEnum string
const (
  SHOP ShopeeTypeAPIEnum = "SHOP"
  PUBLIC ShopeeTypeAPIEnum = "PUBLIC"
  MERCHANT ShopeeTypeAPIEnum = "MERCHANT"
)

type ShopeePathAPIEnum string 
// ::TABLE_METHOD
const (
  SHOP_GET_PROFILE_API ShopeePathAPIEnum   = "/api/v2/shop/get_profile"
  SHOP_GET_SHOP_INFO_API ShopeePathAPIEnum = "/api/v2/shop/get_shop_info"
)

type IReqShopeeAdapter struct {
  PartnerID string 
  // TimeStamp int  // time.Unix
  AccessToken string 
  ShopID    string 
  SecretKey string
  Code *string
  OrderSN []string
}

type IShopeeService interface {

	// waiting reface replace body abd query gen
	// GenerateBodyQueryParams()(,error)

  RequestHTTP(method string, url string , body *[]byte) (*http.Response, error)

	GenerateSignWithPathURL(state string, pathUrl string, partnerID string, partnerKey string, shopID string, code string, accessToken string) (*IResGenerateSignWithUri, error)

	GetAccessToken(partnerID string, shopID string, code string, signCode string) (*IResShopeeAuthResponse, error)
	GetRefreshToken(partnerID string, shopID string, refreshToken string, signCode string) (*dto.IResShopeeAuthRefreshResponse, error)
	// ExchangeToken(ctx context.Context, code string, redirectURI string, partnerID string) (*ShopeeAuthResponse, error)
	GetShopByPartnerPublic(partnerID string, signCode string) (*dto.IResGetShopByPartnerPublic, error)

	GetOrderListByShopID(partnerID string, accessToken string, shopID string, signCode string, optsShopee *dto.IOptionShopeeQuery) (*dto.IResGetOrderListByShopIDShop, error)
  GetOrderDetailByOrderSN(partnerID string, partnerKey string,accessToken string, shopID string, orderList []string, pending bool, option bool) (*dto.IResOrderDetailByOrderSN, error)
  
  // path : */api/v2/order/get_order_detail
  GetOrderDetailListByOrderSN(ctx context.Context, params *IReqShopeeAdapter) ([]dto.IResOrderListWithDetails,error)


  // path : */api/v2/shop/get_profile 
  GetShopProfile(ctx context.Context, params *IReqShopeeAdapter ) (*dto.IResShopGetProfile_ResponseDTO, error)  
  // path : */api/v2/shop/get_shop_info
  GetShopInfo(ctx context.Context, params *IReqShopeeAdapter) (*dto.IResShopGetShopInfoDTO ,error)
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
	TimeStamp string
  URL       *url.URL
}

func (s *shopeeApi) RequestHTTP(method string, url string , body *[]byte) (*http.Response, error) {
  switch method {
  case "GET":
    // s.Logger.Debug("adapter.RequestHTTP.GET", zap.String("url" , url))
    return http.Get(url)
  case "POST":
    req,err := http.NewRequest("POST", url, bytes.NewBuffer(*body))
    if err != nil {return nil, err}
    req.Header.Set("Content-Type","application/json")
    return http.DefaultClient.Do(req)
  default: 
    return nil, errors.New("adapter.RequestHTTP: upsupport method" )
  }
}

// Tip : func auto complete fill  /api/v2/***(shopee)
func (s *shopeeApi) GenerateSignWithPathURL(state string, pathUrl string, partnerID string, partnerKey string, shopID string, code string, accessToken string) (*IResGenerateSignWithUri, error) {
  // var url string
	var method string
  Url,err := url.Parse(s.Config.Shopee.ShopeeApiBaseUrl) 
  if err != nil { return nil , errors.New("error parse url from env")}

	// host := s.Config.Shopee.ShopeeApiBaseUrl
	timest := strconv.FormatInt(time.Now().Unix(), 10)
	path := fmt.Sprintf("%s", pathUrl)
	// s.Logger.Sugar().Debugf("adapter.shopee.GenerateSignWithPathURL: %s", path)
	var baseString string
  var sign string
  q := Url.Query()
  q.Set("partner_id", partnerID)
  q.Set("timestamp", timest)
	// baseString := fmt.Sprintf("%s%s%s", partnerID, path, timest)
	switch state {
	case "PUBLIC":
		// For Public APIs: partner_id, api path, timestamp
		baseString = fmt.Sprintf("%s%s%s", partnerID, path, timest)
	  h := hmac.New(sha256.New, []byte(partnerKey))
	  h.Write([]byte(baseString))
	  sign = hex.EncodeToString(h.Sum(nil))
    // q.Set("sign", sign)

	case "SHOP":
		// For Shop APIs: partner_id, api path, timestamp, access_token, shop_id
    // s.Logger.Debug("adapter.GenerateSignWithPathURL.SHOP", zap.String("val", partnerID+path+timest+accessToken+shopID))
		baseString = fmt.Sprintf("%s%s%s%s%s", partnerID, path, timest, accessToken, shopID)
	  h := hmac.New(sha256.New, []byte(partnerKey))
	  h.Write([]byte(baseString))
	  sign = hex.EncodeToString(h.Sum(nil))

    // q.Set("sign", sign)
    q.Set("shop_id", shopID)
    q.Set("access_token", accessToken)

	case "MERCHANT":
		// Not available
		// For Merchant APIs: partner_id, api path, timestamp, access_token, merchant_id
		merchantID := ""
		baseString = fmt.Sprintf("%s%s%s%s%s", partnerID, path, timest, accessToken, merchantID)
	  h := hmac.New(sha256.New, []byte(partnerKey))
	  h.Write([]byte(baseString))
	  sign = hex.EncodeToString(h.Sum(nil))

    // q.Set("sign", sign)
    q.Set("merchant_id", merchantID)
    q.Set("access_token", accessToken)

	default:
		s.Logger.Error("adapter.shopee.GenerateSignWithPathURL: invalid state")
		return nil, errors.New("adapter.shopee.GenerateSignWithPathURL: invalid state")
	}

  q.Set("sign", sign)
  // Set query 
  Url.RawQuery = q.Encode() 
	// s.Logger.Sugar().Debugf("adapter.shopee.GenerateSignWithPathURL: baseString - %s", baseString)
  // s.Logger.Info("shopee.adapter.GenerateSignWithPathURL", zap.String("val", path))
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

  case "/api/v2/shop/get_profile":
    method = "GET"

  case "/api/v2/shop/get_shop_info":
    method = "GET"
    // url = fmt.Sprintf("%s%spartner_id=%s&timestamp=%s&sign=%s&shop_id=%s&access_token=%s", s.Config.Shopee.ShopeeApiBaseUrl,path,  ) 

	default:
		s.Logger.Error(`adapter.shopee.GenerateSignWithPathURL:invalid path `+ method)
    return nil, errors.New("adapter.shopee.GenerateSignWithPathURL: invalid path" + method + ":" + path )
	}

  Url.Path = path
  
	return &IResGenerateSignWithUri{
		Method:    method,
		Path:      path,
		Sign:      sign,
		Code:      code,
		TimeStamp: timest,
    URL: Url,
  }, nil
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
  // s.Logger.Error("shopee.adapter.GetAccessToken", zap.String("error",url))

	partnerIDInt, err := strconv.ParseInt(partnerID, 10, 64)
	if err != nil {
    // s.Logger.Error("shopee.adapter.GetAccessToken", zap.String("error", "parseIDInt failed"))
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



  // s.Logger.Info("shopee.adapter.GetAccessToken", zap.Any("val", payload) )

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
  s.Logger.Info("shopee.adapter.GetAccessToken", zap.String("val", "xx"))

	resp, err := s.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	s.Logger.Info("adapter.GetAccessToken: Outgoing Request After",
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

  if authResp.Error != "" {
    return nil , errors.New(authResp.Message)
  } 
   
  // s.Logger.Debug("shopee.adapter.GetAccessToken", zap.Any("val", authResp))

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
	url := fmt.Sprintf("%s%s/order/get_order_list?%s&response_optional_fields=order_status", s.BaseURL, s.PrefixURL, queryString)
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

	// s.Logger.Debug("adapter.GetOrderListByShopID.bodyBytes", zap.Any("bodyBytes", bodyBytes))

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

	// s.Logger.Debug("resGetShopByPartnerPublic", zap.Any("resGetShopByPartnerPublic", resGetShopByPartnerPublic))

	return &resGetShopByPartnerPublic, nil
}

func (s *shopeeApi) GetOrderDetailByOrderSN(partnerID string, partnerKey string,accessToken string, shopID string, orderList []string, pending bool, option bool) (*dto.IResOrderDetailByOrderSN, error) {

  genData,err := s.GenerateSignWithPathURL("SHOP", "/api/v2/order/get_order_detail", partnerID, partnerKey, shopID, "", accessToken)
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

	// s.Logger.Debug("GetOrderListByShopID", zap.String("url", url))

	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	// s.Logger.Debug("GetOrderListByShopID", zap.String("url", url))

	resp, err := s.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// s.Logger.Debug("adapter.GetOrderListByShopID.bodyBytes", zap.Any("bodyBytes", bodyBytes))

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

	// s.Logger.Debug("resGetShopByPartnerPublic", zap.Any("resGetShopByPartnerPublic", resGetOrderDetailByOrderSN))

	return &resGetOrderDetailByOrderSN, nil
}

func (s *shopeeApi)GetOrderDetailListByOrderSN(ctx context.Context, params *IReqShopeeAdapter) ([]dto.IResOrderListWithDetails,error) {
  state := "SHOP"
  path  := "/api/v2/order/get_order_detail"
  gen,err := s.GenerateSignWithPathURL(state, path, params.PartnerID, params.SecretKey, params.ShopID, "", params.AccessToken  )
  if err != nil {
    s.Logger.Debug("adapter.GetShopProfile", zap.Error(err))
    return nil , err }

  // cutomize url.Params

  orderSNList := strings.Join(params.OrderSN, ",")

  qUrl := gen.URL.Query()
  qUrl.Set("order_sn_list", orderSNList)
  optsQ := []string{
    "total_amount",
    "pending_terms", 
    "buyer_username",
    "recipient_address",
    "dropshipper",
    "dropshipper_phone",
    "item_list", 
    "package_list", 
    "note", 
  }
  qUrl.Set("response_optional_fields", strings.Join(optsQ, ",")  )

  gen.URL.RawQuery = qUrl.Encode()
  // s.Logger.Debug("adapter.GetOrderDetailListByOrderSN", zap.String("cUrl", gen.URL.String()))

  var resp *http.Response
  resp, err = s.RequestHTTP(gen.Method, gen.URL.String(), nil)
  if err != nil {
    s.Logger.Debug("adapter.GetShopProfile.resp", zap.Error(err))
    return nil,err}
  defer resp.Body.Close()

  bodyBytes,err := io.ReadAll(resp.Body)
  if err != nil {
    s.Logger.Debug("adapter.GetShopProfile.bodyBytes", zap.Error(err))
    return nil, err
  }

  // s.Logger.Debug("adapter.GetShopeeProfile", zap.String("bodyBytes", string(bodyBytes) ))

  var parse dto.IResOrderDetailByOrderSN
  if err := json.Unmarshal(bodyBytes, &parse) ; err != nil { 
    return nil, errors.New("invalidate parse bodyBytes in adapter") 
  }
  if parse.Error != "" {
    return nil , errors.New(parse.Error)
  }

  // s.Logger.Debug("adapter.GetShopeeProfile", zap.Any("parse", parse.Response.OrderList))

  return parse.Response.OrderList, nil
}

func (s *shopeeApi)GetShopProfile(ctx context.Context, params *IReqShopeeAdapter ) (*dto.IResShopGetProfile_ResponseDTO, error)  {
  // 0. generate sign
  state := "SHOP"
  path  := "/api/v2/shop/get_profile"
  gen,err := s.GenerateSignWithPathURL(state, path, params.PartnerID, params.SecretKey, params.ShopID, "", params.AccessToken  )
  if err != nil {
    s.Logger.Debug("adapter.GetShopProfile", zap.Error(err))
    return nil , err }

  var resp *http.Response
  resp, err = s.RequestHTTP(gen.Method, gen.URL.String(), nil)
  if err != nil {
    s.Logger.Debug("adapter.GetShopProfile.resp", zap.Error(err))
    return nil,err}
  defer resp.Body.Close()

  bodyBytes,err := io.ReadAll(resp.Body)
  if err != nil {
    s.Logger.Debug("adapter.GetShopProfile.bodyBytes", zap.Error(err))
    return nil, err}

  // s.Logger.Info("adapter.GetShopeeProfile", zap.String("bodyBytes", string(bodyBytes) ))

  var parse dto.IResShopGetProfileDTO
  if err := json.Unmarshal(bodyBytes, &parse) ; err != nil { return nil, errors.New("invalidate parse bodyBytes") }
  if parse.Error != "" {
    return nil , errors.New(parse.Error)
  }

  // s.Logger.Info("adapter.GetShopeeProfile", zap.Any("parse", parse) )

  // parse to Entity 
  return &parse.Resoponse, nil 
}

func (s *shopeeApi)GetShopInfo(ctx context.Context, params *IReqShopeeAdapter) (*dto.IResShopGetShopInfoDTO ,error) {
  // 0. generate sign
  state := "SHOP"
  path  := "/api/v2/shop/get_shop_info"
  gen,err := s.GenerateSignWithPathURL(state, path, params.PartnerID, params.SecretKey, params.ShopID, "", params.AccessToken  )
  if err != nil {
    s.Logger.Debug("adapter.GetShopInfo", zap.Error(err))
    return nil , err }

  var resp *http.Response
  resp, err = s.RequestHTTP(gen.Method, gen.URL.String(), nil)
  if err != nil {
    s.Logger.Debug("adapter.GetShopInfo.resp", zap.Error(err))
    return nil,err}
  defer resp.Body.Close()

  bodyBytes,err := io.ReadAll(resp.Body)
  if err != nil {
    s.Logger.Debug("adapter.GetShopInfo.bodyBytes", zap.Error(err))
    return nil, err}

  // s.Logger.Info("adapter.GetShopeeInfo", zap.String("bodyBytes", string(bodyBytes) ))
  var parse dto.IResShopGetShopInfoDTO
  if err := json.Unmarshal(bodyBytes, &parse) ; err != nil { return nil, errors.New("invalidate parse bodyBytes") }
  if parse.Error != "" { return nil , errors.New(parse.Error) }

  // s.Logger.Info("adapter.GetShopInfo", zap.Any("parse", parse) )
  // parse to Entity 
  return &parse, nil 
}

// ------------------------------------ Demo template -----------------------------------
// Method:Service:Type
// func (s *shopeeAapi) GetShopByPartnerPublic() (string, error) { return "", nil }
// ------------------------------------ End Demo template -------------------------------
