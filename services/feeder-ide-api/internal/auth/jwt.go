package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(getJWTSecret())

// Claims JWT claims 結構
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Tier   string `json:"tier"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT token
func GenerateToken(userID, email, tier string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // 24 小時過期

	claims := &Claims{
		UserID: userID,
		Email:  email,
		Tier:   tier,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "feeder-platform",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken 驗證 JWT token
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// RefreshToken 刷新 token（如果即將過期）
func RefreshToken(tokenString string) (string, error) {
	claims, err := ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	// 如果 token 還有超過 1 小時才過期，不需要刷新
	if time.Until(claims.ExpiresAt.Time) > time.Hour {
		return tokenString, nil
	}

	// 生成新 token
	return GenerateToken(claims.UserID, claims.Email, claims.Tier)
}

func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// 開發環境默認密鑰（生產環境必須設置）
		return "dev-secret-key-change-in-production"
	}
	return secret
}

