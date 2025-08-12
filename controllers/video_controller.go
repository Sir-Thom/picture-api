package controllers

import (
	"Api-Picture/models"
	"Api-Picture/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type VideoController struct {
	Service *services.VideoService
}

// NewVideoController creates a new VideoController
func NewVideoController(service *services.VideoService) *VideoController {
	return &VideoController{
		Service: service,
	}
}

// GetAllVideos godoc
//
//	@Summary		Get all videos (paginated)
//	@Description	Get paginated list of videos
//	@Tags			videos
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int	false	"Page number"	default(1)
//	@Param			limit	query		int	false	"Results per page"	default(20)
//	@Success		200		{object}	PagedVideoResponse
//	@Router			/videos [get]
func (vc *VideoController) GetAllVideos(ctx *gin.Context) {
	// Parse pagination parameters with defaults
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "20"))

	// Validate parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	videos, total, err := vc.Service.GetVideosPaginated(offset, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if len(videos) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No videos found"})
		return
	}

	// Calculate pagination metadata
	totalPages := (total + int64(limit) - 1) / int64(limit)
	if totalPages == 0 {
		totalPages = 1
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": videos,
		"pagination": gin.H{
			"current_page": page,
			"per_page":     limit,
			"total":        total,
			"total_pages":  totalPages,
		},
	})
}

// GetVideoByName godoc
//
//	@Summary		Get videos by name (paginated)
//	@Description	Search videos by name with fuzzy matching
//	@Tags			videos
//	@Security		Bearer
//	@Accept			json
//	@Produce		json
//	@Param			name	path		string	true	"Video name"
//	@Param			page	query		int		false	"Page number"	default(1)
//	@Param			limit	query		int		false	"Results per page"	default(20)
//	@Success		200		{object}	PagedVideoResponse
//	@Security		APIKeyAuth
//	@Router			/videos/{name} [get]
func (vc *VideoController) GetVideoByName(ctx *gin.Context) {
	name := ctx.Param("name")
	if name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Name parameter is required"})
		return
	}

	// Parse pagination parameters with defaults
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "20"))

	// Validate parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	videos, total, err := vc.Service.GetVideoByNamePaginated(name, offset, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if len(videos) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No videos found"})
		return
	}

	// Calculate pagination metadata
	totalPages := (total + int64(limit) - 1) / int64(limit)
	if totalPages == 0 {
		totalPages = 1
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": videos,
		"pagination": gin.H{
			"current_page": page,
			"per_page":     limit,
			"total":        total,
			"total_pages":  totalPages,
		},
	})
}

type PagedVideoResponse struct {
	Data       []models.Video `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

// PaginationMeta contains pagination information
type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	Total       int64 `json:"total"`
	TotalPages  int64 `json:"total_pages"`
}
