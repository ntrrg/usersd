// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	nthttp "github.com/ntrrg/ntgo/net/http"

	"github.com/ntrrg/usersd/api/rest"
	"github.com/ntrrg/usersd/pkg/usersd"
)

func main() {
	var (
		usersdOpts = usersd.DefaultOptions

		restOpts nthttp.Config

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

	flag.StringVar(&usersdOpts.Admin, "admin", "admin", "Administrator user")
	flag.StringVar(&usersdOpts.Database, "db", "", "Database location")
	flag.BoolVar(&verbose, "verbose", true, "Enable debugging")
	flag.BoolVar(&debug, "debug", false, "Enable debugging")
	flag.StringVar(&logfile, "log", "", "Log file location (default: stderr)")
	flag.Parse()

	usersdOpts.Verbose = verbose
	usersdOpts.Debug = debug

	if logfile != "" {
		lf, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil {
			log.Fatalf("[ERROR][SERVER] Can't create/open the log file. (%v)\n", err)
		}

		defer lf.Close()
		log.SetOutput(lf)
		usersdOpts.Logger = log.New(lf, "", log.LstdFlags)
	}

	if err := usersd.Init(usersdOpts); err != nil {
		log.Fatalf("[ERROR][API] Can't initialize the API. (%v)\n", err)
	}

	defer usersd.Close()
	rest.Server.Setup(restOpts)

	if err := rest.Server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("[ERROR][SERVER] Can't start the server. (%v)\n", err)
	}

	<-rest.Server.Done
}
