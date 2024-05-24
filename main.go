package main

import (
	"Api-Picture/controllers"
	"Api-Picture/docs"
	_ "Api-Picture/docs"
	"Api-Picture/middlewares"
	"Api-Picture/models"
	"Api-Picture/repositories"
	"Api-Picture/services"
	"Api-Picture/utils"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
)

// @license: Apache 2.0
func main() {
	gin.SetMode(gin.ReleaseMode)
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.Title = "Picture API"
	docs.SwaggerInfo.Description = "picture API"
	db, err := models.Database()
	if err != nil {
		log.Println(err)
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Cron job panicked: %v", r)
			}
		}()
		utils.CronJob()
	}()

	db.DB()

	router := gin.Default()
	// cors middleware
	router.Use(middlewares.CORSMiddleware())
	// jwt middleware
	//router.Use(middlewares.JWTAuthMiddleware(db, "secret"))

	v1 := router.Group("/api/v1")
	{
		signin := v1.Group("/signup")
		{
			userController := controllers.NewUserController(services.NewUserService(repositories.NewUserRepository(db)))
			signin.POST("/register", userController.SignUp)
		}

		picture := v1.Group(
			"/pictures")
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
	if gin.Mode() == gin.ReleaseMode {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	log.Fatal(router.Run(":8080"))
}
