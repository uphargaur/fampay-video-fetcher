package repository

import (
	"FamPay-uphar/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
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
func (vr *VideoRepository) FetchVideos(ctx context.Context, limit, offset int64) ([]models.Video, int64, error) {
	opts := options.Find().SetSort(map[string]interface{}{"published_at": -1}).SetSkip(offset).SetLimit(limit)
	cursor, err := vr.db.Find(ctx, map[string]interface{}{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var videos []models.Video
	if err := cursor.All(ctx, &videos); err != nil {
		return nil, 0, err
	}
	totalCount, err := vr.db.CountDocuments(ctx, map[string]interface{}{})
	if err != nil {
		return nil, 0, err
	}

	// Return the videos and the total count
	return videos, totalCount, nil
}
func (vr *VideoRepository) SearchVideos(ctx context.Context, limit, offset int64, query, orderBy, orderType string) ([]models.Video, int64, error) {
	// Map order type to ascending/descending value for MongoDB query
	sortOrder := 1
	if orderType == "d" {
		sortOrder = -1
	}

	// Search query filter
	filter := map[string]interface{}{
		"$or": []interface{}{
			map[string]interface{}{"name": bson.M{"$regex": query, "$options": "i"}},
			map[string]interface{}{"description": bson.M{"$regex": query, "$options": "i"}},
		},
	}

	// Sorting and pagination
	sortOptions := map[string]interface{}{orderBy: sortOrder}
	opts := options.Find().SetSort(sortOptions).SetSkip(offset).SetLimit(limit)

	cursor, err := vr.db.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var videos []models.Video
	if err := cursor.All(ctx, &videos); err != nil {
		return nil, 0, err
	}

	// Count total matching videos in the collection
	totalVideos, err := vr.db.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return videos, totalVideos, nil
}
