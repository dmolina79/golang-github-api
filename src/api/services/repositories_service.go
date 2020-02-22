package services

import (
	"fmt"
	"github.com/dmolina79/golang-github-api/src/api/config"
	"github.com/dmolina79/golang-github-api/src/api/domain/github"
	"github.com/dmolina79/golang-github-api/src/api/domain/repositories"
	"github.com/dmolina79/golang-github-api/src/api/log"
	"github.com/dmolina79/golang-github-api/src/api/providers/github_provider"
	"github.com/dmolina79/golang-github-api/src/api/utils/errors"
	"net/http"
	"sync"
)

type reposService struct{}

type repoServiceInterface interface {
	CreateRepo(request repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.ApiError)
	CreateRepos(request []repositories.CreateRepoRequest) (repositories.CreateReposResponse, errors.ApiError)
}

var (
	RepositoryService repoServiceInterface
)

func init() {
	RepositoryService = &reposService{}
}

func (s *reposService) CreateRepo(input repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.ApiError) {
	clientId := "1"
	if err := input.Validate(); err != nil {
		return nil, err
	}

	request := github.CreateRepoRequest{
		Name:        input.Name,
		Description: input.Description,
		Private:     false,
	}

	log.Info("sending request to external api", fmt.Sprintf("client_id:%s", clientId), "status:pending")
	res, err := github_provider.CreateRepo(config.GetGithubAccessToken(), request)

	if err != nil {
		log.Error("sending request to external api", err, fmt.Sprintf("client_id:%s", clientId), "status:error")
		return nil, errors.NewApiError(err.StatusCode, err.Message)
	}

	log.Info("response obtained from external api", fmt.Sprintf("client_id:%s", clientId), "status:success")
	result := repositories.CreateRepoResponse{
		Id:    res.Id,
		Owner: res.Owner.Login,
		Name:  res.Name,
	}

	return &result, nil
}

func (s *reposService) CreateRepos(req []repositories.CreateRepoRequest) (repositories.CreateReposResponse, errors.ApiError) {
	input := make(chan repositories.CreateReposResult)
	output := make(chan repositories.CreateReposResponse)
	defer close(output)

	var wg sync.WaitGroup
	go s.handleRepoResults(&wg, input, output)

	for _, current := range req {
		wg.Add(1)
		go s.createRepoConcurrent(current, input)
	}

	// wait until all routines are done
	wg.Wait()
	close(input)

	result := <-output

	successCreations := 0
	for _, current := range result.Results {
		if current.Response != nil {
			successCreations++
		}
	}

	switch successCreations {
	case len(req):
		result.StatusCode = http.StatusCreated
	case 0:
		result.StatusCode = result.Results[0].Error.Status()
	default:
		result.StatusCode = http.StatusPartialContent
	}
	fmt.Println(fmt.Sprintf("successCreations %v", successCreations))

	return result, nil
}

func (s *reposService) handleRepoResults(wg *sync.WaitGroup, input chan repositories.CreateReposResult, out chan repositories.CreateReposResponse) {
	var results repositories.CreateReposResponse

	for incomingRes := range input {
		repoResult := repositories.CreateReposResult{
			Response: incomingRes.Response,
			Error:    incomingRes.Error,
		}
		results.Results = append(results.Results, repoResult)

		wg.Done()
	}

	out <- results
}

func (s *reposService) createRepoConcurrent(input repositories.CreateRepoRequest, out chan repositories.CreateReposResult) {
	if err := input.Validate(); err != nil {
		out <- repositories.CreateReposResult{Error: err}
		return
	}

	res, err := s.CreateRepo(input)

	if err != nil {
		out <- repositories.CreateReposResult{Error: err}
		return
	}

	out <- repositories.CreateReposResult{Response: res}

}
