package payment

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/feeder-platform/feeder-ide-api/internal/user"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/webhook"
)

// WebhookHandler Webhook 處理器
type WebhookHandler struct {
	stripeService *StripeService
	paypalService *PayPalService
	userRepo      user.Repository
	userService   *user.Service
}

// NewWebhookHandler 建立新的 webhook 處理器
func NewWebhookHandler(stripeService *StripeService, paypalService *PayPalService, userRepo user.Repository, userService *user.Service) *WebhookHandler {
	return &WebhookHandler{
		stripeService: stripeService,
		paypalService: paypalService,
		userRepo:      userRepo,
		userService:   userService,
	}
}

// HandleStripeWebhook 處理 Stripe webhook
func (h *WebhookHandler) HandleStripeWebhook(payload []byte, signature string) error {
	if h.stripeService == nil {
		return fmt.Errorf("Stripe not configured")
	}

	webhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	if webhookSecret == "" {
		return fmt.Errorf("STRIPE_WEBHOOK_SECRET not configured")
	}

	event, err := webhook.ConstructEvent(payload, signature, webhookSecret)
	if err != nil {
		return fmt.Errorf("failed to verify webhook signature: %w", err)
	}

	switch event.Type {
	case "checkout.session.completed":
		return h.handleStripeCheckoutCompleted(event)
	case "customer.subscription.updated":
		return h.handleStripeSubscriptionUpdated(event)
	case "customer.subscription.deleted":
		return h.handleStripeSubscriptionDeleted(event)
	case "invoice.payment_succeeded":
		return h.handleStripePaymentSucceeded(event)
	default:
		// 忽略其他事件
		return nil
	}
}

// handleStripeCheckoutCompleted 處理 checkout 完成事件
func (h *WebhookHandler) handleStripeCheckoutCompleted(event stripe.Event) error {
	var session stripe.CheckoutSession
	if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
		return fmt.Errorf("failed to parse session: %w", err)
	}

	userID := session.Metadata["user_id"]
	tier := session.Metadata["tier"]

	if userID == "" {
		return fmt.Errorf("user_id not found in metadata")
	}

	// 更新用戶等級
	if err := h.userService.UpdateUserTier(userID, tier); err != nil {
		return fmt.Errorf("failed to update user tier: %w", err)
	}

	// 創建訂閱記錄
	subscription := &user.Subscription{
		UserID:               userID,
		Tier:                 tier,
		Status:               "active",
		PaymentProvider:      stringPtr("stripe"),
		PaymentSubscriptionID: stringPtr(session.Subscription.ID),
		CurrentPeriodStart:   timePtr(time.Now()),
		CurrentPeriodEnd:     timePtr(time.Now().Add(30 * 24 * time.Hour)), // 假設月付
	}

	if err := h.userRepo.CreateSubscription(subscription); err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	// 創建付費記錄
	payment := &user.Payment{
		UserID:            userID,
		SubscriptionID:    stringPtr(subscription.ID),
		Amount:            float64(session.AmountTotal) / 100, // Stripe 使用分
		Currency:          string(session.Currency),
		PaymentProvider:   "stripe",
		PaymentProviderID: session.ID,
		Status:            "completed",
	}

	if err := h.userRepo.CreatePayment(payment); err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}

	return nil
}

// handleStripeSubscriptionUpdated 處理訂閱更新事件
func (h *WebhookHandler) handleStripeSubscriptionUpdated(event stripe.Event) error {
	var sub stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &sub); err != nil {
		return fmt.Errorf("failed to parse subscription: %w", err)
	}

	// 查找訂閱記錄
	subscriptions, err := h.userRepo.GetSubscriptionsByProviderID("stripe", sub.ID)
	if err != nil || len(subscriptions) == 0 {
		return fmt.Errorf("subscription not found")
	}

	subscription := subscriptions[0]

	// 更新訂閱狀態
	if sub.Status == "active" {
		subscription.Status = "active"
	} else if sub.Status == "canceled" {
		subscription.Status = "cancelled"
	}

	if sub.CurrentPeriodStart > 0 {
		subscription.CurrentPeriodStart = timePtr(time.Unix(sub.CurrentPeriodStart, 0))
	}
	if sub.CurrentPeriodEnd > 0 {
		subscription.CurrentPeriodEnd = timePtr(time.Unix(sub.CurrentPeriodEnd, 0))
	}

	return h.userRepo.UpdateSubscription(subscription)
}

// handleStripeSubscriptionDeleted 處理訂閱刪除事件
func (h *WebhookHandler) handleStripeSubscriptionDeleted(event stripe.Event) error {
	var sub stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &sub); err != nil {
		return fmt.Errorf("failed to parse subscription: %w", err)
	}

	// 查找訂閱記錄
	subscriptions, err := h.userRepo.GetSubscriptionsByProviderID("stripe", sub.ID)
	if err != nil || len(subscriptions) == 0 {
		// 如果找不到，跳過（可能是已經取消的訂閱）
		return nil
	}

	subscription := subscriptions[0]
	subscription.Status = "cancelled"

	return h.userRepo.UpdateSubscription(subscription)
}

// handleStripePaymentSucceeded 處理付費成功事件
func (h *WebhookHandler) handleStripePaymentSucceeded(event stripe.Event) error {
	var invoice stripe.Invoice
	if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
		return fmt.Errorf("failed to parse invoice: %w", err)
	}

	// 創建付費記錄（需要從 subscription 取得 user_id）
	if invoice.Subscription != nil {
		subscriptions, err := h.userRepo.GetSubscriptionsByProviderID("stripe", invoice.Subscription.ID)
		if err == nil && len(subscriptions) > 0 {
			payment := &user.Payment{
				UserID:            subscriptions[0].UserID,
				SubscriptionID:    stringPtr(subscriptions[0].ID),
				Amount:            float64(invoice.AmountPaid) / 100,
				Currency:          string(invoice.Currency),
				PaymentProvider:   "stripe",
				PaymentProviderID: invoice.ID,
				Status:            "completed",
			}
			return h.userRepo.CreatePayment(payment)
		}
	}
	return nil
}

// HandlePayPalWebhook 處理 PayPal webhook
func (h *WebhookHandler) HandlePayPalWebhook(payload []byte) error {
	if h.paypalService == nil {
		return fmt.Errorf("PayPal not configured")
	}

	var event map[string]interface{}
	if err := json.Unmarshal(payload, &event); err != nil {
		return fmt.Errorf("failed to parse webhook: %w", err)
	}

	eventType, ok := event["event_type"].(string)
	if !ok {
		return fmt.Errorf("invalid event type")
	}

	switch eventType {
	case "BILLING.SUBSCRIPTION.CREATED":
		return h.handlePayPalSubscriptionCreated(event)
	case "BILLING.SUBSCRIPTION.UPDATED":
		return h.handlePayPalSubscriptionUpdated(event)
	case "BILLING.SUBSCRIPTION.CANCELLED":
		return h.handlePayPalSubscriptionCancelled(event)
	case "PAYMENT.SALE.COMPLETED":
		return h.handlePayPalPaymentCompleted(event)
	default:
		return nil
	}
}

// handlePayPalSubscriptionCreated 處理 PayPal 訂閱創建事件
func (h *WebhookHandler) handlePayPalSubscriptionCreated(event map[string]interface{}) error {
	resource, ok := event["resource"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid resource")
	}

	subscriptionID, _ := resource["id"].(string)
	customID, _ := resource["custom_id"].(string)

	if customID == "" {
		return fmt.Errorf("user_id not found")
	}

	// 更新用戶等級
	if err := h.userService.UpdateUserTier(customID, "premium"); err != nil {
		return fmt.Errorf("failed to update user tier: %w", err)
	}

	// 創建訂閱記錄
	subscription := &user.Subscription{
		UserID:               customID,
		Tier:                 "premium",
		Status:               "active",
		PaymentProvider:      stringPtr("paypal"),
		PaymentSubscriptionID: stringPtr(subscriptionID),
		CurrentPeriodStart:   timePtr(time.Now()),
		CurrentPeriodEnd:     timePtr(time.Now().Add(30 * 24 * time.Hour)),
	}

	return h.userRepo.CreateSubscription(subscription)
}

// handlePayPalSubscriptionUpdated 處理 PayPal 訂閱更新事件
func (h *WebhookHandler) handlePayPalSubscriptionUpdated(event map[string]interface{}) error {
	// 類似 Stripe 的處理邏輯
	return nil
}

// handlePayPalSubscriptionCancelled 處理 PayPal 訂閱取消事件
func (h *WebhookHandler) handlePayPalSubscriptionCancelled(event map[string]interface{}) error {
	// 類似 Stripe 的處理邏輯
	return nil
}

// handlePayPalPaymentCompleted 處理 PayPal 付費完成事件
func (h *WebhookHandler) handlePayPalPaymentCompleted(event map[string]interface{}) error {
	// 類似 Stripe 的處理邏輯
	return nil
}

// 輔助函數
func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}

