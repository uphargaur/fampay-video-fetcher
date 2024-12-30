package controllers

import (
	"FamPay-uphar/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPaginatedVideos(repo *repository.VideoRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
		limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
		offset := (page - 1) * limit

		videos, err := repo.FetchVideos(c, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch videos"})
			return
		}
		c.JSON(http.StatusOK, videos)
	}
}
