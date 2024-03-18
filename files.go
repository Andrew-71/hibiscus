package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// GetFile returns raw contents of a txt file in data directory
func GetFile(filename string, w http.ResponseWriter, r *http.Request) {
	path := "data/" + filename + ".txt" // Can we and should we sanitize this?

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}

	fileContents, err := os.ReadFile(path)
	if err != nil {
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
		w.Write([]byte("error reading body"))
		return
	}

	f, err := os.OpenFile("data/"+filename+".txt", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("error opening/making file")
		w.Write([]byte("error opening or creating file"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := f.Write(body); err != nil {
		fmt.Println("error writing to the file")
		w.Write([]byte("error writing to file"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte("wrote to file"))
	w.WriteHeader(http.StatusOK)
}

// ListFiles returns JSON list of filenames in a directory without extensions or path
func ListFiles(directory string, w http.ResponseWriter, r *http.Request) {
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
	w.Write(filenamesJson)
}

func GetDay(w http.ResponseWriter, r *http.Request) {
	// TODO: This will be different if I move away from chi
	dayString := chi.URLParam(r, "day")
	if dayString == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("day not specified"))
		return
	}
	GetFile("day/"+dayString, w, r)
}

func GetNote(w http.ResponseWriter, r *http.Request) {
	noteString := chi.URLParam(r, "note")
	if noteString == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("note not specified"))
		return
	}
	GetFile("notes/"+noteString, w, r)
}

// GetToday runs GetFile with today's daily txt
func GetToday(w http.ResponseWriter, r *http.Request) {
	GetFile("day/"+time.Now().Format("2006-01-02"), w, r)
}

// PostToday runs PostFile with today's daily txt
func PostToday(w http.ResponseWriter, r *http.Request) {
	PostFile("day/"+time.Now().Format("2006-01-02"), w, r)
}
