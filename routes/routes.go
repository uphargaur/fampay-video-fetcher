package routes

import (
	"FamPay-uphar/controllers"
	"FamPay-uphar/repository"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitializeRoutes(router *gin.Engine, db *mongo.Database) {
	videoRepo := repository.NewVideoRepository(db)

	api := router.Group("/api")
	{
		api.GET("/videos", controllers.GetPaginatedVideos(videoRepo))
		api.GET("/search", controllers.SearchVideos(videoRepo))
	}
}
