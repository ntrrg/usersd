// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the MIT license.

package http

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
)

// Config wraps all the customizable options from Server.
type Config struct {
	// TCP address to listen on. If a file path is given, the server will use a
	// Unix Domain Socket.
	Addr string

	// Requests handler. Just as http.Server, if nil is given,
	// http.DefaultServeMux will be used.
	Handler http.Handler

	ShutdownCtx func() context.Context
}

// Server is a http.Server with some extra functionalities.
type Server struct {
	http.Server

	// Shutdown context used for gracefully shutdown, it is implemented as a
	// function since deadlines will start at server creation and not at shutdown.
	ShutdownCtx func() context.Context

	// Gracefully shutdown done notifier.
	Done chan struct{}
}

// NewServer creates and setups a new Server.
func NewServer(c *Config) *Server {
	s := new(Server)
	s.Done = make(chan struct{})
	s.Setup(c)
	return s
}

// ListenAndServe listen in a TCP address for HTTP requests.
func (s *Server) ListenAndServe() error {
	go s.gracefullyShutdown()
	return s.Server.ListenAndServe()
}

// ListenAndServeTLS listen in a TCP address for HTTPS/H2 requests.
func (s *Server) ListenAndServeTLS(cert, key string) error {
	go s.gracefullyShutdown()
	return s.Server.ListenAndServeTLS(cert, key)
}

// ListenAndServeUDS listen in a Unix Domain Socket for HTTPS/H2 requests.
func (s *Server) ListenAndServeUDS() error {
	uds, err := net.Listen("unix", s.Addr)

	if err != nil {
		return err
	}

	go s.gracefullyShutdown()
	return s.Server.Serve(uds)
}

// Setup prepares the Server with the given Config.
func (s *Server) Setup(c *Config) {
	if c.Addr != "" {
		s.Addr = c.Addr
	}

	if c.Handler != nil {
		s.Handler = c.Handler
	}

	if c.ShutdownCtx != nil {
		s.ShutdownCtx = c.ShutdownCtx
	} else if s.ShutdownCtx == nil {
		s.ShutdownCtx = func() context.Context { return context.Background() }
	}
}

// gracefullyShutdown starts a gracefully shutdown when a SIGTERM signal is
// launched. This must be called with the go keyword.
func (s *Server) gracefullyShutdown() {
	defer close(s.Done)

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint

	if err := s.Shutdown(s.ShutdownCtx()); err != nil {
		panic(err)
	}
}
