package main

import (
	"Api-Picture/controllers"
	"Api-Picture/docs"
	_ "Api-Picture/docs"
	"Api-Picture/middlewares"
	"Api-Picture/models"
	"Api-Picture/repositories"
	"Api-Picture/services"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
)

// @license: Apache 2.0
// @BasePath: /api/v1
// @Schemes: http, https
// @title: Picture API
// @description: picture API
// @version: 1.0.0
// add swagger token bearer
// @securityDefinitions.apikey  API key auth
// @in header
// @name Authorization
// @tokenUrl http://localhost:8080/api/v1/signin
func main() {

	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	docs.SwaggerInfo.Title = "Picture API"
	docs.SwaggerInfo.Description = "picture API"
	db, err := models.Database()
	if err != nil {
		log.Println(err)
	}

	db.DB()

	router := gin.Default()
	// cors middleware
	//router.Use(middlewares.CORSMiddleware())
	v1 := router.Group("/api/v1")
	{
		userController := controllers.NewUserController(services.NewUserService(repositories.NewUserRepository(db)))
		signup := v1.Group("/signup")
		{
			signup.POST("/register", userController.SignUp)
		}
		signin := v1.Group("/signin")
		{
			signin.POST("", userController.SignIn)
		}

		picture := v1.Group("/pictures")
		{
			// let only authenticated user to access this endpoint
			picture.Use(middlewares.JWTAuthMiddleware(db))

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

	log.Println(router.Run(":8080"))
}
