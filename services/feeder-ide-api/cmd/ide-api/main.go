package main

import (
	"log"
	"os"

	"github.com/feeder-platform/feeder-ide-api/api"
	"github.com/feeder-platform/feeder-ide-api/internal/topology"
	"github.com/feeder-platform/feeder-ide-api/internal/profiles"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化資料庫連接（暫時使用記憶體模式）
	topologyRepo := topology.NewInMemoryRepository()
	profileRepo := profiles.NewInMemoryRepository()

	// 初始化 handlers
	topologyHandler := api.NewTopologyHandler(topologyRepo)
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
		// Topology endpoints
		v1.POST("/topologies", topologyHandler.CreateTopology)
		v1.GET("/topologies/:id", topologyHandler.GetTopology)
		v1.PUT("/topologies/:id", topologyHandler.UpdateTopology)
		v1.DELETE("/topologies/:id", topologyHandler.DeleteTopology)
		v1.GET("/topologies", topologyHandler.ListTopologies)

		// Profile endpoints
		v1.GET("/profiles", profileHandler.ListProfiles)
		v1.GET("/profiles/:type", profileHandler.GetProfile)
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

