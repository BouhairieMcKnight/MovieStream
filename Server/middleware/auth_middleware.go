package middleware

import (
	"net/http"
	"time"

	"github.com/BouhairieMcKnight/MovieStream/Server/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := utils.GetCookieToken("access_token", c)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error})
			c.Abort()
			return
		}

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			c.Abort()
			return
		}
		
		claims, err := utils.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims.ExpiresAt.Time.Before(time.Now()) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}



		c.Set("userId", claims.UserId)
		c.Set("role", claims.Role)

		c.Next()
	}
}