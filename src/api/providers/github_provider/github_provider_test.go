package github_provider

import (
	"errors"
	"github.com/dmolina79/golang-github-api/src/api/client/restclient"
	"github.com/dmolina79/golang-github-api/src/api/domain/github"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

func TestMain(t *testing.M) {
	restclient.StartMockups()
	os.Exit(t.Run())
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
	restclient.StopMockups()
}
