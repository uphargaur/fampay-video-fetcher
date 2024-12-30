package services

import (
	"FamPay-uphar/models"
	"FamPay-uphar/repository"
	"FamPay-uphar/utils"
	"context"
	"log"
	"sync"
	"time"
)

type VideoService struct {
	Repo       *repository.VideoRepository
	APIKeys    []string
	Query      string
	Interval   time.Duration
	currentKey int
	mu         sync.Mutex
}

func NewVideoService(repo *repository.VideoRepository, apiKeys []string, query string, interval time.Duration) *VideoService {
	return &VideoService{
		Repo:     repo,
		APIKeys:  apiKeys,
		Query:    query,
		Interval: interval,
	}
}

// Fetch and store YouTube videos
func (vs *VideoService) FetchAndStoreVideos(ctx context.Context) {
	publishedAfter := time.Now().Add(-1 * time.Hour).UTC() // Fetch videos from the last 24 hours
	apiKey := vs.GetAPIKey()

	videos, err := utils.FetchYouTubeVideos(apiKey, vs.Query, publishedAfter)
	if err != nil {
		log.Printf("Error fetching YouTube videos: %v", err)
		if vs.RotateAPIKey() {
			log.Println("Switched to the next API key due to quota exhaustion.")
		}
		return
	}

	var videoModels []models.Video
	for _, video := range videos {
		videoModels = append(videoModels, models.Video{
			VideoID:     video.VideoID,
			Title:       video.Title,
			Description: video.Description,
			PublishedAt: video.PublishedAt,
			Thumbnail:   video.Thumbnail,
			CreatedAt:   time.Now(),
		})
	}

	// Store videos in MongoDB
	err = vs.Repo.SaveVideos(ctx, videoModels)
	if err != nil {
		log.Printf("Error saving videos to MongoDB: %v", err)
	}
}

// Start polling YouTube API at regular intervals
func (vs *VideoService) StartPolling(ctx context.Context) {
	ticker := time.NewTicker(vs.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping video polling...")
			return
		case <-ticker.C:
			log.Println("Fetching and storing YouTube videos...")
			vs.FetchAndStoreVideos(ctx)
		}
	}
}

// Get the current API key
func (vs *VideoService) GetAPIKey() string {
	vs.mu.Lock()
	defer vs.mu.Unlock()
	return vs.APIKeys[vs.currentKey]
}

// Rotate to the next API key (circular rotation)
func (vs *VideoService) RotateAPIKey() bool {
	vs.mu.Lock()
	defer vs.mu.Unlock()

	if len(vs.APIKeys) <= 1 {
		return false
	}
	vs.currentKey = (vs.currentKey + 1) % len(vs.APIKeys)
	return true
}
