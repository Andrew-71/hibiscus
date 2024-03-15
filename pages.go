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

// GetDay gets... a day... page
func GetDay(w http.ResponseWriter, r *http.Request) {
	// TODO: This will be *very* different, `today` func will be needed
	dayString := chi.URLParam(r, "day")
	if dayString == "" {
		dayString = time.Now().Format("02-01-2006") // By default, use today
	}
	GetFile("day/"+dayString, w, r)
}

// GetNote gets... a day... page
func GetNote(w http.ResponseWriter, r *http.Request) {
	// TODO: This will be *very* different, `today` func will be needed
	noteString := chi.URLParam(r, "note")
	if noteString == "" {
		w.WriteHeader(http.StatusNotFound) // TODO: maybe different status fits better?
		w.Write([]byte("note name not given"))
		return
	}
	GetFile("notes/"+noteString, w, r)
}

func PostDayPage(w http.ResponseWriter, r *http.Request) {

}
