package handlers

import "encoding/json"

func NewErrResponse(err error) []byte {
	bytes, _ := json.Marshal(Response{
		Err:     err.Error(),
		Success: false,
	})
	return bytes
}

func NewSuccessResponse(msg string, data any) []byte {
	bytes, _ := json.Marshal(Response{
		Msg:     msg,
		Success: true,
		Data:    data,
	})
	return bytes
}

type Response struct {
	Success bool   `json:"success"`
	Err     string `json:"error,omitempty"`
	Msg     string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}
