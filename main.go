package main

import (
	"Api-Picture/controllers"
	"Api-Picture/docs"
	"Api-Picture/middlewares"
	"Api-Picture/models"
	"Api-Picture/repositories"
	"Api-Picture/services"
	//"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	//"time"
)

// @license: Apache 2.0
// @BasePath: /api/v1
// @Schemes: http, https
// @title: Picture API
// @description: picture API
// @version: 1.0.0
// add swagger token bearer
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
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Cors middleware
	router.Use(gin.Recovery())
	/*corsConfig := cors.Config{

		AllowOrigins:        []string{"http://*"},
		AllowMethods:        []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:        []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"},
		ExposeHeaders:       []string{"Content-Length"},
		AllowCredentials:    true,
		AllowWildcard:       true,
		AllowPrivateNetwork: true,
		MaxAge:              12 * time.Hour,
	}*/
	//	router.Use(cors.New(corsConfig))
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	v1 := router.Group("api/v1")
	{
		userController := controllers.NewUserController(services.NewUserService(repositories.NewUserRepository(db)))
		signup := v1.Group("signup")
		{
			signup.POST("register", userController.SignUp)
		}
		signin := v1.Group("signin")
		{
			signin.POST("", userController.SignIn)
		}

		picture := v1.Group("pictures")
		{
			picture.Use(middlewares.JWTAuthMiddleware(db))

			pictureController := controllers.NewPictureController(services.NewPictureService(repositories.NewPictureRepository(db)))
			picture.GET("", pictureController.GetPictures)
			picture.GET(":id", pictureController.GetPictureById)
			picture.GET("count", pictureController.CountPicture)
			picture.GET("paginated", pictureController.GetPicturesPaginated)
		}
		health := v1.Group("health")
		{
			health.GET("ping", func(c *gin.Context) {
				db, err := models.Database()
				if err != nil {
					c.JSON(500, gin.H{
						"message": "error",
					})
					return
				}
				db.DB()
				c.JSON(200, gin.H{
					"message": "ok",
				})
			})

		}
	}

	// Run database migrations
	err = db.AutoMigrate(&models.Pictures{})
	if err != nil {
		panic(err)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println(router.Run(":8080"))
}
