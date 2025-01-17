package internal

import (
	"html/template"
	"net/http"
	"path/filepath"
)

// TemplateRenderer handles the template rendering logic
type TemplateRenderer struct {
	TemplatesDir string
	LayoutPath   string
}

// NewTemplateRenderer creates a new template renderer
func NewTemplateRenderer(templatesDir, layoutFile string) *TemplateRenderer {
	return &TemplateRenderer{
		TemplatesDir: templatesDir,
		LayoutPath:   filepath.Join(templatesDir, layoutFile),
	}
}

func (tr *TemplateRenderer) Render(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}) {
	tmplPath := filepath.Join(tr.TemplatesDir, tmpl)
	t, err := template.ParseFiles(tr.LayoutPath, tmplPath)
	if err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Determine which template to execute based on HTMX request
	templateName := "layout"
	if r.Header.Get("HX-Request") == "true" {
		templateName = "HX-Response"
	}

	if err := t.ExecuteTemplate(w, templateName, data); err != nil {
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
	}
}

// Handlers struct to keep all HTTP handlers organized
type Handlers struct {
	renderer *TemplateRenderer
}

// NewHandlers creates a new handlers instance
func NewHandlers(renderer *TemplateRenderer) *Handlers {
	return &Handlers{renderer: renderer}
}

// HomePage handles the root path
func (h *Handlers) HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	h.renderer.Render(w, r, "1-tabs-navigation.html", nil)
}

// WebSocket handles the WebSocket page
func (h *Handlers) WebSocket(w http.ResponseWriter, r *http.Request) {
	h.renderer.Render(w, r, "4-web-socket.html", nil)
}
