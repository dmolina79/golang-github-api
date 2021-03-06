package repositories

import (
	"github.com/dmolina79/golang-github-api/src/api/domain/repositories"
	"github.com/dmolina79/golang-github-api/src/api/services"
	"github.com/dmolina79/golang-github-api/src/api/utils/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateRepo(c *gin.Context) {
	var request repositories.CreateRepoRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		apiErr := errors.NewBadRequestError("invalid json body")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	res, err := services.RepositoryService.CreateRepo(request)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusCreated, res)
}

func CreateRepos(c *gin.Context) {
	var requests []repositories.CreateRepoRequest
	if err := c.ShouldBindJSON(&requests); err != nil {
		apiErr := errors.NewBadRequestError("invalid json body")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	res, err := services.RepositoryService.CreateRepos(requests)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(res.StatusCode, res)
}
