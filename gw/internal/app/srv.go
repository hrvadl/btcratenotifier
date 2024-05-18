package app

import (
	"log"
	"net/http"
)

const (
	idleTimeout        = 240 // most optimal time in seconds
	writeHeaderTimeout = 15  // recommended time in seconds for timeout
	readHeaderTimeout  = 30  // sec. most optimal time in seconds for waiting for header
)

func newServer(h http.Handler, addr string, log *log.Logger) *http.Server {
	srv := &http.Server{
		Handler:      h,
		ErrorLog:     log,
		Addr:         addr,
		IdleTimeout:  idleTimeout,
		WriteTimeout: writeHeaderTimeout,
		ReadTimeout:  readHeaderTimeout,
	}
	return srv
}
