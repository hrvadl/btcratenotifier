package handlers

import "encoding/json"

func NewErrResponse(err error) []byte {
	bytes, _ := json.Marshal(ErrorResponse{
		Err:     err.Error(),
		Success: false,
	})
	return bytes
}

func NewEmptyResponse(msg string) []byte {
	bytes, _ := json.Marshal(EmptyResponse{
		Msg:     msg,
		Success: true,
	})
	return bytes
}

func NewSuccessResponse[T any](msg string, data T) []byte {
	bytes, _ := json.Marshal(Response[T]{
		EmptyResponse: EmptyResponse{
			Msg:     msg,
			Success: true,
		},
		Data: data,
	})
	return bytes
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Err     string `json:"error"`
}

type EmptyResponse struct {
	Success bool   `json:"success"`
	Msg     string `json:"message,omitempty"`
}

type Response[T any] struct {
	EmptyResponse
	Data T `json:"data,omitempty"`
}
