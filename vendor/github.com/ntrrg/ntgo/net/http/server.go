// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the MIT license.

package http

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
)

// Config wraps all the customizable options from Server.
type Config struct {
	// TCP address to listen on. If a path to a file is given, the server will use
	// a Unix Domain Socket.
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
func NewServer(c Config) *Server {
	s := new(Server)
	s.Done = make(chan struct{})
	s.Setup(c)
	return s
}

// ListenAndServe starts listening in a TCP address or in a Unix Domain Socket.
func (s *Server) ListenAndServe() error {
	addr := s.Addr

	// Gracefully shutdown
	go func() {
		defer close(s.Done)

		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		log.Println("[INFO][SERVER] Shutting down the server..")

		if err := s.Shutdown(s.ShutdownCtx()); err != nil {
			log.Fatalf("[ERROR][SERVER] Can't close the server gracefully.\n%v", err)
		} else {
			log.Println("[INFO][SERVER] All the pending tasks were done.")
		}

		log.Println("[INFO][SERVER] Server closed.")
	}()

	if strings.Contains(addr, "/") {
		uds, err := net.Listen("unix", addr)

		if err != nil {
			log.Printf("[ERROR][SERVER] Can't use the socket %s.\n%v", addr, err)
			return err
		}

		log.Printf("[INFO][SERVER] Using UDS %v..\n", addr)
		return s.Server.Serve(uds)
	}

	log.Printf("[INFO][SERVER] Listening on %v..\n", addr)
	return s.Server.ListenAndServe()
}

// Setup prepares the Server with the given Config.
func (s *Server) Setup(c Config) {
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
