package errors

import (
	"encoding/json"
	"errors"
	"net/http"
)

type ApiError interface {
	Status() int
	Message() string
	Error() string
}

type apiError struct {
	ErrStatus  int    `json:"status"`
	ErrMessage string `json:"message"`
	ErrError   string `json:"error,omitempty"`
}

func (a *apiError) Error() string {
	return a.ErrError
}

func (a *apiError) Message() string {
	return a.ErrMessage
}

func (a *apiError) Status() int {
	return a.ErrStatus
}

func NewApiError(statusCode int, message string) ApiError {
	return &apiError{ErrStatus: statusCode, ErrMessage: message}
}

func NewApiErrFromBody(body []byte) (ApiError, error) {
	var res apiError
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, errors.New("invalid Json for creating an API error")
	}

	return &res, nil
}

func NewNotFoundError(m string) ApiError {
	return &apiError{
		ErrStatus:  http.StatusNotFound,
		ErrMessage: m,
	}
}

func NewInternalServerError(m string) ApiError {
	return &apiError{
		ErrStatus:  http.StatusInternalServerError,
		ErrMessage: m,
	}
}

func NewBadRequestError(m string) ApiError {
	return &apiError{
		ErrStatus:  http.StatusBadRequest,
		ErrMessage: m,
	}
}
