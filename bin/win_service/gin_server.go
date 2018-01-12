package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.xibao100.com/skyline/skyline/cubes/rest"
)

func run_gin_server() {
	rest.SetLogger(logger)

	gin.SetMode(gin.ReleaseMode)
	gin.DisableBindValidation()
	router := gin.New()
	router.Use(gin.Recovery())
	cubesGroup := router.Group("/")
	{
		cubesGroup.POST("rpt", rest.CubesRptHandler)
		cubesGroup.GET("hello", rest.HelloHandler)
	}
	logger.Infof("Cubes Server started at:0.0.0.0:%d", *portFlag)
	defer func() {
		logger.Infof("Cubes Server exit from:0.0.0.0:%d", *portFlag)
	}()
	router.Run(fmt.Sprintf(":%d", *portFlag))
}
