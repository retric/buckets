package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"gopkg.in/fsnotify.v1"
)

type MyController struct {
	AppController
	page      map[string]string
	templates map[string]*template.Template
	includes  []string
}

/* View Handlers */
func (c *MyController) HomeHandler(w http.ResponseWriter, req *http.Request) {
	c.renderTemplate(w, "home.tmpl", c.page)
}

func (c *MyController) LoginHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		// implement login logic here
		http.Redirect(w, req, "/", 200)
	}
	c.renderTemplate(w, "login.tmpl", c.page)
}

func (c *MyController) LogoutHandler(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, "/", 200)
}

func (c *MyController) StaticHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("retrieving %s", req.URL.Path)
	http.ServeFile(w, req, ".."+req.URL.Path)
}

/* Template Renderer */
func (c *MyController) renderTemplate(w http.ResponseWriter, name string, data map[string]string) {
	// Ensure the template exists in the map.
	fmt.Printf("beginning template %s\n", name)
	tmpl, ok := c.templates[name]
	if !ok {
		fmt.Printf("Template %s does not exist.", name)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Printf("rendering template %s\n", name)
	tmpl.ExecuteTemplate(w, "layout", data)
}

/* Initialize templates */
func (c *MyController) initTemplates() {
	if c.templates == nil {
		c.templates = make(map[string]*template.Template)
	}

	templatesDir := "../templates/"
	files, err := filepath.Glob(templatesDir + "*.tmpl")
	if err != nil {
		panic(fmt.Sprintf("Error: Unable to fetch templates"))
	}

	includes, err := filepath.Glob(templatesDir + "partials/*.tmpl")
	if err != nil {
		panic(fmt.Sprintf("Error: Unable to fetch templates"))
	}
	c.includes = includes

	for _, include := range c.includes {
		fmt.Printf("loaded %s include from %s\n", filepath.Base(include), include)
	}
	for _, file := range files {
		total := append(c.includes, file)
		c.templates[filepath.Base(file)] = template.Must(template.ParseFiles(total...))
		fmt.Printf("loaded %s template from %s\n", filepath.Base(file), file)
	}
}

/* Watcher process to fetch template changes */
func (c *MyController) startWatcher() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(fmt.Sprintf("Error: Unable to init watcher"))
	}
	defer watcher.Close()

	done := make(chan bool)

	go c.watcherEvents(watcher)

	err = watcher.Add("../templates")
	err2 := watcher.Add("../templates/partials")
	if err != nil || err2 != nil {
		panic(fmt.Sprintf("Error: Unable to watch directory"))
	}

	<-done
}

/* Handle watcher events */
func (c *MyController) watcherEvents(watcher *fsnotify.Watcher) {
	for {
		select {
		case event := <-watcher.Events:
			// check if file has been modified
			mod := (event.Op&fsnotify.Create == fsnotify.Create) ||
				(event.Op&fsnotify.Write == fsnotify.Write) ||
				(event.Op&fsnotify.Rename == fsnotify.Rename)

			if mod {
				if filepath.Ext(event.Name) == ".tmpl" {
					fmt.Printf("Op: %d, file: %s, filepath: %s\n",
						event.Op, event.Name, filepath.Ext(event.Name))
					file := event.Name
					total := append(c.includes, file)

					// wait to reload file as some editors delete/rename
					time.Sleep(time.Second)
					fmt.Printf("reloading %s template from %s\n",
						filepath.Base(file), file)
					c.templates[filepath.Base(file)] = template.Must(template.ParseFiles(total...))
					fmt.Printf("reloaded %s template from %s\n",
						filepath.Base(file), file)
				}
			}
		case err := <-watcher.Errors:
			fmt.Printf("error:", err)
		}
	}
}

/* Main */
func main() {
	c := &MyController{page: map[string]string{"Static": "static"}}
	muxer := mux.NewRouter().StrictSlash(true)

	c.initTemplates()
	go c.startWatcher()
	//mongoSession := dbSetup()

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

	n := negroni.Classic()
	n.UseHandler(muxer)
	n.Run(":3000")
}
