package main

import (
	"Api-Picture/controllers"
	"Api-Picture/docs"
	"Api-Picture/middlewares"
	"Api-Picture/models"
	"Api-Picture/repositories"
	"Api-Picture/services"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/time/rate"
)

// @license.name Apache 2.0
// @BasePath /api/v1
// @Schemes http https
// @title Picture API
// @description Picture and Video API
// @version 1.1.0
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @tokenUrl http://localhost:8080/api/v1/signin
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env vars")
	}

	// Initialize rate limiter for authentication routes
	authLimiter := middlewares.NewRateLimiter(rate.Every(12*time.Second), 5, 10*time.Minute)

	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	docs.SwaggerInfo.Title = "Picture API"
	docs.SwaggerInfo.Description = "picture API"

	db, err := models.Database()
	if err != nil {
		log.Println(err)
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(middlewares.CORSMiddleware())
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	v1 := router.Group("api/v1")
	{
		jwtService := services.NewJWTService(os.Getenv("SECRET_KEY"), 90*24*time.Hour)
		userController := controllers.NewUserController(services.NewUserService(repositories.NewUserRepository(db), jwtService))

		signup := v1.Group("signup")
		{
			// Attached rate limiter here
			signup.POST("register", authLimiter.Middleware(), userController.SignUp)
		}

		signin := v1.Group("signin")
		{
			// Attached rate limiter here
			signin.POST("", authLimiter.Middleware(), userController.SignIn)
		}

		picture := v1.Group("pictures")
		{
			picture.Use(middlewares.JWTAuthMiddleware(db, jwtService))

			pictureController := controllers.NewPictureController(services.NewPictureService(repositories.NewPictureRepository(db)))
			picture.GET("", pictureController.GetPictures)
			picture.GET(":id", pictureController.GetPictureById)
			picture.GET("count", pictureController.CountPicture)
			picture.GET("paginated", pictureController.GetPicturesPaginated)
		}

		series := v1.Group("series")
		{
			series.Use(middlewares.JWTAuthMiddleware(db, jwtService))
			seriesController := controllers.NewSeriesController(services.NewSeriesService(repositories.NewSeriesRepository(db)))
			series.GET("", seriesController.GetAllSeries)
			series.GET(":name", seriesController.GetSeriesByName)
		}

		videos := v1.Group("videos")
		{
			videos.Use(middlewares.JWTAuthMiddleware(db, jwtService))
			videosController := controllers.NewVideoController(services.NewVideoService(repositories.NewVideoRepository(db)))
			videos.GET("", videosController.GetAllVideos)
			videos.GET(":name", videosController.GetVideoByName)
		}

		health := v1.Group("health")
		{
			health.GET("ping", func(c *gin.Context) {
				if err != nil {
					c.JSON(500, gin.H{"message": "error"})
					return
				}
				c.JSON(200, gin.H{"message": "ok"})
			})
		}
	}

	// Run database migrations
	err = db.AutoMigrate(&models.Pictures{}, &models.User{}, &models.Series{}, &models.Video{})
	if err != nil {
		panic(err)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println(router.Run(":8080"))
}
