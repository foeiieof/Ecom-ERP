package adapter

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type ShopeeAPI struct {
	baseURL    string
	partnerID  string
	partnerKey string
	httpClient *http.Client
}

type ShopeeAuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpireIn     int    `json:"expire_in"`
	ShopID       int64  `json:"shop_id"`
	Error        string `json:"error"`
	Message      string `json:"message"`
}

// NewShopeeAPI initializes adapter
func NewShopeeAPI(baseURL, partnerID, partnerKey string) *ShopeeAPI {
	return &ShopeeAPI{
		baseURL:    baseURL,
		partnerID:  partnerID,
		partnerKey: partnerKey,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (a *ShopeeAPI) ExchangeToken(ctx context.Context, code string, redirectURI string) (*ShopeeAuthResponse, error) {
	url := fmt.Sprintf("%s/api/v2/auth/token", a.baseURL)

	payload := map[string]string{
		"code":         code,
		"partner_id":   a.partnerID,
		"redirect_uri": redirectURI,
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var authResp ShopeeAuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return nil, err
	}

	if authResp.Error != "" {
		return nil, errors.New(authResp.Message)
	}

	return &authResp, nil
}
