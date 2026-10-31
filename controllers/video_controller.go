package controllers

import (
	"Api-Picture/models"
	"Api-Picture/services"
	"Api-Picture/streaming"
	"errors"
	"io/fs"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const StreamRoutePrefix = "/api/v1/stream/videos/"

type VideoController struct {
	Service *services.VideoService
	Stream  *streaming.Server
}

// NewVideoController creates a new VideoController
func NewVideoController(service *services.VideoService, stream *streaming.Server) *VideoController {
	return &VideoController{
		Service: service,
		Stream:  stream,
	}
}

// withStreamURLs fills the virtual stream_url field for each video.
func (vc *VideoController) withStreamURLs(videos []models.Video) []models.Video {
	if vc.Stream == nil {
		return videos
	}
	for i := range videos {
		videos[i].StreamURL = StreamRoutePrefix + strconv.FormatUint(uint64(videos[i].ID), 10)
	}
	return videos
}

// StreamVideo godoc
//
//	@Summary		Stream a video
//	@Description	Streams the video file with HTTP Range support (seeking, resume). Authenticate with the Bearer header or the httpOnly "token" cookie, so it works directly as a <video src>.
//	@Tags			videos
//	@Produce		video/mp4
//	@Param			id		path	int		true	"Video ID"
//	@Param			Range	header	string	false	"Byte range, e.g. bytes=0-1048575"
//	@Success		200
//	@Success		206
//	@Failure		400,404,416
//	@Security		BearerAuth
//	@Router			/stream/videos/{id} [get]
func (vc *VideoController) StreamVideo(ctx *gin.Context) {
	if vc.Stream == nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "Video streaming is not configured"})
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil || id == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
		return
	}

	stored, err := vc.Service.GetVideoPathByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if err := vc.Stream.Serve(ctx.Writer, ctx.Request, stored); err != nil {
		switch {
		case errors.Is(err, fs.ErrNotExist), errors.Is(err, streaming.ErrNotAFile):
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Video file not found"})
		case errors.Is(err, streaming.ErrOutsideRoot):
			log.Printf("video %d: stored path %q is outside the storage root", id, stored)
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Video file not found"})
		default:
			log.Printf("video %d: %v", id, err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not open video"})
		}
		return
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
