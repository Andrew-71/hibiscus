package routes

import (
	"html/template"
	"net/http"
	"net/url"

	"git.a71.su/Andrew71/hibiscus-txt/internal/files"
	"git.a71.su/Andrew71/hibiscus-txt/internal/lang"
	"git.a71.su/Andrew71/hibiscus-txt/internal/server/util"
	"github.com/go-chi/chi/v5"
)

// getNotes calls getEntries for all notes.
func getNotes(w http.ResponseWriter, r *http.Request) {
	// This is suboptimal, but will do...
	description := template.HTML(
		"<a href=\"#\" onclick='newNote(\"" + template.HTMLEscapeString(lang.Translate("prompt.notes")) + "\")'>" + template.HTMLEscapeString(lang.Translate("button.notes")) + "</a>" +
			" <noscript>(" + template.HTMLEscapeString(lang.Translate("noscript.notes")) + ")</noscript>")
	getEntries(w, r, lang.Translate("title.notes"), description, "notes", func(files []string) []Entry {
		var filesFormatted []Entry
		for _, v := range files {
			// titleString := strings.Replace(v, "-", " ", -1) // This would be cool, but what if I need a hyphen?
			filesFormatted = append(filesFormatted, Entry{Title: v, Link: "notes/" + v})
		}
		return filesFormatted
	})
}

// getNote calls getEntry for a note.
func getNote(w http.ResponseWriter, r *http.Request) {
	noteString := chi.URLParam(r, "note")
	if noteString == "" {
		w.WriteHeader(http.StatusBadRequest)
		util.HandleWrite(w.Write([]byte("note not specified")))
		return
	}
	// Handle non-latin note names
	if decodedNote, err := url.QueryUnescape(noteString); err == nil {
		noteString = decodedNote
	}

	getEntry(w, r, noteString, files.DataFile("notes/"+noteString), true)
}

// postNote calls postEntry for a note.
func postNote(w http.ResponseWriter, r *http.Request) {
	noteString := chi.URLParam(r, "note")
	if noteString == "" {
		w.WriteHeader(http.StatusBadRequest)
		util.HandleWrite(w.Write([]byte("note not specified")))
		return
	}
	postEntry(files.DataFile("notes/"+noteString), w, r)
}
