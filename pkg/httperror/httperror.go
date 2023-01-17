package httperror

import "github.com/gin-gonic/gin"

type Error struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}

func NewError(c *gin.Context, status int, err error) {
	er := Error{
		Code:    status,
		Message: err.Error(),
	}
	c.JSON(status, er)
}
