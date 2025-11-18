package payment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// PayPalService PayPal 付費服務
type PayPalService struct {
	clientID     string
	clientSecret string
	baseURL      string
	accessToken  string
}

// NewPayPalService 建立新的 PayPal 服務
func NewPayPalService() *PayPalService {
	clientID := os.Getenv("PAYPAL_CLIENT_ID")
	clientSecret := os.Getenv("PAYPAL_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		return nil // PayPal 未配置
	}

	// 判斷是沙盒還是生產環境
	baseURL := os.Getenv("PAYPAL_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.sandbox.paypal.com" // 默認使用沙盒
	}

	service := &PayPalService{
		clientID:     clientID,
		clientSecret: clientSecret,
		baseURL:      baseURL,
	}

	// 獲取 access token
	if err := service.refreshAccessToken(); err != nil {
		return nil
	}

	return service
}

// refreshAccessToken 刷新 access token
func (s *PayPalService) refreshAccessToken() error {
	url := s.baseURL + "/v1/oauth2/token"

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(s.clientID, s.clientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	body := bytes.NewBufferString("grant_type=client_credentials")
	req.Body = io.NopCloser(body)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return err
	}

	s.accessToken = tokenResp.AccessToken
	return nil
}

// CreateSubscription 創建訂閱
func (s *PayPalService) CreateSubscription(userID, tier string, returnURL, cancelURL string) (*PayPalSubscription, error) {
	if s == nil {
		return nil, fmt.Errorf("PayPal not configured")
	}

	// 定義計劃 ID（從環境變數取得）
	var planID string
	switch tier {
	case "premium":
		planID = os.Getenv("PAYPAL_PREMIUM_PLAN_ID")
		if planID == "" {
			return nil, fmt.Errorf("PAYPAL_PREMIUM_PLAN_ID not configured")
		}
	default:
		return nil, fmt.Errorf("unsupported tier: %s", tier)
	}

	url := s.baseURL + "/v1/billing/subscriptions"

	payload := map[string]interface{}{
		"plan_id": planID,
		"subscriber": map[string]interface{}{
			"email_address": "", // 需要從用戶資料取得
		},
		"application_context": map[string]interface{}{
			"return_url": returnURL,
			"cancel_url": cancelURL,
		},
		"custom_id": userID,
	}

	jsonData, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+s.accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var subscription PayPalSubscription
	if err := json.NewDecoder(resp.Body).Decode(&subscription); err != nil {
		return nil, err
	}

	return &subscription, nil
}

// GetSubscription 取得訂閱資訊
func (s *PayPalService) GetSubscription(subscriptionID string) (*PayPalSubscription, error) {
	if s == nil {
		return nil, fmt.Errorf("PayPal not configured")
	}

	url := s.baseURL + "/v1/billing/subscriptions/" + subscriptionID

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+s.accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var subscription PayPalSubscription
	if err := json.NewDecoder(resp.Body).Decode(&subscription); err != nil {
		return nil, err
	}

	return &subscription, nil
}

// CancelSubscription 取消訂閱
func (s *PayPalService) CancelSubscription(subscriptionID string, reason string) error {
	if s == nil {
		return fmt.Errorf("PayPal not configured")
	}

	url := s.baseURL + "/v1/billing/subscriptions/" + subscriptionID + "/cancel"

	payload := map[string]interface{}{
		"reason": reason,
	}

	jsonData, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+s.accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to cancel subscription: %d", resp.StatusCode)
	}

	return nil
}

// PayPalSubscription PayPal 訂閱結構
type PayPalSubscription struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	PlanID string `json:"plan_id"`
	// 其他欄位根據 PayPal API 文檔添加
}

