package controllers

import (
	"Api-Picture/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type SeriesController struct {
	Service *services.SeriesService
}

func NewSeriesController(service *services.SeriesService) *SeriesController {
	return &SeriesController{Service: service}
}

func (sc *SeriesController) GetSeriesByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	series, err := sc.Service.GetSeriesByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, series)
}

// GetAllSeries GetSeries godoc
//
//	@Summary		Get all series
//	@Description	Get all series
//	@Tags			series
//	@Accept			json
//	@Produce		json
//	@Success		200 {object}	models.Series
//	@Router			/series [get]
func (sc *SeriesController) GetAllSeries(ctx *gin.Context) {
	series, err := sc.Service.GetAllSeries()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(series) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No series found"})
		return
	}
	ctx.JSON(http.StatusOK, series)
}

// GetSeriesByName godoc
//
//	@Summary		Get series by name
//	@Description	Get series by name
//	@Tags			series
//	@Security		Bearer
//	@Accept			json
//	@Produce		json
//	@Param			name path string true "series name"
//	@Success		200	{object}	models.Series
//
// @Security API key auth
//
//	@Router			/series/{name} [get]
func (sc *SeriesController) GetSeriesByName(ctx *gin.Context) {
	name := ctx.Param("name")
	videos, err := sc.Service.GetSeriesByName(name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(videos) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No series found"})
		return
	}
	println(videos)
	ctx.JSON(http.StatusOK, videos)
}
