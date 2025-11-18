package main

import (
	"log"
	"os"

	"github.com/feeder-platform/feeder-ide-api/api"
	"github.com/feeder-platform/feeder-ide-api/internal/auth"
	"github.com/feeder-platform/feeder-ide-api/internal/database"
	"github.com/feeder-platform/feeder-ide-api/internal/middleware"
	"github.com/feeder-platform/feeder-ide-api/internal/payment"
	"github.com/feeder-platform/feeder-ide-api/internal/topology"
	"github.com/feeder-platform/feeder-ide-api/internal/profiles"
	"github.com/feeder-platform/feeder-ide-api/internal/user"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化資料庫連接
	var topologyRepo topology.Repository
	var err error

	// 檢查是否有 DATABASE_URL，如果有則使用 PostgreSQL，否則使用記憶體模式
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL != "" {
		// 初始化 PostgreSQL
		if err := database.Init(); err != nil {
			log.Fatalf("Failed to initialize database: %v", err)
		}
		defer database.Close()

		topologyRepo, err = topology.NewPostgresRepository()
		if err != nil {
			log.Fatalf("Failed to create postgres repository: %v", err)
		}
		log.Println("Using PostgreSQL database")
	} else {
		// 使用記憶體模式（開發/測試用）
		topologyRepo = topology.NewInMemoryRepository()
		log.Println("Using in-memory database (development mode)")
	}

	profileRepo := profiles.NewInMemoryRepository()

	// 初始化用戶相關服務（僅在 PostgreSQL 模式下）
	var userRepo user.Repository
	var userService *user.Service
	var authHandler *api.AuthHandler
	var paymentHandler *api.PaymentHandler
	var oauthConfig *auth.OAuthConfig

	if databaseURL != "" {
		// 初始化用戶 repository
		userRepo, err = user.NewPostgresUserRepository()
		if err != nil {
			log.Fatalf("Failed to create user repository: %v", err)
		}

		// 初始化用戶服務
		userService = user.NewService(userRepo)
		// 設置拓樸計數器（用於檢查配額）
		userService.SetTopologyCounter(topologyRepo)

		// 初始化 OAuth 配置
		oauthConfig = auth.NewOAuthConfig()

		// 初始化認證處理器
		authHandler = api.NewAuthHandler(oauthConfig, userRepo, userService)

		// 初始化付費服務
		stripeService := payment.NewStripeService()
		paypalService := payment.NewPayPalService()
		webhookHandler := payment.NewWebhookHandler(stripeService, paypalService, userRepo, userService)
		paymentHandler = api.NewPaymentHandler(stripeService, paypalService, webhookHandler, userRepo, userService)
	}

	// 初始化 handlers
	var topologyHandler *api.TopologyHandler
	if userService != nil {
		topologyHandler = api.NewTopologyHandler(topologyRepo, userService)
	} else {
		topologyHandler = api.NewTopologyHandler(topologyRepo, nil)
	}
	profileHandler := api.NewProfileHandler(profileRepo)

	// 設定 Gin router
	router := gin.Default()

	// CORS middleware（開發環境用）
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API routes
	v1 := router.Group("/api/v1")
	{
		// Auth endpoints (僅在資料庫模式下可用)
		if authHandler != nil {
			authGroup := v1.Group("/auth")
			{
				authGroup.POST("/oauth/callback", authHandler.OAuthCallback)
				authGroup.POST("/oauth/url", authHandler.GetAuthURL)
				authGroup.POST("/refresh", authHandler.RefreshToken)
				authGroup.GET("/me", auth.AuthMiddleware(), authHandler.GetMe)
			}
		}

		// Topology endpoints (使用可選認證中間件和配額檢查)
		if authHandler != nil && userService != nil {
			v1.Use(auth.OptionalAuthMiddleware())
			// 為創建拓樸添加配額檢查
			v1.POST("/topologies", middleware.QuotaMiddleware("topology", userService), topologyHandler.CreateTopology)
		} else {
			v1.POST("/topologies", topologyHandler.CreateTopology)
		}
		v1.GET("/topologies/:id", topologyHandler.GetTopology)
		v1.PUT("/topologies/:id", topologyHandler.UpdateTopology)
		v1.DELETE("/topologies/:id", topologyHandler.DeleteTopology)
		v1.GET("/topologies", topologyHandler.ListTopologies)

		// Profile endpoints
		v1.GET("/profiles", profileHandler.ListProfiles)
		v1.GET("/profiles/:type", profileHandler.GetProfile)

		// Payment endpoints (僅在資料庫模式下可用)
		if paymentHandler != nil {
			payments := v1.Group("/payments")
			payments.Use(auth.AuthMiddleware())
			{
				payments.POST("/create-checkout", paymentHandler.CreateCheckout)
				payments.GET("/history", paymentHandler.GetPaymentHistory)
			}

			// Webhook 端點（不需要認證）
			v1.POST("/payments/webhook/stripe", paymentHandler.HandleStripeWebhook)
			v1.POST("/payments/webhook/paypal", paymentHandler.HandlePayPalWebhook)
		}
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 啟動 server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

