package main

import (
	"log"
	"os"

	"github.com/feeder-platform/security-gateway/internal/proxy"
	"github.com/feeder-platform/security-gateway/internal/mTLS"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化 mTLS（簡化版本，使用自簽證書）
	mtlsManager, err := mTLS.NewManager()
	if err != nil {
		log.Printf("Warning: mTLS not fully configured: %v", err)
	}

	// 初始化反向 proxy
	proxyHandler := proxy.NewHandler(mtlsManager)

	// 設定 Gin router
	router := gin.Default()

	// Health check (必須在所有中間件之前定義，確保不被攔截)
	router.Any("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 記錄所有請求
	router.Use(proxyHandler.LoggingMiddleware())

	// mTLS 驗證 middleware（可選，但跳過 /health）
	if mtlsManager != nil {
		router.Use(func(c *gin.Context) {
			// 跳過 /health 端點的 mTLS 驗證
			if c.Request.URL.Path == "/health" {
				c.Next()
				return
			}
			// 對其他端點應用 mTLS 驗證
			proxyHandler.MTLSMiddleware()(c)
		})
	}

	// Proxy routes
	router.Any("/api/*path", proxyHandler.ProxyRequest)

	// 啟動 server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8443"
	}

	// Render 需要綁定到 0.0.0.0 才能檢測到端口
	addr := "0.0.0.0:" + port
	log.Printf("Security Gateway starting on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

