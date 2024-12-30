package repository

import (
	"FamPay-uphar/models"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type VideoRepository struct {
	db *mongo.Collection
}

func NewVideoRepository(db *mongo.Database) *VideoRepository {
	return &VideoRepository{
		db: db.Collection("videos"),
	}
}

// Save videos to MongoDB
func (vr *VideoRepository) SaveVideos(ctx context.Context, videos []models.Video) error {
	opts := options.Update().SetUpsert(true)
	for _, video := range videos {
		filter := map[string]interface{}{"video_id": video.VideoID}
		update := map[string]interface{}{"$set": video}
		_, err := vr.db.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			return err
		}
	}
	return nil
}

// Fetch videos with pagination
func (vr *VideoRepository) FetchVideos(ctx context.Context, limit, offset int64) ([]models.Video, error) {
	opts := options.Find().SetSort(map[string]interface{}{"published_at": -1}).SetSkip(offset).SetLimit(limit)
	cursor, err := vr.db.Find(ctx, map[string]interface{}{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var videos []models.Video
	if err := cursor.All(ctx, &videos); err != nil {
		return nil, err
	}
	return videos, nil
}
