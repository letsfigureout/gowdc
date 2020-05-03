package web

import (
	"net/http"
	"time"
)

//BackendServer
type BackendServer struct {
	Router          http.Handler
	Addr            string
	WriteTimeout    time.Duration
	ReadTimeout     time.Duration
	ShutdownTimeout time.Duration
	Server          http.Server
}

func NewBackendServer(router http.Handler, addr string, writeTimeout, readTimeout, shutdownTimeout int) BackendServer {
	be := BackendServer{
		Router:          router,
		Addr:            addr,
		WriteTimeout:    time.Duration(writeTimeout) * time.Second,
		ReadTimeout:     time.Duration(readTimeout) * time.Second,
		ShutdownTimeout: time.Duration(shutdownTimeout) * time.Second,
		Server: http.Server{
			Addr:         addr,
			WriteTimeout: time.Duration(writeTimeout) * time.Second,
			ReadTimeout:  time.Duration(readTimeout) * time.Second,
			Handler:      router,
		},
	}

	return be
}
