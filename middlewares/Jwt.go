package middlewares

import (
	"Api-Picture/models"
	"Api-Picture/services"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	lru "github.com/hashicorp/golang-lru/v2/expirable"
	"gorm.io/gorm"
)

func JWTAuthMiddleware(db *gorm.DB, jwtService *services.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		users := lru.NewLRU[uint, models.User](1024, nil, time.Minute)
		var tokenString string

		// 1. Try extracting token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		}

		// 2. Fall back to reading the token from cookie if header was empty
		if tokenString == "" {
			tokenString, _ = c.Cookie("token")
		}

		// If neither header nor cookie contains a token, reject
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// 3. Validate token
		token, err := jwtService.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// 4. Verify expiration
		expClaim, ok := claims["exp"].(float64)
		if !ok || time.Now().After(time.Unix(int64(expClaim), 0)) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// 5. Extract user_id and load user from DB
		userID, ok := claims["user_id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		user, cached := users.Get(uint(userID))
		if !cached {
			if err := db.First(&user, uint(userID)).Error; err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				c.Abort()
				return
			}
			users.Add(uint(userID), user)
		}

		c.Set("user", user)
		c.Next()
	}
}
