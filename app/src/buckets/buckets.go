package buckets

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"gopkg.in/fsnotify.v1"
)

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
