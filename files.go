package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// GetFile returns raw contents of a txt file in data directory
func GetFile(filename string, w http.ResponseWriter, r *http.Request) {
	path := "data/" + filename + ".txt" // Can we and should we sanitize this?

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		NotFound(w, r)
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
