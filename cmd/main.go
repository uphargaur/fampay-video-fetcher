package main

import (
	"FamPay-uphar/config"
	"FamPay-uphar/repository"
	"FamPay-uphar/routes"
	"FamPay-uphar/services"
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"time"
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
	// CORS configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},                   // Allow requests from your Angular app
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Allowed methods
		AllowHeaders:     []string{"Content-Type", "Authorization"},           // Allowed headers
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
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
