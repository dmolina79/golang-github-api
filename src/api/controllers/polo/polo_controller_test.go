package polo

import (
	"github.com/dmolina79/golang-github-api/src/api/utils/test_utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestConstants(t *testing.T) {
	assert.EqualValues(t, "polo", polo)
}

func TestPolo(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/marco", nil)
	c := test_utils.GetMockContext(req, res)

	Marco(c)

	assert.EqualValues(t, http.StatusOK, res.Code)
	assert.EqualValues(t, "polo", res.Body.String())

}