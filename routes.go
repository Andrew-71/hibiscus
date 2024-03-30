package main

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type EntryList struct {
	Title   string
	Entries []Entry
}

type Entry struct {
	Title   string
	Content string
	Link    string
}

// NotFound returns a user-friendly 404 error page
func NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	http.ServeFile(w, r, "./pages/error/404.html")
}

// InternalError returns a user-friendly 500 error page
func InternalError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(500)
	http.ServeFile(w, r, "./pages/error/500.html")
}

// GetToday renders HTML page for today's view
func GetToday(w http.ResponseWriter, r *http.Request) {
	day, err := ReadToday()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			day = []byte("")
		} else {
			slog.Error("error reading today's file", "error", err)
			InternalError(w, r)
			return
		}
	}

	files := []string{"./pages/base.html", "./pages/edit.html"}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		InternalError(w, r)
		return
	}

	err = ts.ExecuteTemplate(w, "base", Entry{Title: "Your day so far", Content: string(day)})
	if err != nil {
		InternalError(w, r)
		return
	}
}

// PostToday saves today's entry from form and redirects back to GET
func PostToday(w http.ResponseWriter, r *http.Request) {
	err := SaveToday([]byte(r.FormValue("text")))
	if err != nil {
		slog.Error("error saving today's file", "error", err)
	}
	http.Redirect(w, r, r.Header.Get("Referer"), 302)
}

// GetDays renders HTML page for list of previous days
func GetDays(w http.ResponseWriter, r *http.Request) {
	day, err := ListFiles("day")
	if err != nil {
		slog.Error("error reading today's file", "error", err)
		InternalError(w, r)
		return
	}
	var daysFormatted []Entry
	for i, _ := range day {
		v := day[len(day)-1-i] // This is suboptimal, but reverse order is better here
		dayString := v
		t, err := time.Parse(time.DateOnly, v)
		if err == nil {
			dayString = t.Format("02 Jan 2006")
		}
		if v == time.Now().Format(time.DateOnly) {
			dayString = "Today"
		}
		daysFormatted = append(daysFormatted, Entry{Title: dayString, Link: "day/" + v})
	}

	files := []string{"./pages/base.html", "./pages/list.html"}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		slog.Error("Error parsing template files", "error", err)
		InternalError(w, r)
		return
	}

	err = ts.ExecuteTemplate(w, "base", EntryList{Title: "Previous days", Entries: daysFormatted})
	if err != nil {
		slog.Error("Error executing template", "error", err)
		InternalError(w, r)
		return
	}
}

// GetDay renders HTML page for a specific day
func GetDay(w http.ResponseWriter, r *http.Request) {
	dayString := chi.URLParam(r, "day")
	if dayString == "" {
		w.WriteHeader(http.StatusBadRequest)
		HandleWrite(w.Write([]byte("day not specified")))
		return
	}
	if dayString == time.Now().Format(time.DateOnly) { // today can still be edited
		http.Redirect(w, r, "/", 302)
		return
	}
	day, err := ReadFile("day/" + dayString)
	if err != nil {
		slog.Error("error reading day's file", "error", err, "day", dayString)
		InternalError(w, r)
		return
	}

	files := []string{"./pages/base.html", "./pages/entry.html"}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		InternalError(w, r)
		return
	}

	t, err := time.Parse(time.DateOnly, dayString)
	if err == nil { // This is low priority so silently fail
		dayString = t.Format("02 Jan 2006")
	}

	err = ts.ExecuteTemplate(w, "base", Entry{Content: string(day), Title: dayString})
	if err != nil {
		InternalError(w, r)
		return
	}
}
