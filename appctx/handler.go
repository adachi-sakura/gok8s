package appctx

import (
	"context"
	"github.com/gin-gonic/gin"
)

type HandlerFunc func(context.Context) error
type GinHandlerFunc func(*gin.Context) error

func Handler(handler HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := GetContextFromGin(c)
		err := handler(ctx)
		if err != nil {
			c.Error(err)
			c.Abort()
		}
	}
}

func GinHandler(handler GinHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := handler(c)
		if err != nil {
			c.Error(err)
			c.Abort()
		}
	}
}