package main

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// HandleWrite checks for error in ResponseWriter.Write output
func HandleWrite(_ int, err error) {
	if err != nil {
		slog.Error("error writing response", "error", err)
	}
}

// GetFile returns raw contents of a txt file in data directory
func GetFile(filename string, w http.ResponseWriter) {
	filename = "data/" + path.Clean(filename) + ".txt" // Does this sanitize the path?

	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}

	fileContents, err := os.ReadFile(filename)
	if err != nil {
		slog.Error("error reading file",
			"error", err,
			"file", filename)
		http.Error(w, "error reading file", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(fileContents)
	if err != nil {
		http.Error(w, "error sending file", http.StatusInternalServerError)
	}
}

// PostFile Writes request's contents to a txt file in data directory
// TODO: Save to trash to prevent malicious/accidental overrides?
func PostFile(filename string, w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		HandleWrite(w.Write([]byte("error reading body")))
		return
	}

	filename = "data/" + filename + ".txt"
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		slog.Error("error opening/making file",
			"error", err,
			"file", filename)
		HandleWrite(w.Write([]byte("error opening or creating file")))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := f.Write(body); err != nil {
		slog.Error("error writing to file",
			"error", err,
			"file", filename)
		HandleWrite(w.Write([]byte("error writing to file")))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	HandleWrite(w.Write([]byte("wrote to file")))
	w.WriteHeader(http.StatusOK)
}

// ListFiles returns JSON list of filenames in a directory without extensions or path
func ListFiles(directory string, w http.ResponseWriter) {
	filenames, err := filepath.Glob("data/" + directory + "/*.txt")
	if err != nil {
		http.Error(w, "error searching for files", http.StatusInternalServerError)
		return
	}
	for i, file := range filenames {
		file, _ := strings.CutSuffix(filepath.Base(file), filepath.Ext(file))
		filenames[i] = file
	}
	filenamesJson, err := json.Marshal(filenames)
	HandleWrite(w.Write(filenamesJson))
}

// GetDay returns a day specified in URL
func GetDay(w http.ResponseWriter, r *http.Request) {
	dayString := chi.URLParam(r, "day")
	if dayString == "" {
		w.WriteHeader(http.StatusBadRequest)
		HandleWrite(w.Write([]byte("day not specified")))
		return
	}
	GetFile("day/"+dayString, w)
}

// GetToday runs GetFile with today's daily txt
func GetToday(w http.ResponseWriter, _ *http.Request) {
	GetFile("day/"+time.Now().Format("2006-01-02"), w)
}

// PostToday runs PostFile with today's daily txt
func PostToday(w http.ResponseWriter, r *http.Request) {
	PostFile("day/"+time.Now().Format("2006-01-02"), w, r)
}

// GetNote returns a note specified in URL
func GetNote(w http.ResponseWriter, r *http.Request) {
	noteString := chi.URLParam(r, "note")
	if noteString == "" {
		w.WriteHeader(http.StatusBadRequest)
		HandleWrite(w.Write([]byte("note not specified")))
		return
	}
	GetFile("notes/"+noteString, w)
}

// PostNote writes request's contents to a note specified in URL
func PostNote(w http.ResponseWriter, r *http.Request) {
	noteString := chi.URLParam(r, "note")
	if noteString == "" {
		w.WriteHeader(http.StatusBadRequest)
		HandleWrite(w.Write([]byte("note not specified")))
		return
	}
	PostFile("notes/"+noteString, w, r)
}
