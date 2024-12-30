package config

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Load environment variables
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}
}

// Get MongoDB URI
func GetMongoURI() string {
	return os.Getenv("MONGODB_URI")
}

// Connect to MongoDB
func ConnectDB() (*mongo.Database, error) {
	clientOptions := options.Client().ApplyURI(GetMongoURI())
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return client.Database("fampay"), nil
}

// Get YouTube API keys
func GetAPIKeys() []string {
	return strings.Split(os.Getenv("YOUTUBE_API_KEYS"), ",")
}

// Get application port
func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}
