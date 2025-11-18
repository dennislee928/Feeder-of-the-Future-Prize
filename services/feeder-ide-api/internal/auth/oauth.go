package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/github"
)

// OAuthConfig OAuth 配置
type OAuthConfig struct {
	Google *oauth2.Config
	GitHub *oauth2.Config
}

// NewOAuthConfig 建立 OAuth 配置
func NewOAuthConfig() *OAuthConfig {
	config := &OAuthConfig{}

	// Google OAuth
	googleClientID := os.Getenv("OAUTH_GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("OAUTH_GOOGLE_CLIENT_SECRET")
	if googleClientID != "" && googleClientSecret != "" {
		config.Google = &oauth2.Config{
			ClientID:     googleClientID,
			ClientSecret: googleClientSecret,
			RedirectURL:  os.Getenv("OAUTH_GOOGLE_REDIRECT_URL"),
			Scopes:       []string{"openid", "profile", "email"},
			Endpoint:     google.Endpoint,
		}
	}

	// GitHub OAuth
	githubClientID := os.Getenv("OAUTH_GITHUB_CLIENT_ID")
	githubClientSecret := os.Getenv("OAUTH_GITHUB_CLIENT_SECRET")
	if githubClientID != "" && githubClientSecret != "" {
		config.GitHub = &oauth2.Config{
			ClientID:     githubClientID,
			ClientSecret: githubClientSecret,
			RedirectURL:  os.Getenv("OAUTH_GITHUB_REDIRECT_URL"),
			Scopes:       []string{"user:email"},
			Endpoint:     github.Endpoint,
		}
	}

	return config
}

// GetAuthURL 取得 OAuth 授權 URL
func (c *OAuthConfig) GetAuthURL(provider string, state string) (string, error) {
	var config *oauth2.Config

	switch provider {
	case "google":
		if c.Google == nil {
			return "", fmt.Errorf("Google OAuth not configured")
		}
		config = c.Google
	case "github":
		if c.GitHub == nil {
			return "", fmt.Errorf("GitHub OAuth not configured")
		}
		config = c.GitHub
	default:
		return "", fmt.Errorf("unsupported provider: %s", provider)
	}

	return config.AuthCodeURL(state), nil
}

// ExchangeCode 交換授權碼取得 token
func (c *OAuthConfig) ExchangeCode(provider, code string) (*oauth2.Token, error) {
	var config *oauth2.Config

	switch provider {
	case "google":
		if c.Google == nil {
			return nil, fmt.Errorf("Google OAuth not configured")
		}
		config = c.Google
	case "github":
		if c.GitHub == nil {
			return nil, fmt.Errorf("GitHub OAuth not configured")
		}
		config = c.GitHub
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	ctx := context.Background()
	token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	return token, nil
}

// GetUserInfo 取得用戶資訊
func (c *OAuthConfig) GetUserInfo(provider string, token *oauth2.Token) (*OAuthUserInfo, error) {
	switch provider {
	case "google":
		return c.getGoogleUserInfo(token)
	case "github":
		return c.getGitHubUserInfo(token)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

// OAuthUserInfo OAuth 用戶資訊
type OAuthUserInfo struct {
	Provider       string
	ProviderUserID string
	Email          string
	Name           string
	AvatarURL      string
}

// getGoogleUserInfo 取得 Google 用戶資訊
func (c *OAuthConfig) getGoogleUserInfo(token *oauth2.Token) (*OAuthUserInfo, error) {
	client := c.Google.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var userInfo struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}

	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	return &OAuthUserInfo{
		Provider:       "google",
		ProviderUserID: userInfo.ID,
		Email:          userInfo.Email,
		Name:           userInfo.Name,
		AvatarURL:      userInfo.Picture,
	}, nil
}

// getGitHubUserInfo 取得 GitHub 用戶資訊
func (c *OAuthConfig) getGitHubUserInfo(token *oauth2.Token) (*OAuthUserInfo, error) {
	client := c.GitHub.Client(context.Background(), token)

	// 取得用戶基本資訊
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var userInfo struct {
		ID    int    `json:"id"`
		Login string `json:"login"`
		Name  string `json:"name"`
		Email string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	}

	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	// 如果 email 為空，嘗試從 email API 取得
	if userInfo.Email == "" {
		emailResp, err := client.Get("https://api.github.com/user/emails")
		if err == nil {
			defer emailResp.Body.Close()
			emailBody, _ := io.ReadAll(emailResp.Body)

			var emails []struct {
				Email   string `json:"email"`
				Primary bool   `json:"primary"`
			}

			if json.Unmarshal(emailBody, &emails) == nil {
				for _, email := range emails {
					if email.Primary {
						userInfo.Email = email.Email
						break
					}
				}
			}
		}
	}

	return &OAuthUserInfo{
		Provider:       "github",
		ProviderUserID: fmt.Sprintf("%d", userInfo.ID),
		Email:          userInfo.Email,
		Name:           userInfo.Name,
		AvatarURL:      userInfo.AvatarURL,
	}, nil
}

