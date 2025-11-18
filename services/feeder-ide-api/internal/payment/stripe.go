package payment

import (
	"fmt"
	"os"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/customer"
	"github.com/stripe/stripe-go/v76/subscription"
)

// StripeService Stripe 付費服務
type StripeService struct {
	apiKey string
}

// NewStripeService 建立新的 Stripe 服務
func NewStripeService() *StripeService {
	apiKey := os.Getenv("STRIPE_SECRET_KEY")
	if apiKey == "" {
		return nil // Stripe 未配置
	}

	stripe.Key = apiKey
	return &StripeService{
		apiKey: apiKey,
	}
}

// CreateCheckoutSession 創建付費 session
func (s *StripeService) CreateCheckoutSession(userID, userEmail, tier string, successURL, cancelURL string) (*stripe.CheckoutSession, error) {
	if s == nil {
		return nil, fmt.Errorf("Stripe not configured")
	}

	// 定義價格（單位：分）
	var priceID string
	switch tier {
	case "premium":
		// 從環境變數取得價格 ID，或使用默認值
		priceID = os.Getenv("STRIPE_PREMIUM_PRICE_ID")
		if priceID == "" {
			return nil, fmt.Errorf("STRIPE_PREMIUM_PRICE_ID not configured")
		}
	default:
		return nil, fmt.Errorf("unsupported tier: %s", tier)
	}

	// 創建或取得客戶
	customerParams := &stripe.CustomerParams{
		Email: stripe.String(userEmail),
		Metadata: map[string]string{
			"user_id": userID,
		},
	}
	customer, err := customer.New(customerParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	// 創建 checkout session
	params := &stripe.CheckoutSessionParams{
		Customer: stripe.String(customer.ID),
		Mode:      stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String(successURL),
		CancelURL:  stripe.String(cancelURL),
		Metadata: map[string]string{
			"user_id": userID,
			"tier":    tier,
		},
	}

	sess, err := session.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create checkout session: %w", err)
	}

	return sess, nil
}

// GetSubscription 取得訂閱資訊
func (s *StripeService) GetSubscription(subscriptionID string) (*stripe.Subscription, error) {
	if s == nil {
		return nil, fmt.Errorf("Stripe not configured")
	}

	sub, err := subscription.Get(subscriptionID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	return sub, nil
}

// CancelSubscription 取消訂閱
func (s *StripeService) CancelSubscription(subscriptionID string) (*stripe.Subscription, error) {
	if s == nil {
		return nil, fmt.Errorf("Stripe not configured")
	}

	params := &stripe.SubscriptionCancelParams{}
	sub, err := subscription.Cancel(subscriptionID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel subscription: %w", err)
	}

	return sub, nil
}

// VerifyWebhookSignature 驗證 webhook 簽名（已在 webhook.go 中實現）
// 此方法保留用於向後兼容，實際驗證在 webhook handler 中進行

