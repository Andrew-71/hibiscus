package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var DEFAULT_EXTENSION = "txt"

func GetFile(filename string, w http.ResponseWriter, r *http.Request) {
	filenames, err := filepath.Glob("data/" + filename + ".*") // .txt, .md, anything
	if err != nil {
		http.Error(w, "error finding file", http.StatusInternalServerError)
		return
	}

	if len(filenames) == 0 {
		http.Error(w, "no matching files found", http.StatusNotFound)
		return
	} else if len(filenames) > 1 {
		http.Error(w, "several matching files found ("+strconv.Itoa(len(filenames))+")", http.StatusInternalServerError) // TODO: Better handling, duh
		return
	}

	fileContents, err := os.ReadFile(filenames[0])
	if err != nil {
		http.Error(w, "error reading file", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(fileContents)
	if err != nil {
		http.Error(w, "error sending file", http.StatusInternalServerError)
	}
}

// PostFile TODO: Save to trash to prevent malicious/accidental ovverrides?
func PostFile(filename string, w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error reading body"))
		return
	}

	filenames, err := filepath.Glob("data/" + filename + ".*") // .txt, .md, anything
	if err != nil {
		http.Error(w, "error searching for file", http.StatusInternalServerError)
		return
	}

	var filenameFinal string
	if len(filenames) == 0 {
		// Create new file and write
		filenameFinal = "data/" + filename + "." + DEFAULT_EXTENSION
	} else if len(filenames) > 1 {
		http.Error(w, "several matching files found ("+strconv.Itoa(len(filenames))+")", http.StatusInternalServerError) // TODO: Better handling, duh
		return
	} else {
		filenameFinal = filenames[0]
		fmt.Println(filenameFinal)
		fmt.Println(filenames)
	}

	f, err := os.OpenFile(filenameFinal, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening/making file")
		return
	}

	if _, err := f.Write(body); err != nil {
		fmt.Println("Error writing to the file")
	}
}

// ListFiles returns JSON of filenames in a directory without extensions or path
func ListFiles(directory string, w http.ResponseWriter, r *http.Request) {
	filenames, err := filepath.Glob("data/" + directory + "/*") // .txt, .md, anything
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
