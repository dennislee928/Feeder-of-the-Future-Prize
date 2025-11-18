package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT 認證中間件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// 提取 token（格式：Bearer <token>）
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// 將用戶資訊存入 context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_tier", claims.Tier)

		c.Next()
	}
}

// OptionalAuthMiddleware 可選認證中間件（不強制要求認證）
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// 沒有 token，繼續執行（demo 模式）
			c.Next()
			return
		}

		// 嘗試解析 token
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			tokenString := parts[1]
			claims, err := ValidateToken(tokenString)
			if err == nil {
				// Token 有效，設置用戶資訊
				c.Set("user_id", claims.UserID)
				c.Set("user_email", claims.Email)
				c.Set("user_tier", claims.Tier)
			}
		}

		c.Next()
	}
}

// GetUserID 從 context 取得用戶 ID
func GetUserID(c *gin.Context) *string {
	userID, exists := c.Get("user_id")
	if !exists {
		return nil
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return nil
	}

	return &userIDStr
}

// GetUserTier 從 context 取得用戶等級
func GetUserTier(c *gin.Context) string {
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

