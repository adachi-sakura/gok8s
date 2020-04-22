package handlers

import (
	"github.com/buzaiguna/gok8s/appctx"
	"github.com/buzaiguna/gok8s/handlers/algorithm"
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
	router.GET("/apiserver/nodes/:nodeName", auth.AuthHandler, appctx.Handler(apiserver.GetNode))
	router.POST("/apiserver/deployments", auth.AuthHandler, appctx.Handler(apiserver.CreateDeployment))
	//metrics
	router.GET("/metrics/nodes", auth.AuthHandler, appctx.Handler(metrics.GetNodesMetrics))
	router.GET("/metrics/nodes/:nodeName", auth.AuthHandler, appctx.Handler(metrics.GetNodeMetrics))
	router.GET("/metrics/pods", auth.AuthHandler, appctx.Handler(metrics.GetPodsMetrics))
	//algorithm
	router.GET("/algorithm/parameters", auth.AuthHandler, appctx.Handler(algorithm.GetParameters))
	router.POST("/algorithm/deployments", auth.AuthHandler, appctx.Handler(algorithm.CreateDeployments))
}
