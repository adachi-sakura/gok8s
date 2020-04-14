package basic

import (
	"context"
	"github.com/buzaiguna/gok8s/appctx"
	"github.com/gin-gonic/gin"
)

func Context() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		ctx = appctx.WithGinContext(ctx, c)
		appctx.SetContext(c, ctx)
		//defer c.Set("ServiceContext", nil)
		c.Next()
	}
}
