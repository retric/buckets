package main

import "net/http"

type Action func(w http.ResponseWriter, r *http.Request)

type AppController struct{}

func (c *AppController) Action(a Action) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a(w, r)
	})
}
