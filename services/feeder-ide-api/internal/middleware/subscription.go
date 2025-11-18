package middleware

import (
	"net/http"
	"strings"

	"github.com/feeder-platform/feeder-ide-api/internal/auth"
	"github.com/feeder-platform/feeder-ide-api/internal/user"
	"github.com/gin-gonic/gin"
)

// SubscriptionMiddleware 訂閱等級檢查中間件
func SubscriptionMiddleware(requiredTier string, userService *user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := auth.GetUserID(c)
		tier, err := userService.GetUserTier(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user tier"})
			c.Abort()
			return
		}

		// 檢查等級是否符合要求
		if !isTierAllowed(tier, requiredTier) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Insufficient subscription tier",
				"required": requiredTier,
				"current": tier,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// isTierAllowed 檢查用戶等級是否允許訪問
func isTierAllowed(userTier, requiredTier string) bool {
	tierLevels := map[string]int{
		"demo":   0,
		"free":   1,
		"premium": 2,
	}

	userLevel, userOk := tierLevels[userTier]
	requiredLevel, requiredOk := tierLevels[requiredTier]

	if !userOk || !requiredOk {
		return false
	}

	return userLevel >= requiredLevel
}

// FeatureMiddleware 功能權限檢查中間件
func FeatureMiddleware(feature string, userService *user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := auth.GetUserID(c)
		canUse, err := userService.CanUseFeature(userID, feature)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check feature permission"})
			c.Abort()
			return
		}

		if !canUse {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Feature not available for your subscription tier",
				"feature": feature,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// QuotaMiddleware 配額檢查中間件
func QuotaMiddleware(quotaType string, userService *user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := auth.GetUserID(c)

		switch quotaType {
		case "topology":
			canCreate, used, max, err := userService.CheckTopologyQuota(userID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check topology quota"})
				c.Abort()
				return
			}

			if !canCreate {
				c.JSON(http.StatusForbidden, gin.H{
					"error": "Topology quota exceeded",
					"used": used,
					"max": max,
				})
				c.Abort()
				return
			}

			// 將配額資訊存入 context，供 handler 使用
			c.Set("quota_used", used)
			c.Set("quota_max", max)

		case "simulation":
			canSimulate, err := userService.CheckSimulationQuota(userID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check simulation quota"})
				c.Abort()
				return
			}

			if !canSimulate {
				c.JSON(http.StatusForbidden, gin.H{
					"error": "Simulation quota exceeded for today",
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// OptionalSubscriptionMiddleware 可選訂閱檢查（不強制要求，但會設置 context）
func OptionalSubscriptionMiddleware(userService *user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := auth.GetUserID(c)
		if userID != nil {
			tier, _ := userService.GetUserTier(userID)
			c.Set("user_tier", tier)

			quota, _ := userService.GetUserQuota(*userID)
			if quota != nil {
				c.Set("user_quota", quota)
			}
		} else {
			c.Set("user_tier", "demo")
		}

		c.Next()
	}
}

// GetUserTierFromContext 從 context 取得用戶等級
func GetUserTierFromContext(c *gin.Context) string {
	tier, exists := c.Get("user_tier")
	if !exists {
		return "demo"
	}

	tierStr, ok := tier.(string)
	if !ok {
		return "demo"
	}

	return tierStr
}

// GetUserQuotaFromContext 從 context 取得用戶配額
func GetUserQuotaFromContext(c *gin.Context) *user.UserQuota {
	quota, exists := c.Get("user_quota")
	if !exists {
		return nil
	}

	userQuota, ok := quota.(*user.UserQuota)
	if !ok {
		return nil
	}

	return userQuota
}

// RequireAuthOrDemo 要求認證或允許 demo 模式
func RequireAuthOrDemo() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// 允許 demo 模式，設置為 demo 用戶
			c.Set("user_id", nil)
			c.Set("user_tier", "demo")
			c.Next()
			return
		}

		// 嘗試解析 token
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			claims, err := auth.ValidateToken(parts[1])
			if err == nil {
				c.Set("user_id", claims.UserID)
				c.Set("user_email", claims.Email)
				c.Set("user_tier", claims.Tier)
			}
		}

		c.Next()
	}
}

