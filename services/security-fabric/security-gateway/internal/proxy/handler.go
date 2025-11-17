package proxy

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/feeder-platform/security-gateway/internal/mTLS"
	"github.com/gin-gonic/gin"
)

// Handler 處理反向 proxy 邏輯
type Handler struct {
	mtlsManager *mTLS.Manager
	client      *http.Client
}

// NewHandler 建立新的 proxy handler
func NewHandler(mtlsManager *mTLS.Manager) *Handler {
	return &Handler{
		mtlsManager: mtlsManager,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// LoggingMiddleware 記錄所有請求
func (h *Handler) LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// 讀取 request body（用於記錄）
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// 處理請求
		c.Next()

		// 記錄
		latency := time.Since(start)
		log.Printf("[SECURITY] %s %s | Status: %d | Latency: %v | IP: %s",
			method, path, c.Writer.Status(), latency, c.ClientIP())

		// 記錄敏感操作
		if h.isSensitiveOperation(method, path) {
			log.Printf("[SECURITY-ALERT] Sensitive operation detected: %s %s from %s",
				method, path, c.ClientIP())
		}
	}
}

// MTLSMiddleware mTLS 驗證 middleware
func (h *Handler) MTLSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if h.mtlsManager == nil {
			c.Next()
			return
		}

		// 簡化版本：檢查 client certificate
		if !h.mtlsManager.VerifyClient(c.Request) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "mTLS authentication required"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// ProxyRequest 代理請求到後端服務
func (h *Handler) ProxyRequest(c *gin.Context) {
	path := c.Param("path")
	
	// 決定後端服務（簡化：根據路徑判斷）
	targetURL := h.determineTarget(path)
	if targetURL == "" {
		c.JSON(http.StatusBadGateway, gin.H{"error": "No target service available"})
		return
	}

	// 建立新請求
	req, err := http.NewRequest(c.Request.Method, targetURL+path, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 複製 headers
	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// 發送請求
	resp, err := h.client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	// 複製 response headers
	for key, values := range resp.Header {
		for _, value := range values {
			c.Writer.Header().Add(key, value)
		}
	}

	// 複製 response body
	c.Writer.WriteHeader(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)
}

// determineTarget 決定目標服務 URL
func (h *Handler) determineTarget(path string) string {
	// 簡化：根據路徑前綴判斷
	if len(path) > 0 {
		// 這裡可以根據實際需求路由到不同服務
		// 例如：/topologies -> feeder-ide-api, /apps -> feeder-os-controller
		return "http://feeder-ide-api:8080"
	}
	return ""
}

// isSensitiveOperation 判斷是否為敏感操作
func (h *Handler) isSensitiveOperation(method, path string) bool {
	// 定義敏感操作
	sensitiveMethods := map[string]bool{
		"POST":   true,
		"PUT":    true,
		"DELETE": true,
		"PATCH":  true,
	}

	sensitivePaths := []string{
		"/apps/install",
		"/apps/enable",
		"/apps/disable",
		"/topologies",
		"/commands",
	}

	if sensitiveMethods[method] {
		for _, sensitivePath := range sensitivePaths {
			if len(path) >= len(sensitivePath) && path[:len(sensitivePath)] == sensitivePath {
				return true
			}
		}
	}

	return false
}

