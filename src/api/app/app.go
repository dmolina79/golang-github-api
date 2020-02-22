package app

import (
	"github.com/dmolina79/golang-github-api/src/api/log"
	"github.com/gin-gonic/gin"
)

var (
	router *gin.Engine
)

func init() {
	router = gin.Default()
}

func StartApp() {
	log.Info("setting up routes...")
	setupRoutes()
	log.Info("routes setup completed")
	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
