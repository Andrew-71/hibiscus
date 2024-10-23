package routes

import (
	"errors"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	"git.a71.su/Andrew71/hibiscus-txt/internal/files"
	"git.a71.su/Andrew71/hibiscus-txt/internal/templates"
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

// getEntries handles showing a list.
func getEntries(w http.ResponseWriter, r *http.Request, title string, description template.HTML, dir string, format formatEntries) {
	filesList, err := files.List(dir)
	if err != nil {
		slog.Error("error reading file list", "directory", dir, "error", err)
		InternalError(w, r)
		return
	}
	var filesFormatted = format(filesList)

	err = templates.List.ExecuteTemplate(w, "base", EntryList{Title: title, Description: description, Entries: filesFormatted})
	if err != nil {
		slog.Error("error executing template", "error", err)
		InternalError(w, r)
		return
	}
}

// getEntry handles showing a single file, editable or otherwise.
func getEntry(w http.ResponseWriter, r *http.Request, title string, filename string, editable bool) {
	entry, err := files.Read(filename)
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
		err = templates.Edit.ExecuteTemplate(w, "base", Entry{Title: title, Content: string(entry)})
	} else {
		err = templates.View.ExecuteTemplate(w, "base", Entry{Title: title, Content: string(entry)})
	}
	if err != nil {
		InternalError(w, r)
		return
	}
}

// postEntry saves value of "text" HTML form component to a file and redirects back to Referer if present.
func postEntry(filename string, w http.ResponseWriter, r *http.Request) {
	err := files.Save(filename, []byte(r.FormValue("text")))
	if err != nil {
		slog.Error("error saving file", "error", err, "file", filename)
	}
	if r.Referer() != "" {
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
		return
	}
}
