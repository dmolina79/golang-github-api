package repositories

import (
	"encoding/json"
	"github.com/dmolina79/golang-github-api/src/api/client/restclient"
	"github.com/dmolina79/golang-github-api/src/api/domain/repositories"
	"github.com/dmolina79/golang-github-api/src/api/utils/errors"
	"github.com/dmolina79/golang-github-api/src/api/utils/test_utils"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	restclient.StartMockups()
	os.Exit(m.Run())
}

func TestCreateRepoInvalidJsonRequest(t *testing.T) {
	request, _ := http.NewRequest(http.MethodPost, "/repo", strings.NewReader(``))
	response := httptest.NewRecorder()
	c := test_utils.GetMockContext(request, response)

	CreateRepo(c)

	assert.EqualValues(t, http.StatusBadRequest, response.Code)

	apiErr, err := errors.NewApiErrFromBody(response.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusBadRequest, apiErr.Status())
	assert.EqualValues(t, "invalid json body", apiErr.Message())
}

func TestCreateRepo_ErrorInGH(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockUp(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(strings.NewReader(`{"message":"Requires authentication","documentation_url":"https://developer.github.com/docs"}`)),
		},
	})
	request, _ := http.NewRequest(http.MethodPost, "/repo", strings.NewReader(`{ "name": "github-repo"}`))
	response := httptest.NewRecorder()
	c := test_utils.GetMockContext(request, response)

	CreateRepo(c)

	assert.EqualValues(t, http.StatusUnauthorized, response.Code)

	apiErr, err := errors.NewApiErrFromBody(response.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusUnauthorized, apiErr.Status())
	assert.EqualValues(t, "Requires authentication", apiErr.Message())
}

func TestCreateRepo_GoGood(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockUp(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(strings.NewReader(`{ "id": 123 }`)),
		},
	})
	request, _ := http.NewRequest(http.MethodPost, "/repo", strings.NewReader(`{ "name": "github-repo"}`))
	response := httptest.NewRecorder()
	c := test_utils.GetMockContext(request, response)

	CreateRepo(c)

	assert.EqualValues(t, http.StatusCreated, response.Code)
	var result repositories.CreateRepoResponse
	err := json.Unmarshal(response.Body.Bytes(), &result)
	assert.Nil(t, err)
	assert.EqualValues(t, 123, result.Id)

}
