package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
)

// NotFound returns a user-friendly 404 error page
func NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	http.ServeFile(w, r, "./pages/error/404.html")
}

func GetDay(w http.ResponseWriter, r *http.Request) {
	// TODO: This will be *very* different, `today` func will be needed
	dayString := chi.URLParam(r, "day")
	if dayString == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("day not specified"))
		return
	}
	GetFile("day/"+dayString, w, r)
}

func GetNote(w http.ResponseWriter, r *http.Request) {
	// TODO: This will be *very* different, `today` func will be needed
	noteString := chi.URLParam(r, "note")
	if noteString == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("note not specified"))
		return
	}
	GetFile("notes/"+noteString, w, r)
}

func GetToday(w http.ResponseWriter, r *http.Request) {
	GetFile("day/"+time.Now().Format("2006-01-02"), w, r)
}

func PostToday(w http.ResponseWriter, r *http.Request) {
	PostFile("day/"+time.Now().Format("2006-01-02"), w, r)
}
