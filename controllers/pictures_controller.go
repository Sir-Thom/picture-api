package controllers

import (
	"Api-Picture/services"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/swaggo/gin-swagger"
)

type PictureController struct {
	Service *services.PictureService
}

func NewPictureController(service *services.PictureService) *PictureController {
	return &PictureController{Service: service}
}

// GetPictures godoc
//
//	@Summary		Get all pictures
//	@Description	Get all pictures
//	@Tags			pictures
//	@Accept			json
//	@Produce		json
//	@Success		200 {object}	models.Pictures
//	@Router			/pictures [get]
func (pc *PictureController) GetPictures(ctx *gin.Context) {
	limit := 10 // default limit
	pictures, err := pc.Service.GetAllPictures(limit)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, pictures)
}

// GetPictureById godoc
//
//	@Summary		Get picture by id
//	@Description	Get picture by id
//	@Tags			pictures
//	@Security		Bearer
//	@Accept			json
//	@Produce		json
//	@Param			id path string true "Picture ID"
//	@Success		200	{object}	models.Pictures
//
// @Security API key auth
//
//	@Router			/pictures/{id} [get]
func (pc *PictureController) GetPictureById(ctx *gin.Context) {
	id := ctx.Param("id")
	picture, err := pc.Service.GetPictureById(id)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, picture)
}

// CountPicture godoc
//
//	@Summary		Count pictures
//	@Description	Count pictures
//	@Tags			pictures
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	int64
//	@Router			/pictures/count [get]
//
// @Security API key auth
func (pc *PictureController) CountPicture(ctx *gin.Context) {
	count, err := pc.Service.CountPictures()
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, count)
}

// GetPicturesPaginated godoc
//
//	@Summary		Get pictures paginated
//	@Description	Get pictures paginated
//	@Tags			pictures
//	@Accept			json
//	@Produce		json
//	@Success		200 {object}	models.Pictures
//	@Param			last_seen_id query int false "Page number" default(1)
//	@Param			limit query int false "Limit per page"	default(12)
//	@Router			/pictures/paginated [get]
func (pc *PictureController) GetPicturesPaginated(ctx *gin.Context) {
	lastSeenID, err := strconv.Atoi(ctx.Query("last_seen_id"))
	if err != nil {
		lastSeenID = 0
	}

	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		limit = 12
	}
	if limit < 1 || limit > 100 {
		limit = 12
	}

	pictures, err := pc.Service.GetPicturesPaginated(lastSeenID, limit)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	if len(pictures) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No pictures found"})
		return
	}

	ctx.JSON(http.StatusOK, pictures)
}
