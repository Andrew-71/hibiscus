package main

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"io"
	"log/slog"
	"net/http"
	"os"
)

// HandleWrite handles error in output of ResponseWriter.Write.
func HandleWrite(_ int, err error) {
	if err != nil {
		slog.Error("error writing response", "error", err)
	}
}

// GetFileApi returns raw contents of a file.
func GetFileApi(filename string, w http.ResponseWriter) {
	fileContents, err := ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			http.Error(w, "file not found", http.StatusNotFound)
		} else {
			http.Error(w, "error reading found", http.StatusNotFound)
		}
		return
	}
	HandleWrite(w.Write(fileContents))
}

// PostFileApi writes contents of Request.Body to a file.
func PostFileApi(filename string, w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		HandleWrite(w.Write([]byte("error reading body")))
		return
	}
	err = SaveFile(filename, body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		HandleWrite(w.Write([]byte("error saving file")))
		return
	}
	HandleWrite(w.Write([]byte("wrote to file")))
	w.WriteHeader(http.StatusOK)
}

// GetFileList returns JSON list of filenames in a directory without extensions or path.
func GetFileList(directory string, w http.ResponseWriter) {
	filenames, err := ListFiles(directory)
	if err != nil {
		http.Error(w, "error searching for files", http.StatusInternalServerError)
		return
	}
	filenamesJson, err := json.Marshal(filenames)
	if err != nil {
		http.Error(w, "error marshaling json", http.StatusInternalServerError)
		return
	}
	HandleWrite(w.Write(filenamesJson))
}

// GetDayApi returns raw contents of a daily file specified in URL.
func GetDayApi(w http.ResponseWriter, r *http.Request) {
	dayString := chi.URLParam(r, "day")
	if dayString == "" {
		w.WriteHeader(http.StatusBadRequest)
		HandleWrite(w.Write([]byte("day not specified")))
		return
	}
	GetFileApi(DataFile("day/"+dayString), w)
}

// GetNoteApi returns contents of a note specified in URL.
func GetNoteApi(w http.ResponseWriter, r *http.Request) {
	noteString := chi.URLParam(r, "note")
	if noteString == "" {
		w.WriteHeader(http.StatusBadRequest)
		HandleWrite(w.Write([]byte("note not specified")))
		return
	}
	GetFileApi(DataFile("notes/"+noteString), w)
}

// PostNoteApi writes contents of Request.Body to a note specified in URL.
func PostNoteApi(w http.ResponseWriter, r *http.Request) {
	noteString := chi.URLParam(r, "note")
	if noteString == "" {
		w.WriteHeader(http.StatusBadRequest)
		HandleWrite(w.Write([]byte("note not specified")))
		return
	}
	PostFileApi(DataFile("notes/"+noteString), w, r)
}

// GraceActiveApi returns "true" if grace period is active, and "false" otherwise.
func GraceActiveApi(w http.ResponseWriter, r *http.Request) {
	value := "false"
	if GraceActive() {
		value = "true"
	}
	HandleWrite(w.Write([]byte(value)))
	w.WriteHeader(http.StatusOK)
}
