package controllers

import (
	"FamPay-uphar/repository"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPaginatedVideos(repo *repository.VideoRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse page and limit from query params
		page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
		limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)

		// Calculate offset for pagination
		offset := (page - 1) * limit

		// Fetch videos and total count
		videos, totalVideos, err := repo.FetchVideos(c, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch videos"})
			return
		}

		// Calculate total pages
		totalPages := int64(0)
		if totalVideos > 0 {
			totalPages = totalVideos / limit
			if totalVideos%limit > 0 {
				totalPages++
			}
		}

		// Generate links for next and previous page
		nextPageLink := ""
		prevPageLink := ""

		if page < totalPages {
			nextPageLink = fmt.Sprintf("%s?page=%d&limit=%d", c.Request.URL.Path, page+1, limit)
		}

		if page > 1 {
			prevPageLink = fmt.Sprintf("%s?page=%d&limit=%d", c.Request.URL.Path, page-1, limit)
		}

		// Response with videos and pagination metadata
		response := gin.H{
			"page":         page,
			"limit":        limit,
			"totalVideos":  totalVideos,
			"totalPages":   totalPages,
			"videos":       videos,
			"nextPageLink": nextPageLink,
			"prevPageLink": prevPageLink,
		}

		// Send the paginated response with links for next/previous pages
		c.JSON(http.StatusOK, response)
	}
}

func SearchVideos(repo *repository.VideoRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse query parameters
		limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
		offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)
		orderBy := c.DefaultQuery("orderBy", "published_at")
		orderType := c.DefaultQuery("orderType", "d")
		query := c.DefaultQuery("q", "")

		// Validate orderType, default to 'd' (descending)
		if orderType != "a" && orderType != "d" {
			orderType = "d"
		}

		if query == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Search query 'q' is required"})
			return
		}

		// Fetch search results from repository
		videos, totalVideos, err := repo.SearchVideos(c, limit, offset, query, orderBy, orderType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch videos"})
			return
		}

		// Calculate total pages for pagination
		totalPages := int64(0)
		if totalVideos > 0 {
			totalPages = totalVideos / limit
			if totalVideos%limit > 0 {
				totalPages++
			}
		}

		// Build pagination response
		response := gin.H{
			"videos":      videos,
			"totalVideos": totalVideos,
			"totalPages":  totalPages,
			"limit":       limit,
			"offset":      offset,
		}

		// Respond with the videos and pagination data
		c.JSON(http.StatusOK, response)
	}
}
