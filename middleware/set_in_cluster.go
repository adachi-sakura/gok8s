package middleware

import (
	"github.com/buzaiguna/gok8s/appctx"
	"github.com/buzaiguna/gok8s/config"
	"github.com/gin-gonic/gin"
)

func SetInClusterConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := appctx.GetContextFromGin(c)
		inClusterConfig := config.InClusterConfig
		newCtx := appctx.WithInClusterConfig(c, inClusterConfig)

		appctx.SetContext(c, newCtx)
		c.Next()
	}
}