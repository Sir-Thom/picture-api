package main

import (
	"Api-Picture/controllers"
	"Api-Picture/docs"
	_ "Api-Picture/docs"
	"Api-Picture/models"
	"Api-Picture/repositories"
	"Api-Picture/services"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
)

// @license: Apache 2.0
func main() {
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.Title = "Picture API"
	docs.SwaggerInfo.Description = "picture API"
	db, err := models.Database()
	if err != nil {
		log.Println(err)
	}
	gin.SetMode(gin.ReleaseMode)
	db.DB()

	router := gin.Default()
	// cors middleware
	router.Use(CORSMiddleware())

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
	err = db.AutoMigrate(&models.Pictures{})
	if err != nil {
		panic(err)
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Fatal(router.Run(":8080"))
}
