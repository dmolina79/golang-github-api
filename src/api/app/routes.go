package app

import (
	"github.com/dmolina79/golang-github-api/src/api/controllers/polo"
	"github.com/dmolina79/golang-github-api/src/api/controllers/repositories"
)

func setupRoutes() {
	router.POST("/repo", repositories.CreateRepo)
	router.GET("/marco", polo.Marco)
}
