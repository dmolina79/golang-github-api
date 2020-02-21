package github_provider

import (
	"encoding/json"
	"fmt"
	"github.com/dmolina79/golang-github-api/src/api/client/restclient"
	"github.com/dmolina79/golang-github-api/src/api/domain/github"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	headerAuthorization       = "Authorization"
	headerAuthorizationFormat = "token %s"
	urlCreateRepo             = "https://api.github.com/user/repos"
)

func getAuthorizationHeader(accessToken string) string {
	return fmt.Sprintf(headerAuthorizationFormat, accessToken)
}

func CreateRepo(accessToken string, request github.CreateRepoRequest) (*github.CreateRepoResponse, *github.GithubErrorResponse) {
	headers := http.Header{}
	headers.Set(headerAuthorization, getAuthorizationHeader(accessToken))

	resp, err := restclient.Post(urlCreateRepo, request, headers)

	if err != nil {
		log.Println(fmt.Sprintf("Error when trying to create new repo in github: %s", err.Error()))
		return nil, &github.GithubErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}

	bytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, &github.GithubErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "invalid  response body",
		}
	}

	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		var errorResp github.GithubErrorResponse
		if err := json.Unmarshal(bytes, &errorResp); err != nil {
			return nil, &github.GithubErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "invalid  json error response body",
			}
		}
		errorResp.StatusCode = resp.StatusCode
		return nil, &errorResp
	}

	var result github.CreateRepoResponse
	if err := json.Unmarshal(bytes, &result); err != nil {
		log.Println(fmt.Sprintf("Error when trying to unmarshal create repo success response: %s", err.Error()))
		return nil, &github.GithubErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error when trying to unmarshal github create repo response",
		}
	}

	return &result, nil
}
