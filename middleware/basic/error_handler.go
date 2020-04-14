package basic

import (
	"github.com/buzaiguna/gok8s/apperror"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) != 0 {
			errs := c.Errors
			err := errs[len(errs)-1].Err
			appErr := apperror.Wrap(err)
			if !c.Writer.Written() {
				c.JSON(appErr.StatusCode(), appErr)
			}
		}
	}
}
