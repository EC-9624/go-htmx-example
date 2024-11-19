package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Define the struct to map the JSON response
type Pokemon struct {
	Name           string  `json:"name"`
	BaseExperience int     `json:"base_experience"`
	Height         int     `json:"height"`
	Weight         int     `json:"weight"`
	Sprites        Sprites `json:"sprites"`
	Types          []Type  `json:"types"`
}

type Sprites struct {
	Other        OtherSprites `json:"other"`
	FrontDefault string       `json:"front_default"`
}

type OtherSprites struct {
	Showdown Showdown `json:"showdown"`
}

type Showdown struct {
	FrontDefault string `json:"front_default"`
}

type Type struct {
	Slot       int        `json:"slot"`
	TypeDetail TypeDetail `json:"type"`
}

type TypeDetail struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func (h *Handlers) ExternalApiHandler(w http.ResponseWriter, r *http.Request){
	
	h.renderer.Render(w, r, "3-external-api.html", nil)
}

func (h *Handlers) HandlePokeRequest(w http.ResponseWriter, r *http.Request){
	if err := r.ParseForm(); err != nil {
		http.Error(w,"unable to parse form", http.StatusInternalServerError)
		return
	}

	resp, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + strings.ToLower(r.FormValue("pokemon")))
	if err != nil {
		http.Error(w, "Unable to fetch new pokemon", http.StatusNotFound)
		fmt.Println(err)
		return
	}
	data := Pokemon{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		http.Error(w, "Unable to parse the Pokemon data", http.StatusUnprocessableEntity)
		return
	}
	h.renderer.Render(w, r, "poke-response.html", data)
}
