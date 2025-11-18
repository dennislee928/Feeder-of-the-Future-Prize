package api

import (
	"net/http"
	"os"

	"github.com/feeder-platform/feeder-ide-api/internal/auth"
	"github.com/feeder-platform/feeder-ide-api/internal/payment"
	"github.com/feeder-platform/feeder-ide-api/internal/user"
	"github.com/gin-gonic/gin"
)

// PaymentHandler 付費處理器
type PaymentHandler struct {
	stripeService *payment.StripeService
	paypalService *payment.PayPalService
	webhookHandler *payment.WebhookHandler
	userRepo      user.Repository
	userService   *user.Service
}

// NewPaymentHandler 建立新的付費處理器
func NewPaymentHandler(stripeService *payment.StripeService, paypalService *payment.PayPalService, webhookHandler *payment.WebhookHandler, userRepo user.Repository, userService *user.Service) *PaymentHandler {
	return &PaymentHandler{
		stripeService:  stripeService,
		paypalService:  paypalService,
		webhookHandler: webhookHandler,
		userRepo:       userRepo,
		userService:    userService,
	}
}

// CreateCheckoutRequest 創建付費 session 請求
type CreateCheckoutRequest struct {
	Tier      string `json:"tier" binding:"required,oneof=premium"` // 目前只支持 premium
	Provider  string `json:"provider" binding:"required,oneof=stripe paypal"`
}

// CreateCheckout 創建付費 session
func (h *PaymentHandler) CreateCheckout(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	var req CreateCheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 取得用戶資訊
	user, err := h.userRepo.GetUserByID(*userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 構建成功和取消 URL
	baseURL := os.Getenv("FRONTEND_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:3001"
	}
	successURL := baseURL + "/payment/success"
	cancelURL := baseURL + "/payment/cancel"

	switch req.Provider {
	case "stripe":
		if h.stripeService == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Stripe not configured"})
			return
		}

		session, err := h.stripeService.CreateCheckoutSession(*userID, user.Email, req.Tier, successURL, cancelURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create checkout session: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"checkout_url": session.URL,
			"session_id":   session.ID,
		})

	case "paypal":
		if h.paypalService == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "PayPal not configured"})
			return
		}

		subscription, err := h.paypalService.CreateSubscription(*userID, req.Tier, successURL, cancelURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subscription: " + err.Error()})
			return
		}

		// PayPal 返回的訂閱包含 approval URL
		c.JSON(http.StatusOK, gin.H{
			"subscription_id": subscription.ID,
			"status":          subscription.Status,
		})

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported provider"})
	}
}

// HandleStripeWebhook 處理 Stripe webhook
func (h *PaymentHandler) HandleStripeWebhook(c *gin.Context) {
	payload, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read payload"})
		return
	}

	signature := c.GetHeader("Stripe-Signature")
	if signature == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing signature"})
		return
	}

	if err := h.webhookHandler.HandleStripeWebhook(payload, signature); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// HandlePayPalWebhook 處理 PayPal webhook
func (h *PaymentHandler) HandlePayPalWebhook(c *gin.Context) {
	payload, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read payload"})
		return
	}

	if err := h.webhookHandler.HandlePayPalWebhook(payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// GetPaymentHistory 取得付費歷史
func (h *PaymentHandler) GetPaymentHistory(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	payments, err := h.userRepo.GetPaymentsByUserID(*userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get payments: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"payments": payments})
}

