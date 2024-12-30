package main

import (
	"FamPay-uphar/config"
	"FamPay-uphar/repository"
	"FamPay-uphar/routes"
	"FamPay-uphar/services"
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	config.LoadEnv()

	// Initialize MongoDB
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Start Gin server
	router := gin.Default()

	// Repository and service setup
	videoRepo := repository.NewVideoRepository(db)
	apiKeys := config.GetAPIKeys()
	searchQuery := "official" // Example query
	pollInterval := 10 * time.Second

	videoService := services.NewVideoService(videoRepo, apiKeys, searchQuery, pollInterval)

	// Start video polling in a goroutine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go videoService.StartPolling(ctx)

	// Initialize routes
	routes.InitializeRoutes(router, db)

	// Run server
	port := config.GetPort()
	log.Printf("Server running at http://localhost:%s", port)
	router.Run(":" + port)
}
