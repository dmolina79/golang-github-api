package repositories

import (
	"github.com/dmolina79/golang-github-api/src/api/utils/errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateRepoInvalidJsonRequest(t *testing.T) {
	response :=  httptest.NewRecorder()
	c, _ := gin.CreateTestContext(response)

	request, _ := http.NewRequest(http.MethodPost, "/repo", strings.NewReader(``))
	c.Request = request

	CreateRepo(c)

	assert.EqualValues(t, http.StatusBadRequest, response.Code)

	apiErr, err := errors.NewApiErrFromBody(response.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusBadRequest, apiErr.Status())
	assert.EqualValues(t, "invalid json body", apiErr.Message())
}

