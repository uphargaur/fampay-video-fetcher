package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const youtubeAPIURL = "https://www.googleapis.com/youtube/v3/search"

type VideoData struct {
	VideoID     string    `json:"video_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	Thumbnail   string    `json:"thumbnail"`
}

func FetchYouTubeVideos(apiKey, query string, publishedAfter time.Time) ([]VideoData, error) {
	url := fmt.Sprintf("%s?part=snippet&type=video&order=date&q=%s&publishedAfter=%s&key=%s",
		youtubeAPIURL, query, publishedAfter.Format(time.RFC3339), apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch YouTube data: %s", resp.Status)
	}

	var data struct {
		Items []struct {
			ID struct {
				VideoID string `json:"videoId"`
			} `json:"id"`
			Snippet struct {
				Title       string    `json:"title"`
				Description string    `json:"description"`
				PublishedAt time.Time `json:"publishedAt"`
				Thumbnails  struct {
					Default struct {
						URL string `json:"url"`
					} `json:"default"`
				} `json:"thumbnails"`
			} `json:"snippet"`
		} `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var videos []VideoData
	for _, item := range data.Items {
		videos = append(videos, VideoData{
			VideoID:     item.ID.VideoID,
			Title:       item.Snippet.Title,
			Description: item.Snippet.Description,
			PublishedAt: item.Snippet.PublishedAt,
			Thumbnail:   item.Snippet.Thumbnails.Default.URL,
		})
	}

	return videos, nil
}
