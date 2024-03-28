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

type DayData struct {
	Day  string
	Date string
}

type List struct {
	Title   string
	Entries []ListEntry
}

type ListEntry struct {
	Name string
	Link string
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

	files := []string{"./pages/base.html", "./pages/index.html"}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		InternalError(w, r)
		return
	}

	err = ts.ExecuteTemplate(w, "base", DayData{Day: string(day)})
	if err != nil {
		InternalError(w, r)
		return
	}
}

func PostToday(w http.ResponseWriter, r *http.Request) {
	err := SaveToday([]byte(r.FormValue("day")))
	if err != nil {
		slog.Error("error saving today's file", "error", err)
	}
	http.Redirect(w, r, r.Header.Get("Referer"), 302)
}

func GetDays(w http.ResponseWriter, r *http.Request) {
	day, err := ListFiles("day")
	if err != nil {
		slog.Error("error reading today's file", "error", err)
		InternalError(w, r)
		return
	}
	var daysFormatted []ListEntry
	for _, v := range day {
		dayString := v
		t, err := time.Parse(time.DateOnly, v)
		if err == nil {
			dayString = t.Format("02 Jan 2006")
		}
		daysFormatted = append(daysFormatted, ListEntry{Name: dayString, Link: v})
	}

	files := []string{"./pages/base.html", "./pages/days.html"}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		slog.Error("Error parsing template files", "error", err)
		InternalError(w, r)
		return
	}

	err = ts.ExecuteTemplate(w, "base", List{Title: "Previous days", Entries: daysFormatted})
	if err != nil {
		slog.Error("Error executing template", "error", err)
		InternalError(w, r)
		return
	}
}

func GetDay(w http.ResponseWriter, r *http.Request) {
	dayString := chi.URLParam(r, "day")
	if dayString == "" {
		w.WriteHeader(http.StatusBadRequest)
		HandleWrite(w.Write([]byte("day not specified")))
		return
	}
	day, err := ReadFile("day/" + dayString)
	if err != nil {
		slog.Error("error reading day's file", "error", err, "day", dayString)
		InternalError(w, r)
		return
	}

	files := []string{"./pages/base.html", "./pages/day.html"}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		InternalError(w, r)
		return
	}

	t, err := time.Parse(time.DateOnly, dayString)
	if err == nil { // This is low priority so silently fail
		dayString = t.Format("02 Jan 2006")
	}

	err = ts.ExecuteTemplate(w, "base", DayData{Day: string(day), Date: dayString})
	if err != nil {
		InternalError(w, r)
		return
	}
}
