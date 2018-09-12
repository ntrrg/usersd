// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	nthttp "github.com/ntrrg/ntgo/net/http"

	"github.com/ntrrg/usersd/internal/rest"
	"github.com/ntrrg/usersd/pkg/usersd"
)

var lf *os.File

func main() {
	defer lf.Close()
	defer usersd.Close()

	err := rest.Server.ListenAndServe()

	if err != http.ErrServerClosed {
		log.Fatalf("[ERROR][SERVER] Can't start the server. (%v)\n", err)
	}

	<-rest.Server.Done
}

func init() {
	var (
		cfg   nthttp.Config
		debug bool

		dir, logfile string
	)

	flag.StringVar(
		&cfg.Addr,
		"addr",
		":4000",
		"TCP address to listen on. If a path to a file is given, the server will "+
			"use a Unix Domain Socket.",
	)

	flag.StringVar(&dir, "db", "", "Database location")
	flag.BoolVar(&debug, "debug", false, "Enable debugging")
	flag.StringVar(&logfile, "log", "", "Log file location (default: stderr)")
	flag.Parse()

	if logfile != "" {
		var err error
		lf, err = os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil {
			log.Fatalf("[ERROR][SERVER] Can't create/open the log file. (%v)\n", err)
		}

		log.SetOutput(lf)
	}

	rest.Server.Setup(cfg)

	if err := usersd.Init(dir); err != nil {
		log.Fatalf("[ERROR][API] Can't initialize the API. (%v)\n", err)
	}
}
