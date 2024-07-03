package controllers

import (
	"Api-Picture/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type VideoController struct {
	Service *services.VideoService
}

func NewVideoController(service *services.VideoService) *VideoController {
	return &VideoController{Service: service}
}

// GetAllVideos GetVideo godoc
//
//	@Summary		Get all videos
//	@Description	Get all videos
//	@Tags			videos
//	@Accept			json
//	@Produce		json
//	@Success		200 {object}	models.Video
//	@Router			/videos [get]
func (vc *VideoController) GetAllVideos(ctx *gin.Context) {
	videos, err := vc.Service.GetAllVideo()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(videos) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No videos found"})
		return
	}

	ctx.JSON(http.StatusOK, videos)
}

// GetVideoByName godoc
//
//	@Summary		Get video by name
//	@Description	Get picture by name
//	@Tags			videos
//	@Security		Bearer
//	@Accept			json
//	@Produce		json
//	@Param			name path string true "Video name"
//	@Success		200	{object}	models.Video
//
// @Security API key auth
//
//	@Router			/videos/{name} [get]
func (vc *VideoController) GetVideoByName(ctx *gin.Context) {
	name := ctx.Param("name")
	videos, err := vc.Service.GetVideoByName(name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(videos) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No videos found"})
		return
	}
	println(videos)
	ctx.JSON(http.StatusOK, videos)
}
