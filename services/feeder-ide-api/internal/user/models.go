package user

import "time"

// User 用戶模型
type User struct {
	ID                  string     `json:"id"`
	Email               string     `json:"email"`
	Name                *string    `json:"name,omitempty"`
	AvatarURL           *string    `json:"avatar_url,omitempty"`
	SubscriptionTier    string     `json:"subscription_tier"`    // demo, free, premium
	SubscriptionStatus  string     `json:"subscription_status"`  // active, cancelled, expired
	SubscriptionExpiresAt *time.Time `json:"subscription_expires_at,omitempty"`
	APIKey              *string    `json:"api_key,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

// Subscription 訂閱模型
type Subscription struct {
	ID                   string     `json:"id"`
	UserID               string     `json:"user_id"`
	Tier                 string     `json:"tier"` // free, premium
	Status               string     `json:"status"` // active, cancelled, expired, pending
	PaymentProvider      *string    `json:"payment_provider,omitempty"` // stripe, paypal
	PaymentSubscriptionID *string   `json:"payment_subscription_id,omitempty"`
	CurrentPeriodStart   *time.Time `json:"current_period_start,omitempty"`
	CurrentPeriodEnd     *time.Time `json:"current_period_end,omitempty"`
	CancelAtPeriodEnd    bool       `json:"cancel_at_period_end"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

// UserOAuth OAuth 關聯模型
type UserOAuth struct {
	ID             string     `json:"id"`
	UserID         string     `json:"user_id"`
	Provider       string     `json:"provider"` // google, github
	ProviderUserID string     `json:"provider_user_id"`
	AccessToken    *string    `json:"access_token,omitempty"`
	RefreshToken   *string    `json:"refresh_token,omitempty"`
	TokenExpiresAt *time.Time `json:"token_expires_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// UserQuota 用戶配額模型
type UserQuota struct {
	ID                      string    `json:"id"`
	UserID                  string    `json:"user_id"`
	MaxTopologies           int       `json:"max_topologies"`
	UsedTopologies          int       `json:"used_topologies"`
	MaxSimulationsPerDay    int       `json:"max_simulations_per_day"`
	UsedSimulationsToday    int       `json:"used_simulations_today"`
	LastSimulationResetDate time.Time `json:"last_simulation_reset_date"`
	CanUse3DRendering       bool      `json:"can_use_3d_rendering"`
	CanUseAIPrediction      bool      `json:"can_use_ai_prediction"`
	CanUseAdvancedSecurity  bool      `json:"can_use_advanced_security"`
	CanAccessAPI            bool      `json:"can_access_api"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

// Payment 付費記錄模型
type Payment struct {
	ID                string                 `json:"id"`
	UserID            string                 `json:"user_id"`
	SubscriptionID    *string                `json:"subscription_id,omitempty"`
	Amount            float64                `json:"amount"`
	Currency          string                 `json:"currency"`
	PaymentProvider   string                 `json:"payment_provider"` // stripe, paypal, usdt
	PaymentProviderID string                 `json:"payment_provider_id"`
	Status            string                 `json:"status"` // pending, completed, failed, refunded
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
}

