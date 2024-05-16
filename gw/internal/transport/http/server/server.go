package server

import "net/http"

func NewHTTP() *http.Server {
	return &http.Server{}
}
