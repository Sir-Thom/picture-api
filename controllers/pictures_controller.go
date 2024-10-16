package controllers

import (
	"Api-Picture/services"
	"compress/gzip"
	"encoding/json"
	"github.com/gin-gonic/gin"
	_ "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	"strconv"
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
	// Get query parameters for lastSeenID and limit

	lastSeenID, err := strconv.Atoi(ctx.Query("last_seen_id"))
	if err != nil {
		// Handle error, or set default value if not provided
		lastSeenID = 0
	}

	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		// Handle error, or set default value if not provided
		limit = 12
	}
	err = ctx.BindQuery(&lastSeenID)
	if err != nil {
		log.Println(err)
		return
	}
	err = ctx.BindQuery(&limit)
	if err != nil {
		log.Println(err)
		return
	}
	pictures, err := pc.Service.GetPicturesPaginated(lastSeenID, limit)
	if err != nil {
		log.Println(err)

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(pictures) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No pictures found"})
		return
	}

	// Compress response using gzip
	gz := gzip.NewWriter(ctx.Writer)
	defer func(gz *gzip.Writer) {
		err := gz.Close()
		if err != nil {
			// Handle error
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
	}(gz)

	ctx.Writer.Header().Set("Content-Encoding", "gzip")
	ctx.Writer.Header().Set("Content-Type", "application/json")

	// Serialize pictures to JSON and write to the compressed response
	if err := json.NewEncoder(gz).Encode(pictures); err != nil {
		// Handle error
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.Status(http.StatusOK)
}
