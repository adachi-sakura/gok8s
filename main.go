package main

import (
	"flag"
	cfg "github.com/buzaiguna/gok8s/config"
	"github.com/buzaiguna/gok8s/handlers"
	"github.com/buzaiguna/gok8s/middleware"
	"github.com/buzaiguna/gok8s/middleware/basic"
	"github.com/gin-gonic/gin"
)


var (
	bandWidthFileLine int
)

func init() {
	flag.IntVar(&bandWidthFileLine, "line", 5, "the line where bandwidth value to be used in algorithm")
	flag.Parse()
}

func main() {
	cfg.InitBandwidth(bandWidthFileLine)
	cfg.InitInClusterConfig()
	listenAndServe()

}

func listenAndServe() {
	router := gin.Default()
	router.Use(UseMiddleWares()...)
	handlers.LoadRoutes(router)
	router.Run(":8080")
}

func UseMiddleWares() []gin.HandlerFunc {
	return []gin.HandlerFunc {
		basic.Context(),
		basic.ErrorHandler(),
		middleware.SetInClusterConfig(),
	}
}