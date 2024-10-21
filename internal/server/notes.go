package server

import (
	"html/template"
	"net/http"
	"net/url"

	"git.a71.su/Andrew71/hibiscus-txt/internal/files"
	"git.a71.su/Andrew71/hibiscus-txt/internal/lang"
	"github.com/go-chi/chi/v5"
)

// GetNotes calls GetEntries for all notes.
func GetNotes(w http.ResponseWriter, r *http.Request) {
	// This is suboptimal, but will do...
	description := template.HTML(
		"<a href=\"#\" onclick='newNote(\"" + template.HTMLEscapeString(lang.Translate("prompt.notes")) + "\")'>" + template.HTMLEscapeString(lang.Translate("button.notes")) + "</a>" +
			" <noscript>(" + template.HTMLEscapeString(lang.Translate("noscript.notes")) + ")</noscript>")
	GetEntries(w, r, lang.Translate("title.notes"), description, "notes", func(files []string) []Entry {
		var filesFormatted []Entry
		for _, v := range files {
			// titleString := strings.Replace(v, "-", " ", -1) // This would be cool, but what if I need a hyphen?
			filesFormatted = append(filesFormatted, Entry{Title: v, Link: "notes/" + v})
		}
		return filesFormatted
	})
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

	GetEntry(w, r, noteString, files.DataFile("notes/"+noteString), true)
}

// PostNote calls PostEntry for a note.
func PostNote(w http.ResponseWriter, r *http.Request) {
	noteString := chi.URLParam(r, "note")
	if noteString == "" {
		w.WriteHeader(http.StatusBadRequest)
		HandleWrite(w.Write([]byte("note not specified")))
		return
	}
	PostEntry(files.DataFile("notes/"+noteString), w, r)
}
