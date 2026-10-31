package main

import (
	"Api-Picture/controllers"
	"Api-Picture/docs"
	"Api-Picture/middlewares"
	"Api-Picture/models"
	"Api-Picture/repositories"
	"Api-Picture/services"
	"Api-Picture/streaming"
	"log"
	"net/http"
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

	// Video files are served from STORAGE_ROOT (default ./storage, which is
	// /app/storage in the container). Set VIDEO_ACCEL_PREFIX (e.g.
	// /protected-videos/) to let nginx send the bytes via X-Accel-Redirect.
	storageRoot := os.Getenv("STORAGE_ROOT")
	if storageRoot == "" {
		storageRoot = "storage"
	}
	videoStream, streamErr := streaming.NewServer(storageRoot, os.Getenv("VIDEO_ACCEL_PREFIX"))
	if streamErr != nil {
		log.Printf("video streaming disabled: %v", streamErr)
	} else {
		defer videoStream.Close()
		log.Printf("video streaming from %s (nginx accel: %v)", storageRoot, videoStream.Accelerated())
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	// Never gzip video: it is already compressed, gzip breaks byte ranges,
	// and it forces the whole stream through a compressor.
	router.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithExcludedPaths([]string{"/api/v1/stream/"})))
	router.Use(middlewares.CORSMiddleware())

	v1 := router.Group("api/v1")
	{
		jwtService := services.NewJWTService(os.Getenv("SECRET_KEY"), 90*24*time.Hour)
		// One shared instance so all groups share the user cache.
		authRequired := middlewares.JWTAuthMiddleware(db, jwtService)
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
			picture.Use(authRequired)

			pictureController := controllers.NewPictureController(services.NewPictureService(repositories.NewPictureRepository(db)))
			picture.GET("", pictureController.GetPictures)
			picture.GET(":id", pictureController.GetPictureById)
			picture.GET("count", pictureController.CountPicture)
			picture.GET("paginated", pictureController.GetPicturesPaginated)
		}

		series := v1.Group("series")
		{
			series.Use(authRequired)
			seriesController := controllers.NewSeriesController(services.NewSeriesService(repositories.NewSeriesRepository(db)))
			series.GET("", seriesController.GetAllSeries)
			series.GET(":name", seriesController.GetSeriesByName)
		}

		videosController := controllers.NewVideoController(
			services.NewVideoService(repositories.NewVideoRepository(db)), videoStream)

		videos := v1.Group("videos")
		{
			videos.Use(authRequired)
			videos.GET("", videosController.GetAllVideos)
			videos.GET(":name", videosController.GetVideoByName)
		}

		// Kept outside /videos so the numeric id route cannot collide with
		// the /videos/:name search route.
		stream := v1.Group("stream")
		{
			stream.Use(authRequired)
			stream.GET("videos/:id", videosController.StreamVideo)
			stream.HEAD("videos/:id", videosController.StreamVideo)
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
	models.EnsureSearchIndexes(db)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       120 * time.Second,
		// No WriteTimeout: it would cut off long video streams.
	}
	log.Println(srv.ListenAndServe())
}
