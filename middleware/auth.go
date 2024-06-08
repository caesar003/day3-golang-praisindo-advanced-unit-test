package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, password, ok := c.Request.BasicAuth()
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization basic token required"})
			c.Abort()
			return
		}

		const (
			validUser     = "superadmin"
			validPassword = "supersecretpassword"
		)
		isValid := (username == validUser) && (password == validPassword)

		if !isValid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token"})
			c.Abort()
			return
		}

		c.Next()
	}
}
