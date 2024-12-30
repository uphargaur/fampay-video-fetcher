package models

import "time"

// Video represents a YouTube video stored in MongoDB
type Video struct {
	ID          string    `bson:"_id,omitempty"`
	VideoID     string    `bson:"video_id"`
	Title       string    `bson:"title"`
	Description string    `bson:"description"`
	PublishedAt time.Time `bson:"published_at"`
	Thumbnail   string    `bson:"thumbnail"`
	CreatedAt   time.Time `bson:"created_at"`
}
