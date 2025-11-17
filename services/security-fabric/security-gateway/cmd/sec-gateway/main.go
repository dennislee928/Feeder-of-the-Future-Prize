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

	// 記錄所有請求
	router.Use(proxyHandler.LoggingMiddleware())

	// mTLS 驗證 middleware（可選）
	if mtlsManager != nil {
		router.Use(proxyHandler.MTLSMiddleware())
	}

	// Proxy routes
	router.Any("/api/*path", proxyHandler.ProxyRequest)
	router.Any("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 啟動 server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8443"
	}

	log.Printf("Security Gateway starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

