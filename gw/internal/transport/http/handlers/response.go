package handlers

import "encoding/json"

// NewErrResponse constructs the JSON encoded response,
// which represents request failure from API.
// NOTE: it expects error to be not nil value.
func NewErrResponse(err error) ([]byte, error) {
	bytes, err := json.Marshal(Response[struct{}]{
		Err:     err.Error(),
		Success: false,
	})
	return bytes, err
}

// NewEmptyResponse constructs the JSON encoded response,
// which represents request success from API, but doesn't have any
// data, so it omits data field in response.
func NewEmptyResponse(msg string) ([]byte, error) {
	bytes, err := json.Marshal(Response[any]{
		Msg:     msg,
		Success: true,
	})
	return bytes, err
}

// NewSuccessResponse constructs the JSON encoded response,
// which represents request success from API.
// NOTE: if data is nil then it will be omitter in the result.
func NewSuccessResponse[T any](msg string, data T) ([]byte, error) {
	bytes, err := json.Marshal(Response[T]{
		Msg:     msg,
		Success: true,
		Data:    data,
	})
	return bytes, err
}

// Response struct is a JSON encoded response,
// which represents request success from API.
// NOTE: if data is nil then it will be omitter in the result.
type Response[T any] struct {
	Success bool   `json:"success"`
	Msg     string `json:"message,omitempty"`
	Data    T      `json:"data,omitempty"`
	Err     string `json:"error,omitempty"`
}
