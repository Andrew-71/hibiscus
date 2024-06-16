package main

import (
	"embed"
	"errors"
	"github.com/go-chi/chi/v5"
	"html/template"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type EntryList struct {
	Title       string
	Description template.HTML
	Entries     []Entry
}

type Entry struct {
	Title   string
	Content string
	Link    string
}

type formatEntries func([]string) []Entry

// Public contains the static files e.g. CSS, JS.
//
//go:embed public
var Public embed.FS

// Pages contains the HTML templates used by the app.
//
//go:embed pages
var Pages embed.FS

// EmbeddedPage returns contents of a file in Pages while "handling" potential errors.
func EmbeddedPage(name string) []byte {
	data, err := Pages.ReadFile(name)
	if err != nil {
		slog.Error("error reading embedded file", "err", err)
	}
	return data
}

var templateFuncs = map[string]interface{}{
	"translatableText": TranslatableText,
	"info":             func() AppInfo { return Info },
	"config":           func() Config { return Cfg },
}
var editTemplate = template.Must(template.New("").Funcs(templateFuncs).ParseFS(Pages, "pages/base.html", "pages/edit.html"))
var viewTemplate = template.Must(template.New("").Funcs(templateFuncs).ParseFS(Pages, "pages/base.html", "pages/entry.html"))
var listTemplate = template.Must(template.New("").Funcs(templateFuncs).ParseFS(Pages, "pages/base.html", "pages/list.html"))

// NotFound returns a user-friendly 404 error page.
func NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	HandleWrite(w.Write(EmbeddedPage("pages/error/404.html")))
}

// InternalError returns a user-friendly 500 error page.
func InternalError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(500)
	HandleWrite(w.Write(EmbeddedPage("pages/error/500.html")))
}

// GetEntries handles showing a list.
func GetEntries(w http.ResponseWriter, r *http.Request, title string, description template.HTML, dir string, format formatEntries) {
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

// GetDays calls GetEntries for previous days' entries.
func GetDays(w http.ResponseWriter, r *http.Request) {
	description := template.HTML(
		"<a href=\"#footer\">" + template.HTMLEscapeString(TranslatableText("prompt.days")) + "</a>")
	GetEntries(w, r, TranslatableText("title.days"), description, "day", func(files []string) []Entry {
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
			// (chances we ever run into tomorrow are really low)
			if v == TodayDate() {
				dayString = TranslatableText("link.today")
				dayString = strings.ToTitle(string([]rune(dayString)[0])) + string([]rune(dayString)[1:])
			} else if v > TodayDate() {
				dayString = TranslatableText("link.tomorrow")
				dayString = strings.ToTitle(string([]rune(dayString)[0])) + string([]rune(dayString)[1:])
			}
			filesFormatted = append(filesFormatted, Entry{Title: dayString, Link: "day/" + v})
		}
		return filesFormatted
	})
}

// GetNotes calls GetEntries for all notes.
func GetNotes(w http.ResponseWriter, r *http.Request) {
	// This is suboptimal, but will do...
	description := template.HTML(
		"<a href=\"#\" onclick='newNote(\"" + template.HTMLEscapeString(TranslatableText("prompt.notes")) + "\")'>" + template.HTMLEscapeString(TranslatableText("button.notes")) + "</a>" +
			" <noscript>(" + template.HTMLEscapeString(TranslatableText("noscript.notes")) + ")</noscript>")
	GetEntries(w, r, TranslatableText("title.notes"), description, "notes", func(files []string) []Entry {
		var filesFormatted []Entry
		for _, v := range files {
			// titleString := strings.Replace(v, "-", " ", -1) // This would be cool, but what if I need a hyphen?
			filesFormatted = append(filesFormatted, Entry{Title: v, Link: "notes/" + v})
		}
		return filesFormatted
	})
}

// GetEntry handles showing a single file, editable or otherwise.
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

// PostEntry saves value of "text" HTML form component to a file and redirects back to Referer if present.
func PostEntry(filename string, w http.ResponseWriter, r *http.Request) {
	err := SaveFile(filename, []byte(r.FormValue("text")))
	if err != nil {
		slog.Error("error saving file", "error", err, "file", filename)
	}
	if r.Referer() != "" {
		http.Redirect(w, r, r.Header.Get("Referer"), 302)
		return
	}
}

// GetDay calls GetEntry for a day entry.
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

	GetEntry(w, r, title, DataFile("day/"+dayString), false)
}

// GetNote calls GetEntry for a note.
func GetNote(w http.ResponseWriter, r *http.Request) {
	noteString := chi.URLParam(r, "note")
	if noteString == "" {
		w.WriteHeader(http.StatusBadRequest)
		HandleWrite(w.Write([]byte("note not specified")))
		return
	}
	// Handle non-latin note names
	if decodedNote, err := url.QueryUnescape(noteString); err == nil {
		noteString = decodedNote
	}

	GetEntry(w, r, noteString, DataFile("notes/"+noteString), true)
}

// PostNote calls PostEntry for a note.
func PostNote(w http.ResponseWriter, r *http.Request) {
	noteString := chi.URLParam(r, "note")
	if noteString == "" {
		w.WriteHeader(http.StatusBadRequest)
		HandleWrite(w.Write([]byte("note not specified")))
		return
	}
	PostEntry(DataFile("notes/"+noteString), w, r)
}
