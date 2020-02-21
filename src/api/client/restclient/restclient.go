package restclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	enabledMocks = false
	mocks        = make(map[string]*Mock)
)

type Mock struct {
	Url        string
	HttpMethod string
	Response   *http.Response
	Err        error
}

func getMockId(httpMethod string, url string) string {
	return fmt.Sprintf("%s %s", httpMethod, url)
}

func StartMockups() {
	enabledMocks = true
}

func StopMockups() {
	enabledMocks = false
}

func FlushMockups() {
	mocks = make(map[string]*Mock)
}

func AddMockUp(m Mock) {
	mocks[getMockId(m.HttpMethod, m.Url)] = &m
}

func Post(url string, body interface{}, headers http.Header) (*http.Response, error) {
	if enabledMocks {
		mock := mocks[getMockId(http.MethodPost, url)]
		if mock == nil {
			return nil, errors.New("no mockup found for given url request")
		}

		return mock.Response, mock.Err
	}

	jsonBytes, err := json.Marshal(body)

	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBytes))
	request.Header = headers

	client := http.Client{}

	return client.Do(request)
}
