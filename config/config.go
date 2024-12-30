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

// LoadEnv loads environment variables
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}
}

// GetMongoURI retrieves the MongoDB Atlas connection URI
func GetMongoURI() string {
	return os.Getenv("MONGODB_URI")
}

// ConnectDB connects to the MongoDB Atlas instance
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

	// Ping MongoDB to verify the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	log.Println("Successfully connected to MongoDB Atlas")
	return client.Database("fampay"), nil
}

func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}

func GetAPIKeys() []string {
	apiKeys := os.Getenv("YOUTUBE_API_KEYS")
	if apiKeys == "" {
		return []string{} // Return an empty slice if no keys are provided
	}
	return strings.Split(apiKeys, ",")
}
