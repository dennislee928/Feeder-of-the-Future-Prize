package api

import (
	"net/http"
	"time"

	"github.com/feeder-platform/feeder-ide-api/internal/auth"
	"github.com/feeder-platform/feeder-ide-api/internal/user"
	"github.com/gin-gonic/gin"
)

// AuthHandler 認證處理器
type AuthHandler struct {
	oauthConfig *auth.OAuthConfig
	userRepo    user.Repository
	userService *user.Service
}

// NewAuthHandler 建立新的認證處理器
func NewAuthHandler(oauthConfig *auth.OAuthConfig, userRepo user.Repository, userService *user.Service) *AuthHandler {
	return &AuthHandler{
		oauthConfig: oauthConfig,
		userRepo:    userRepo,
		userService: userService,
	}
}

// OAuthCallbackRequest OAuth 回調請求
type OAuthCallbackRequest struct {
	Provider string `json:"provider" binding:"required"` // google, github
	Code     string `json:"code" binding:"required"`
	State    string `json:"state,omitempty"`
}

// OAuthResponse OAuth 回應
type OAuthResponse struct {
	Token     string      `json:"token"`
	User      *user.User  `json:"user"`
	ExpiresIn int         `json:"expires_in"` // 秒
}

// OAuthCallback 處理 OAuth 回調
func (h *AuthHandler) OAuthCallback(c *gin.Context) {
	var req OAuthCallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 交換授權碼取得 token
	oauthToken, err := h.oauthConfig.ExchangeCode(req.Provider, req.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange code: " + err.Error()})
		return
	}

	// 取得用戶資訊
	oauthUserInfo, err := h.oauthConfig.GetUserInfo(req.Provider, oauthToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get user info: " + err.Error()})
		return
	}

	// 檢查是否已有 OAuth 關聯
	existingOAuth, err := h.userRepo.GetOAuthByProvider(req.Provider, oauthUserInfo.ProviderUserID)
	if err != nil {
		// 沒有找到，創建新用戶
		newUser := &user.User{
			Email:              oauthUserInfo.Email,
			Name:               &oauthUserInfo.Name,
			AvatarURL:          &oauthUserInfo.AvatarURL,
			SubscriptionTier:   "free", // 註冊用戶默認為免費會員
			SubscriptionStatus: "active",
		}

		if err := h.userRepo.CreateUser(newUser); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
			return
		}

		// 創建 OAuth 關聯
		userOAuth := &user.UserOAuth{
			UserID:         newUser.ID,
			Provider:       req.Provider,
			ProviderUserID: oauthUserInfo.ProviderUserID,
		}
		if oauthToken.AccessToken != "" {
			userOAuth.AccessToken = &oauthToken.AccessToken
		}
		if oauthToken.RefreshToken != "" {
			userOAuth.RefreshToken = &oauthToken.RefreshToken
		}
		if !oauthToken.Expiry.IsZero() {
			userOAuth.TokenExpiresAt = &oauthToken.Expiry
		}

		if err := h.userRepo.CreateOrUpdateOAuth(userOAuth); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create oauth: " + err.Error()})
			return
		}

		// 創建默認配額
		quota, err := h.userService.GetUserQuota(newUser.ID)
		if err != nil || quota == nil {
			// 如果獲取失敗，創建默認配額
			quota = &user.UserQuota{
				UserID:                  newUser.ID,
				MaxTopologies:           999999, // 免費會員無限
				UsedTopologies:          0,
				MaxSimulationsPerDay:    100,
				UsedSimulationsToday:    0,
				LastSimulationResetDate: time.Now().Truncate(24 * time.Hour),
				CanUse3DRendering:       true,
				CanUseAIPrediction:      true,
				CanUseAdvancedSecurity:  true,
				CanAccessAPI:            false,
			}
		}
		if err := h.userRepo.CreateOrUpdateQuota(quota); err != nil {
			// 配額創建失敗不影響登入
			_ = err
		}

		// 生成 JWT token
		token, err := auth.GenerateToken(newUser.ID, newUser.Email, newUser.SubscriptionTier)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, OAuthResponse{
			Token:     token,
			User:      newUser,
			ExpiresIn: 24 * 60 * 60, // 24 小時
		})
		return
	}

	// 已有 OAuth 關聯，取得用戶
	existingUser, err := h.userRepo.GetUserByID(existingOAuth.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user: " + err.Error()})
		return
	}

	// 更新 OAuth token
	existingOAuth.AccessToken = &oauthToken.AccessToken
	if oauthToken.RefreshToken != "" {
		existingOAuth.RefreshToken = &oauthToken.RefreshToken
	}
	if !oauthToken.Expiry.IsZero() {
		existingOAuth.TokenExpiresAt = &oauthToken.Expiry
	}
	if err := h.userRepo.CreateOrUpdateOAuth(existingOAuth); err != nil {
		// OAuth 更新失敗不影響登入
		_ = err
	}

	// 生成 JWT token
	token, err := auth.GenerateToken(existingUser.ID, existingUser.Email, existingUser.SubscriptionTier)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, OAuthResponse{
		Token:     token,
		User:      existingUser,
		ExpiresIn: 24 * 60 * 60, // 24 小時
	})
}

// GetAuthURLRequest 取得授權 URL 請求
type GetAuthURLRequest struct {
	Provider string `json:"provider" binding:"required"` // google, github
	State    string `json:"state,omitempty"`
}

// GetAuthURL 取得 OAuth 授權 URL
func (h *AuthHandler) GetAuthURL(c *gin.Context) {
	var req GetAuthURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	state := req.State
	if state == "" {
		state = "random-state-" + time.Now().Format("20060102150405")
	}

	authURL, err := h.oauthConfig.GetAuthURL(req.Provider, state)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"auth_url": authURL,
		"state":    state,
	})
}

// RefreshTokenRequest 刷新 token 請求
type RefreshTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

// RefreshToken 刷新 token
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newToken, err := auth.RefreshToken(req.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":      newToken,
		"expires_in": 24 * 60 * 60, // 24 小時
	})
}

// GetMe 取得當前用戶資訊
func (h *AuthHandler) GetMe(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	user, err := h.userRepo.GetUserByID(*userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 取得訂閱資訊
	subscription, _ := h.userRepo.GetActiveSubscriptionByUserID(*userID)

	// 取得配額資訊
	quota, _ := h.userRepo.GetQuotaByUserID(*userID)

	c.JSON(http.StatusOK, gin.H{
		"user":        user,
		"subscription": subscription,
		"quota":       quota,
	})
}

