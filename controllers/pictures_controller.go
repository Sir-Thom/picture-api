package controllers

import (
	"Api-Picture/services"
	"github.com/gin-gonic/gin"
	_ "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
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
//	@Accept			json
//	@Produce		json
//	@Param			id path string true "Picture ID"
//	@Success		200	{object}	models.Pictures
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
//	@Router			/pictures/paginated [get]
func (pc *PictureController) GetPicturesPaginated(ctx *gin.Context) {
	pictures, err := pc.Service.GetPicturesPaginated(0, 12)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, pictures)
}
