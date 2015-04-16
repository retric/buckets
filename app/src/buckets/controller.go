package buckets

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	//"strconv"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
)

type Action func(w http.ResponseWriter, r *http.Request)

type AppController struct{}

func (c *AppController) Action(a Action) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a(w, r)
	})
}

type MyController struct {
	AppController
	Page      map[string]string
	templates map[string]*template.Template
	includes  []string
	session   *mgo.Session
}

/* Setup template handling and session */
func (c *MyController) Init() {
	c.initTemplates()
	go c.startWatcher()
	c.initSession(DbSetup())
}

func (c *MyController) initSession(session *mgo.Session) {
	c.session = session
}

/* View Handlers */
func (c *MyController) HomeHandler(w http.ResponseWriter, req *http.Request) {
	c.renderTemplate(w, "home.tmpl", c.Page)
}

func (c *MyController) LoginHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		// implement login logic here
		http.Redirect(w, req, "/", 200)
	}
	c.renderTemplate(w, "login.tmpl", c.Page)
}

func (c *MyController) LogoutHandler(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, "/", 200)
}

func (c *MyController) StaticHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("retrieving %s", req.URL.Path)
	http.ServeFile(w, req, ".."+req.URL.Path)
}

/* Serialize JSON response and return it */
func sendJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return
}

/* API requests */
func (c *MyController) BucketsHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		buckets := getBuckets(c.session)
		sendJSON(w, buckets)
		return
	case "POST":
		var bForm BucketPart
		err := json.NewDecoder(req.Body).Decode(&bForm)
		if err != nil {
			log.Fatal("jsonError:", err)
		}
		bucket := createBucket(c.session, bForm)
		w.WriteHeader(201)
		sendJSON(w, bucket)
		return
	}
}

func (c *MyController) BucketHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, _ := vars["id"]

	switch req.Method {
	case "GET", "":
		bucket, _ := getBucket(c.session, id)
		sendJSON(w, bucket)
		return
	case "POST":
		return
	case "PUT":
		var bForm BucketPart
		err := json.NewDecoder(req.Body).Decode(&bForm)
		if err != nil {
			log.Fatal("jsonError:", err)
		}
		updateBucket(c.session, id, bForm)
		w.WriteHeader(204)
		sendJSON(w, nil)
		return
	case "DELETE":
		removeBucket(c.session, id)
		w.WriteHeader(204)
		sendJSON(w, nil)
		return
	}
}

func (c *MyController) TasksHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		buckets := getBuckets(c.session)
		sendJSON(w, buckets)
		return
	case "POST":
		var tForm TaskPart
		err := json.NewDecoder(req.Body).Decode(&tForm)
		if err != nil {
			log.Fatal("jsonError:", err)
		}
		task := createTask(c.session, tForm)
		w.WriteHeader(201)
		sendJSON(w, task)
	}
}

func (c *MyController) TaskHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, _ := vars["id"]

	switch req.Method {
	case "GET":
		task, _ := getTask(c.session, id)
		sendJSON(w, task)
		return
	case "POST":
		return
	case "PUT":
		var tForm TaskPart
		err := json.NewDecoder(req.Body).Decode(&tForm)
		if err != nil {
			log.Fatal("jsonError:", err)
		}
		updateTask(c.session, id, tForm)
		w.WriteHeader(204)
		sendJSON(w, nil)
		return
	case "DELETE":
		removeTask(c.session, id)
		w.WriteHeader(204)
		sendJSON(w, nil)
		return
	}

}
