package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

var templates map[string]*template.Template
var page = map[string]string{"Static": "static"}

/* View Handlers */
func HomeHandler(w http.ResponseWriter, req *http.Request) {
	renderTemplate(w, "home.tmpl", page)
}

func LoginHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		// implement login logic here
		http.Redirect(w, req, "/", 200)
	}
	renderTemplate(w, "login.tmpl", page)
}

func LogoutHandler(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, "/", 200)
}

func StaticHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	asset, file := vars["asset"], vars["file"]
	http.ServeFile(w, req, "static/"+asset+"/"+file)
}

/* Template Renderer */
func renderTemplate(w http.ResponseWriter, name string, data map[string]string) {
	// Ensure the template exists in the map.
	fmt.Printf("beginning template %s\n", name)
	tmpl, ok := templates[name]
	if !ok {
		fmt.Printf("Template %s does not exist.", name)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Printf("rendering template %s\n", name)
	tmpl.ExecuteTemplate(w, "layout", data)

	return nil
}

/* Initialize templates */
func initTemplates() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	templatesDir := "./templates/"
	files, err := filepath.Glob(templatesDir + "*.tmpl")
	if err != nil {
		panic(fmt.Sprintf("Error: Unable to fetch templates"))
	}

	includes, err := filepath.Glob(templatesDir + "partials/*.tmpl")
	if err != nil {
		panic(fmt.Sprintf("Error: Unable to fetch templates"))
	}

	for _, include := range includes {
		fmt.Printf("loaded %s include from %s\n", filepath.Base(include), include)
	}
	for _, file := range files {
		total := append(includes, file)
		templates[filepath.Base(file)] = template.Must(template.ParseFiles(total...))
		fmt.Printf("loaded %s template from %s\n", filepath.Base(file), file)
	}
}

/* Main */
func main() {
	muxer := mux.NewRouter().StrictSlash(true)
	initTemplates()

	/* Routes */
	muxer.HandleFunc("/", HomeHandler)
	muxer.Path("/login").
		Methods("GET", "POST").
		Handler(http.HandlerFunc(LoginHandler))
	muxer.Path("/logout").
		Handler(http.HandlerFunc(LogoutHandler))
	muxer.Path("/static/{asset}/{file}").
		Methods("GET").
		Handler(http.HandlerFunc(StaticHandler))

	n := negroni.Classic()
	n.UseHandler(muxer)
	n.Run(":3000")
}
