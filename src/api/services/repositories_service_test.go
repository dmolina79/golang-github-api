package services

import (
	"github.com/dmolina79/golang-github-api/src/api/client/restclient"
	"github.com/dmolina79/golang-github-api/src/api/domain/repositories"
	"github.com/dmolina79/golang-github-api/src/api/utils/errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
)

func TestMain(m *testing.M) {
	restclient.StartMockups()
	os.Exit(m.Run())
}

func TestReposService_CreateRepo_InvalidInputName(t *testing.T) {
	req := repositories.CreateRepoRequest{}

	res, err := RepositoryService.CreateRepo(req)

	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.Status())
	assert.EqualValues(t, "Invalid repository name", err.Message())

}

func TestReposService_CreateRepo_HandleErrorFromGH(t *testing.T) {
	// setup
	restclient.FlushMockups()

	restclient.AddMockUp(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(strings.NewReader(`{"message":"Requires authentication","documentation_url":"https://developer.github.com/docs"}`)),
		},
	})

	req := repositories.CreateRepoRequest{
		Name: "github-repo",
	}

	// execute
	res, err := RepositoryService.CreateRepo(req)

	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusUnauthorized, err.Status())
	assert.EqualValues(t, "Requires authentication", err.Message())

}

func TestReposService_CreateRepo_GoGood(t *testing.T) {
	// setup
	restclient.FlushMockups()

	restclient.AddMockUp(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(strings.NewReader(`{"id": 123, "name": "github-repo", "owner": { "login": "dmolina79" } }`)),
		},
	})

	req := repositories.CreateRepoRequest{
		Name: "github-repo",
	}

	// execute
	res, err := RepositoryService.CreateRepo(req)

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.EqualValues(t, 123, res.Id)
	assert.EqualValues(t, "github-repo", res.Name)
	assert.EqualValues(t, "dmolina79", res.Owner)
}

func TestReposService_CreateRepoConcurrent_InvalidRequest(t *testing.T) {
	request := repositories.CreateRepoRequest{}
	output := make(chan repositories.CreateReposResult)
	service := reposService{}

	go service.createRepoConcurrent(request, output)

	result := <-output
	assert.NotNil(t, result)
	assert.Nil(t, result.Response)
	assert.NotNil(t, result.Error)
	assert.EqualValues(t, http.StatusBadRequest, result.Error.Status())
	assert.EqualValues(t, "Invalid repository name", result.Error.Message())
}

func TestReposService_CreateRepoConcurrent_ErrorFromGH(t *testing.T) {
	// setup
	restclient.FlushMockups()

	restclient.AddMockUp(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(strings.NewReader(`{"message":"Requires authentication","documentation_url":"https://developer.github.com/docs"}`)),
		},
	})
	request := repositories.CreateRepoRequest{Name: "my-github-repo"}
	output := make(chan repositories.CreateReposResult)
	service := reposService{}

	go service.createRepoConcurrent(request, output)

	result := <-output
	assert.NotNil(t, result)
	assert.Nil(t, result.Response)
	assert.NotNil(t, result.Error)
	assert.EqualValues(t, http.StatusUnauthorized, result.Error.Status())
	assert.EqualValues(t, "Requires authentication", result.Error.Message())
}

func TestReposService_CreateRepoConcurrent_GoGood(t *testing.T) {
	// setup
	restclient.FlushMockups()

	restclient.AddMockUp(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(strings.NewReader(`{"id": 123, "name": "my-github-repo", "owner": { "login": "dmolina79" } }`)),
		},
	})
	request := repositories.CreateRepoRequest{Name: "my-github-repo"}
	output := make(chan repositories.CreateReposResult)
	service := reposService{}

	go service.createRepoConcurrent(request, output)

	result := <-output
	assert.NotNil(t, result)
	assert.Nil(t, result.Error)
	assert.NotNil(t, result.Response)
	assert.EqualValues(t, 123, result.Response.Id)
	assert.EqualValues(t, "my-github-repo", result.Response.Name)
	assert.EqualValues(t, "dmolina79", result.Response.Owner)
}

func TestReposService_HandleRepoResults(t *testing.T) {
	input := make(chan repositories.CreateReposResult)
	output := make(chan repositories.CreateReposResponse)
	var wg sync.WaitGroup

	service := reposService{}
	go service.handleRepoResults(&wg, input, output)

	wg.Add(1)
	go func() {
		input <- repositories.CreateReposResult{
			Error: errors.NewBadRequestError("invalid repository name"),
		}
	}()

	wg.Wait()
	close(input)

	result := <-output
	assert.NotNil(t, result)
	assert.EqualValues(t, 0, result.StatusCode)
	assert.EqualValues(t, 1, len(result.Results))
	assert.NotNil(t, result.Results[0].Error)
	assert.EqualValues(t, http.StatusBadRequest, result.Results[0].Error.Status())
	assert.EqualValues(t, "invalid repository name", result.Results[0].Error.Message())
}

func TestReposService_CreateRepos_InvalidRequests(t *testing.T) {
	badRequests := []repositories.CreateRepoRequest{
		{},
		{Name: "  "},
	}

	res, err := RepositoryService.CreateRepos(badRequests)

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.EqualValues(t, 2, len(res.Results))
	assert.EqualValues(t, http.StatusBadRequest, res.StatusCode)

	assert.Nil(t, res.Results[0].Response)
	assert.EqualValues(t, http.StatusBadRequest, res.Results[0].Error.Status())
	assert.EqualValues(t, "Invalid repository name", res.Results[0].Error.Message())

	assert.Nil(t, res.Results[1].Response)
	assert.EqualValues(t, http.StatusBadRequest, res.Results[1].Error.Status())
	assert.EqualValues(t, "Invalid repository name", res.Results[1].Error.Message())
}

func TestReposService_CreateRepos_PartialSuccess(t *testing.T) {
	// setup
	restclient.FlushMockups()

	restclient.AddMockUp(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(strings.NewReader(`{"id": 123, "name": "my-github-repo", "owner": { "login": "dmolina79" } }`)),
		},
	})
	requests := []repositories.CreateRepoRequest{
		{},
		{Name: "my-github-repo"},
	}

	res, err := RepositoryService.CreateRepos(requests)

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.EqualValues(t, 2, len(res.Results))
	assert.EqualValues(t, http.StatusPartialContent, res.StatusCode)

	for _, result := range res.Results {
		if result.Error != nil {
			assert.EqualValues(t, http.StatusBadRequest, result.Error.Status())
			assert.EqualValues(t, "Invalid repository name", result.Error.Message())
			continue
		}

		assert.EqualValues(t, 123, result.Response.Id)
		assert.EqualValues(t, "my-github-repo", result.Response.Name)
		assert.EqualValues(t, "dmolina79", result.Response.Owner)
	}
}

// TODO: fix this test to refactor mocking
func TestReposService_CreateRepos_AllGood(t *testing.T) {
	// setup
	restclient.FlushMockups()

	restclient.AddMockUp(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(strings.NewReader(`{"id": 123, "name": "my-github-repo", "owner": { "login": "dmolina79" } }`)),
		},
	})
	requests := []repositories.CreateRepoRequest{
		{Name: "my-github-repo"},
		{Name: "my-github-repo"},
	}

	res, err := RepositoryService.CreateRepos(requests)

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.EqualValues(t, 2, len(res.Results))
	/*assert.EqualValues(t, http.StatusCreated, res.StatusCode)

	for _, result := range res.Results {
		assert.Nil(t, result.Error)
		assert.EqualValues(t, 123, result.Response.Id)
		assert.EqualValues(t, "my-github-repo", result.Response.Name)
		assert.EqualValues(t, "dmolina79", result.Response.Owner)
	}*/
}
