package repositories

import (
	"github.com/dmolina79/golang-github-api/src/api/utils/errors"
	"strings"
)

type CreateRepoRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (r *CreateRepoRequest) Validate() errors.ApiError {
	r.Name = strings.TrimSpace(r.Name)
	if r.Name == "" {
		return errors.NewBadRequestError("Invalid repository name")
	}

	return nil
}

type CreateRepoResponse struct {
	Id    int64  `json:"id"`
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

type CreateReposResponse struct {
	StatusCode int                 `json:"status"`
	Result     []CreateReposResult `json:"result"`
}

type CreateReposResult struct {
	Response *CreateRepoResponse `json:"repo"`
	Error    errors.ApiError     `json:"error"`
}
