package main

import (
	"net/http"
	"time"
)

func (app *application) serve() error {

	srv := http.Server{
		Addr:         "localhost:8080",
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		Handler:      app.routes(),
	}
	return srv.ListenAndServe()
}
