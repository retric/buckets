package main

import (
	"buckets"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

/* Main */
func main() {
	c := &buckets.MyController{Page: map[string]string{"Static": "static"}}
	muxer := mux.NewRouter().StrictSlash(true)
	c.Init()

	/* Routes */
	muxer.HandleFunc("/", c.HomeHandler)
	muxer.Path("/login").
		Methods("GET", "POST").
		Handler(c.Action(c.LoginHandler))
	muxer.Path("/logout").
		Handler(c.Action(c.LogoutHandler))
	muxer.PathPrefix("/static/").
		Methods("GET").
		Handler(c.Action(c.StaticHandler))

	/* API routes */
	muxer.Path("/api/buckets/").
		Handler(c.Action(c.BucketsHandler))
	muxer.Path("/api/buckets/{id}").
		Handler(c.Action(c.BucketHandler))
	muxer.Path("/api/tasks").
		Handler(c.Action(c.TasksHandler))
	muxer.Path("/api/task/{id}").
		Handler(c.Action(c.TaskHandler))

	n := negroni.Classic()
	n.UseHandler(muxer)
	n.Run(":3000")
}
