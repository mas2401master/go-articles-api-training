package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mas2401master/go-articles-api-training/pkg/encryption"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("authorization")
		if !encryption.ValidateToken(authorizationHeader) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "message": "user not authorized"})
			return
		} else {
			c.Next()
		}
	}
}
