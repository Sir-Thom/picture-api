package main

import (
	"Api-Picture/controllers"
	_ "Api-Picture/docs"
	"Api-Picture/models"
	"Api-Picture/repositories"
	"Api-Picture/services"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
)

//	@title			Picture API
//	@version		2.0
//	@description	picture API
// @host		127.0.0.1:8080
// @BasePath	/api/v1

func main() {

	db, err := models.Database()
	if err != nil {
		log.Println(err)
	}
	db.DB()

	router := gin.Default()

	v1 := router.Group("/api/v1")
	{

		picture := v1.Group("/pictures")
		{

			pictureController := controllers.NewPictureController(services.NewPictureService(repositories.NewPictureRepository(db)))
			picture.GET("", pictureController.GetPictures)
			picture.GET("/:id", pictureController.GetPictureById)
			picture.GET("/count", pictureController.CountPicture)
			picture.GET("/paginated", pictureController.GetPicturesPaginated)
		}

	}
	// integer swagger

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	ginSwagger.URL("http://127.0.0.1/swagger/doc.json")

	log.Fatal(router.Run("127.0.0.1:8080"))
}
