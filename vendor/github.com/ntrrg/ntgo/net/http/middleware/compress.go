// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the MIT license.

package middleware

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"
)

type compressedResponseWriter struct {
	http.ResponseWriter
	io.WriteCloser
}

func (w compressedResponseWriter) Write(b []byte) (int, error) {
	return w.WriteCloser.Write(b)
}

// Gzip compresses the response body. The compression level is given as an
// integer value following the compress/falte package values.
func Gzip(level int) Adapter {
	return func(h http.Handler) http.Handler {
		nh := func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				gz, err := gzip.NewWriterLevel(w, level)

				if err != nil {
					log.Printf("[ERROR] Bad compression level: %v\n%v", level, err)
				}

				gzw := compressedResponseWriter{w, gz}

				defer func() {
					if err := gzw.Close(); err != nil {
						log.Printf("[ERROR] Can't close the GZIP writer.\n%v", err)
					}
				}()

				w.Header().Set("Content-Encoding", "gzip")
				h.ServeHTTP(gzw, r)
			} else {
				h.ServeHTTP(w, r)
			}
		}

		return http.HandlerFunc(nh)
	}
}
