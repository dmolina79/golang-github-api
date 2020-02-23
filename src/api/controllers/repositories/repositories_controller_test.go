package repositories

import (
	"encoding/json"
	"github.com/dmolina79/golang-github-api/src/api/domain/repositories"
	"github.com/dmolina79/golang-github-api/src/api/services"
	"github.com/dmolina79/golang-github-api/src/api/utils/errors"
	"github.com/dmolina79/golang-github-api/src/api/utils/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type repoServiceMock struct {
	mock.Mock
}

// stubs for mock
func (r repoServiceMock) CreateRepo(request repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.ApiError) {
	args := r.Called(request)
	if args.Error(1) != nil {
		return nil, args.Error(1).(errors.ApiError)
	}

	return args.Get(0).(*repositories.CreateRepoResponse), nil

}

func (r repoServiceMock) CreateRepos(request []repositories.CreateRepoRequest) (repositories.CreateReposResponse, errors.ApiError) {
	return repositories.CreateReposResponse{}, nil
}

func TestCreateRepo_Success(t *testing.T) {
	mockService := new(repoServiceMock)
	mockService.On("CreateRepo", mock.Anything).Return(
		&repositories.CreateRepoResponse{
			Id:    321,
			Owner: "vbuterin",
			Name:  "github-repo-test-mock",
		},
		nil)

	services.RepositoryService = mockService

	request, _ := http.NewRequest(http.MethodPost, "/repo", strings.NewReader(`{ "name": "github-repo"}`))
	response := httptest.NewRecorder()
	c := test_utils.GetMockContext(request, response)

	CreateRepo(c)

	assert.EqualValues(t, http.StatusCreated, response.Code)
	var result repositories.CreateRepoResponse
	err := json.Unmarshal(response.Body.Bytes(), &result)
	assert.Nil(t, err)
	assert.EqualValues(t, 321, result.Id)
	assert.EqualValues(t, "vbuterin", result.Owner)
	assert.EqualValues(t, "github-repo-test-mock", result.Name)
}

func TestCreateRepo_HandleError(t *testing.T) {
	mockService := new(repoServiceMock)
	mockService.On("CreateRepo", mock.Anything).Return(
		nil, errors.NewApiError(http.StatusBadRequest, "error on request"))

	services.RepositoryService = mockService

	request, _ := http.NewRequest(http.MethodPost, "/repo", strings.NewReader(`{ "name": "github-repo"}`))
	response := httptest.NewRecorder()
	c := test_utils.GetMockContext(request, response)

	CreateRepo(c)

	assert.EqualValues(t, http.StatusBadRequest, response.Code)

	apiErr, err := errors.NewApiErrFromBody(response.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusBadRequest, apiErr.Status())
	assert.EqualValues(t, "error on request", apiErr.Message())
}
