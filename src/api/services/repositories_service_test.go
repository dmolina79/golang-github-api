package services

import (
	"github.com/dmolina79/golang-github-api/src/api/client/restclient"
	"github.com/dmolina79/golang-github-api/src/api/domain/repositories"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
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
		Name:"github-repo",
	}

	// execute
	res, err := RepositoryService.CreateRepo(req)

	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusUnauthorized, err.Status())
	assert.EqualValues(t, "Requires authentication", err.Message())

}

func TestReposService_CreateRepo_GoGood(t *testing.T) {
	restclient.FlushMockups()

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
		Name:"github-repo",
	}

	// execute
	res, err := RepositoryService.CreateRepo(req)

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.EqualValues(t, 123, res.Id)
	assert.EqualValues(t, "github-repo", res.Name)
	assert.EqualValues(t, "dmolina79", res.Owner)
}
