package github_provider

import (
	"errors"
	"github.com/dmolina79/golang-github-api/src/api/client/restclient"
	"github.com/dmolina79/golang-github-api/src/api/domain/github"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestMain(t *testing.M) {
	restclient.StartMockups()
	os.Exit(t.Run())
}

func TestConstants(t *testing.T) {
	assert.EqualValues(t, "Authorization", headerAuthorization)
	assert.EqualValues(t, "token %s", headerAuthorizationFormat)
	assert.EqualValues(t, "https://api.github.com/user/repos", urlCreateRepo)
}

func Test_getAuthorizationHeader(t *testing.T) {
	header := getAuthorizationHeader("abc123")
	assert.EqualValues(t, "token abc123", header)
}

func TestCreateRepoErrorRestclient(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockUp(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Err:        errors.New("Invalid rest client response"),
	})

	response, err := CreateRepo("", github.CreateRepoRequest{})

	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, "Invalid rest client response", err.Message)
}

func TestCreateRepoErrorUnauthorized(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockUp(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body: ioutil.NopCloser(strings.NewReader(`{"message":"Requires authentication","documentation_url":"https://developer.github.com/v3/repos/#create"}`)),
		},
	})

	response, err := CreateRepo("", github.CreateRepoRequest{})

	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusUnauthorized, err.StatusCode)
	assert.EqualValues(t, "Requires authentication", err.Message)
}

func TestCreateRepoSuccess(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockUp(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body: ioutil.NopCloser(strings.NewReader(`{"id": 123, "name": "golang-repo", "full_name": "dmolina79/golang-github-api"}`)),
		},
	})

	r, err := CreateRepo("", github.CreateRepoRequest{})

	assert.Nil(t, err)
	assert.NotNil(t, r)
	assert.EqualValues(t, 123, r.Id)
	assert.EqualValues(t, "golang-repo", r.Name)
	assert.EqualValues(t, "dmolina79/golang-github-api", r.FullName)
}
