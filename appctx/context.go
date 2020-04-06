package appctx

import (
	"context"
	"github.com/gin-gonic/gin"
)

const (
	keyGinContext 	= "GinContext"
	keyContext		= "ServiceContext"
)

func WithGinContext(ctx context.Context, ginContext *gin.Context) context.Context {
	return context.WithValue(ctx, keyGinContext, ginContext)
}

func GinContext(ctx context.Context) *gin.Context {
	val := ctx.Value(keyGinContext)
	if val == nil {
		return nil
	}
	return val.(*gin.Context)
}

func GetContextFromGin(c *gin.Context) context.Context {
	val, _ := c.Get(keyContext)
	return val.(context.Context)
}
