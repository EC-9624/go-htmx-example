package internal

import "net/http"

// Count struct to hold the count variable
type Count struct {
    Count int
}

// Global variable count
var count = Count{Count: 0}

func (h *Handlers) OobUpdate(w http.ResponseWriter, r *http.Request) {
	h.renderer.Render(w, r, "5-oob-update.html", count)
}

func (h *Handlers) AddCount(w http.ResponseWriter, r *http.Request) {

    count.Count++
	h.renderer.Render(w, r, "oob-response.html", count)
}
