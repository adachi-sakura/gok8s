package handlers

import (
	"github.com/buzaiguna/gok8s/appctx"
	"github.com/buzaiguna/gok8s/handlers/auth"
	"github.com/buzaiguna/gok8s/handlers/apiserver"
	"github.com/buzaiguna/gok8s/handlers/metrics"
	"github.com/gin-gonic/gin"
	"net/http"
)

func LoadRoutes(router gin.IRouter) {
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})
	//apiserver
	router.GET("/apiserver/pods", auth.AuthHandler, appctx.Handler(apiserver.GetPods))
	//metrics
	router.GET("/metrics/nodes", auth.AuthHandler, appctx.Handler(metrics.GetNodesMetrics))
	router.GET("/metrics/:nodeName", auth.AuthHandler, )
}
