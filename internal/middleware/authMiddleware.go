package middleware

import (
	helper "contractfinder/internal/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Unauthenticated"})
			c.Abort()
			return
		}
		claims, err := helper.ValidateToken(token)
		if err != "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err})
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Set("uuid", claims.Uid)
		c.Next()
	}
}
