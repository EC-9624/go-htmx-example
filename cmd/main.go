package main

import (
	"html/template"
	"net/http"
	"path/filepath"
)


func renderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}) {
	tmplPath := filepath.Join("htmx", "templates", tmpl)
	t, err := template.ParseFiles(
		filepath.Join("htmx", "templates", "layout.html"),
		tmplPath,
	)
	if err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the request is an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		// Render only the content block
		if err := t.ExecuteTemplate(w, "content", data); err != nil {
			http.Error(w, "Error executing content template: "+err.Error(), http.StatusInternalServerError)
		}
	} else {
		// Render the full layout
		if err := t.ExecuteTemplate(w, "layout", data); err != nil {
			http.Error(w, "Error executing layout template: "+err.Error(), http.StatusInternalServerError)
		}
	}
	
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r, "home.html", nil)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r, "about.html", nil)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r, "contact.html", nil)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		homeHandler(w, r)
	})
	http.HandleFunc("/home", homeHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/contact", contactHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
