package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NoRouteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": "The processing function of the request route was not found"})
	}
}
