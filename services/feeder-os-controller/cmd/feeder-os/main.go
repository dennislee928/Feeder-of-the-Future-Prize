package main

import (
	"log"
	"os"

	"github.com/feeder-platform/feeder-os-controller/api"
	"github.com/feeder-platform/feeder-os-controller/internal/apps"
	"github.com/feeder-platform/feeder-os-controller/internal/bus"
	"github.com/feeder-platform/feeder-os-controller/internal/config"
	"github.com/gin-gonic/gin"
)

func main() {
	// 載入配置
	cfg := config.Load()

	// 初始化 MQTT bus
	mqttBus, err := bus.NewMQTTBus(cfg.MQTT)
	if err != nil {
		log.Fatal("Failed to initialize MQTT bus:", err)
	}
	defer mqttBus.Close()

	// 初始化 app manager
	appManager := apps.NewManager(mqttBus, cfg)

	// 初始化 handlers
	appHandler := api.NewAppHandler(appManager)

	// 設定 Gin router
	router := gin.Default()

	// CORS middleware
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
		// App lifecycle endpoints
		v1.POST("/apps/install", appHandler.InstallApp)
		v1.POST("/apps/enable", appHandler.EnableApp)
		v1.POST("/apps/disable", appHandler.DisableApp)
		v1.GET("/apps", appHandler.ListApps)
		v1.GET("/apps/:id", appHandler.GetApp)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 啟動 server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("Feeder OS Controller starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

