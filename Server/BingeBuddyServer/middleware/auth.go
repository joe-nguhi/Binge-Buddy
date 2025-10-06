package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joe-nguhi/Binge-Buddy/Server/BingeBuddyServer/utils"
)

const (
	Bearer = len("Bearer ")
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {

		header := c.GetHeader("Authorization")

		if header == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization header"})
			c.Abort()
			return
		}

		token := header[Bearer:]

		if header == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing bearer token"})
			c.Abort()
			return
		}

		userData, err := utils.ValidateToken(token)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("userID", userData.UserID)
		c.Set("role", userData.Role)
		c.Next()
	}
}
