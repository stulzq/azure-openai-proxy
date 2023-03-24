package apis

import "github.com/gin-gonic/gin"

func SendError(c *gin.Context, err error) {
	c.JSON(500, ApiResponse{
		Error: ErrorDescription{
			Code:    "500",
			Message: err.Error(),
		},
	})
}
