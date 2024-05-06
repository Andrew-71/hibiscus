package main

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

type EntryList struct {
	Title       string
	Description string
	Entries     []Entry
}

type Entry struct {
	Title   string
	Content string
	Link    string
}

type formatEntries func([]string) []Entry

var templateFuncs = map[string]interface{}{"translatableText": TranslatableText}
var editTemplate = template.Must(template.New("").Funcs(templateFuncs).ParseFiles("./pages/base.html", "./pages/edit.html"))
var viewTemplate = template.Must(template.New("").Funcs(templateFuncs).ParseFiles("./pages/base.html", "./pages/entry.html"))
var listTemplate = template.Must(template.New("").Funcs(templateFuncs).ParseFiles("./pages/base.html", "./pages/list.html"))

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

// GetToday renders HTML page for today's entry
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

	err = editTemplate.ExecuteTemplate(w, "base", Entry{Title: TranslatableText("title.today"), Content: string(day)})
	if err != nil {
		slog.Error("error executing template", "error", err)
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

// GetEntries is a generic HTML renderer for a list
func GetEntries(w http.ResponseWriter, r *http.Request, title string, description string, dir string, format formatEntries) {
	filesList, err := ListFiles(dir)
	if err != nil {
		slog.Error("error reading file list", "directory", dir, "error", err)
		InternalError(w, r)
		return
	}
	var filesFormatted = format(filesList)

	err = listTemplate.ExecuteTemplate(w, "base", EntryList{Title: title, Description: description, Entries: filesFormatted})
	if err != nil {
		slog.Error("error executing template", "error", err)
		InternalError(w, r)
		return
	}
}

// GetDays renders HTML list of previous days' entries
func GetDays(w http.ResponseWriter, r *http.Request) {
	GetEntries(w, r, TranslatableText("title.days"), "", "day", func(files []string) []Entry {
		var filesFormatted []Entry
		for i := range files {
			v := files[len(files)-1-i] // This is suboptimal, but reverse order is better here
			dayString := v
			t, err := time.Parse(time.DateOnly, v)
			if err == nil {
				dayString = t.Format("02 Jan 2006")
			}

			// Fancy text for today and tomorrow
			// This looks bad, but strings.Title is deprecated, and I'm not importing a golang.org/x package for this...
			// ( chances we ever run into tomorrow are really low)
			if v == TodayDate() {
				dayString = TranslatableText("link.today")
				dayString = strings.ToTitle(string([]rune(dayString)[0])) + string([]rune(dayString)[1:])
			} else if GraceActive() && v > TodayDate() {
				dayString = TranslatableText("link.tomorrow")
				dayString = strings.ToTitle(string([]rune(dayString)[0])) + string([]rune(dayString)[1:])
			}
			filesFormatted = append(filesFormatted, Entry{Title: dayString, Link: "day/" + v})
		}
		return filesFormatted
	})
}

// GetNotes renders HTML list of all notes
func GetNotes(w http.ResponseWriter, r *http.Request) {
	GetEntries(w, r, TranslatableText("title.notes"), TranslatableText("description.notes"), "notes", func(files []string) []Entry {
		var filesFormatted []Entry
		for _, v := range files {
			titleString := strings.Replace(v, "-", " ", -1) // FIXME: what if I need a hyphen?
			filesFormatted = append(filesFormatted, Entry{Title: titleString, Link: "notes/" + v})
		}
		return filesFormatted
	})
}

func GetEntry(w http.ResponseWriter, r *http.Request, title string, filename string, editable bool) {
	entry, err := ReadFile(filename)
	if err != nil {
		if editable && errors.Is(err, os.ErrNotExist) {
			entry = []byte("")
		} else {
			slog.Error("error reading entry file", "error", err, "file", filename)
			InternalError(w, r)
			return
		}
	}

	files := []string{"./pages/base.html"}
	if editable {
		files = append(files, "./pages/edit.html")
	} else {
		files = append(files, "./pages/entry.html")
	}

	if editable {
		err = editTemplate.ExecuteTemplate(w, "base", Entry{Title: title, Content: string(entry)})
	} else {
		err = viewTemplate.ExecuteTemplate(w, "base", Entry{Title: title, Content: string(entry)})
	}
	if err != nil {
		InternalError(w, r)
		return
	}
}

// GetDay renders HTML page for a specific day entry
func GetDay(w http.ResponseWriter, r *http.Request) {
	dayString := chi.URLParam(r, "day")
	if dayString == "" {
		w.WriteHeader(http.StatusBadRequest)
		HandleWrite(w.Write([]byte("day not specified")))
		return
	}
	if dayString == TodayDate() { // Today can still be edited
		http.Redirect(w, r, "/", 302)
		return
	}

	title := dayString
	t, err := time.Parse(time.DateOnly, dayString)
	if err == nil { // This is low priority so silently fail
		title = t.Format("02 Jan 2006")
	}

	GetEntry(w, r, title, "day/"+dayString, false)
}

// GetNote renders HTML page for a note
func GetNote(w http.ResponseWriter, r *http.Request) {
	noteString := chi.URLParam(r, "note")
	if noteString == "" {
		w.WriteHeader(http.StatusBadRequest)
		HandleWrite(w.Write([]byte("note not specified")))
		return
	}

	GetEntry(w, r, noteString, "notes/"+noteString, true)
}

// PostNote saves a note form and redirects back to GET
func PostNote(w http.ResponseWriter, r *http.Request) {
	noteString := chi.URLParam(r, "note")
	if noteString == "" {
		w.WriteHeader(http.StatusBadRequest)
		HandleWrite(w.Write([]byte("note not specified")))
		return
	}
	err := SaveFile("notes/"+noteString, []byte(r.FormValue("text")))
	if err != nil {
		slog.Error("error saving a note", "note", noteString, "error", err)
	}
	http.Redirect(w, r, r.Header.Get("Referer"), 302)
}
