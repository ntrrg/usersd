// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package main

import (
	"fmt"
)

func main() {
	fmt.Println("hello, world!")
}

/*
import (
	"flag"
	"log"
	"net/http"
	"os"
	"strings"

	nthttp "nt.web.ve/go/ntgo/net/http"

	"nt.web.ve/go/usersd/api/rest"
	"nt.web.ve/go/usersd/pkg/usersd"
)

func main() {
	var (
		usersdOpts = usersd.DefaultOptions

		restOpts  nthttp.Config
		key, cert string

		verbose bool
		debug   bool
		logfile string
	)

	flag.StringVar(
		&restOpts.Addr,
		"addr",
		":4000",
		"TCP address to listen on. If a path to a file is given, the server will "+
			"use a Unix Domain Socket.",
	)

	flag.StringVar(&key, "key", "", "TLS private key file")
	flag.StringVar(&cert, "cert", "", "TLS certificate file")
	flag.StringVar(&usersdOpts.Database, "db", "", "Database location")
	flag.BoolVar(&verbose, "verbose", true, "Enable verbosing")
	flag.BoolVar(&debug, "debug", false, "Enable debugging")
	flag.StringVar(&logfile, "log", "", "Log file location (default: stderr)")
	flag.Parse()

	if logfile != "" {
		lf, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)

		if err != nil {
			log.Fatalf("[FATAL][SERVER] Can't create/open the log file -> %v", err)
		}

		defer func() {
			if err2 := lf.Close(); err2 != nil {
				log.Fatalf("[FATAL][SERVER] Can't close the log file -> %v", err2)
			}
		}()

		log.SetOutput(lf)
	}

	if err := usersd.Init(usersdOpts); err != nil {
		log.Fatalf("[FATAL][USERSD] Can't initialize the API -> %v", err)
	}

	restOpts.Handler = rest.Mux()
	server := nthttp.NewServer(&restOpts)

	var err error

	switch {
	case strings.Contains(restOpts.Addr, "/"):
		err = server.ListenAndServeUDS()
	case key != "" && cert != "":
		err = server.ListenAndServeTLS(cert, key)
	default:
		err = server.ListenAndServe()
	}

	if err != http.ErrServerClosed {
		log.Fatalf("[FATAL][SERVER] Can't start the server -> %v", err)
	}

	<-server.Done
}
*/
