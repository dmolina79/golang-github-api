package services

import (
	"github.com/dmolina79/golang-github-api/src/api/config"
	"github.com/dmolina79/golang-github-api/src/api/domain/github"
	"github.com/dmolina79/golang-github-api/src/api/domain/repositories"
	"github.com/dmolina79/golang-github-api/src/api/providers/github_provider"
	"github.com/dmolina79/golang-github-api/src/api/utils/errors"
	"strings"
)

type reposService struct{}

type repoServiceInterface interface {
	CreateRepo(request repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.ApiError)
}

var (
	RepositoryService repoServiceInterface
)

func init() {
	RepositoryService = &reposService{}
}

func (s *reposService) CreateRepo(input repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.ApiError) {
	input.Name = strings.TrimSpace(input.Name)
	if input.Name == "" {
		return nil, errors.NewBadRequestError("Invalid repository name")
	}

	request := github.CreateRepoRequest{
		Name:        input.Name,
		Description: input.Description,
		Private:     false,
	}

	res, err := github_provider.CreateRepo(config.GetGithubAccessToken(), request)

	if err != nil {
		return nil, errors.NewApiError(err.StatusCode, err.Message)
	}

	result := repositories.CreateRepoResponse{
		Id:    res.Id,
		Owner: res.Owner.Login,
		Name:  res.Name,
	}

	return &result, nil
}
