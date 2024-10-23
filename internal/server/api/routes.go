package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"

	"git.a71.su/Andrew71/hibiscus-txt/internal/config"
	"git.a71.su/Andrew71/hibiscus-txt/internal/files"
	"git.a71.su/Andrew71/hibiscus-txt/internal/server/util"
	"github.com/go-chi/chi/v5"
)

// getFile returns raw contents of a file.
func getFile(filename string, w http.ResponseWriter) {
	fileContents, err := files.Read(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			http.Error(w, "file not found", http.StatusNotFound)
		} else {
			http.Error(w, "error reading found", http.StatusNotFound)
		}
		return
	}
	util.HandleWrite(w.Write(fileContents))
}

// postFile writes contents of Request.Body to a file.
func postFile(filename string, w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		util.HandleWrite(w.Write([]byte("error reading body")))
		return
	}
	err = files.Save(filename, body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		util.HandleWrite(w.Write([]byte("error saving file")))
		return
	}
	util.HandleWrite(w.Write([]byte("wrote to file")))
	w.WriteHeader(http.StatusOK)
}

// fileList returns JSON list of filenames in a directory without extensions or path.
func fileList(directory string, w http.ResponseWriter) {
	filenames, err := files.List(directory)
	if err != nil {
		http.Error(w, "error searching for files", http.StatusInternalServerError)
		return
	}
	filenamesJson, err := json.Marshal(filenames)
	if err != nil {
		http.Error(w, "error marshaling json", http.StatusInternalServerError)
		return
	}
	util.HandleWrite(w.Write(filenamesJson))
}

// getDay returns raw contents of a daily file specified in URL.
func getDay(w http.ResponseWriter, r *http.Request) {
	dayString := chi.URLParam(r, "day")
	if dayString == "" {
		w.WriteHeader(http.StatusBadRequest)
		util.HandleWrite(w.Write([]byte("day not specified")))
		return
	}
	getFile(files.DataFile("day/"+dayString), w)
}

// getNote returns contents of a note specified in URL.
func getNote(w http.ResponseWriter, r *http.Request) {
	noteString := chi.URLParam(r, "note")
	if noteString == "" {
		w.WriteHeader(http.StatusBadRequest)
		util.HandleWrite(w.Write([]byte("note not specified")))
		return
	}
	getFile(files.DataFile("notes/"+noteString), w)
}

// postNote writes contents of Request.Body to a note specified in URL.
func postNote(w http.ResponseWriter, r *http.Request) {
	noteString := chi.URLParam(r, "note")
	if noteString == "" {
		w.WriteHeader(http.StatusBadRequest)
		util.HandleWrite(w.Write([]byte("note not specified")))
		return
	}
	postFile(files.DataFile("notes/"+noteString), w, r)
}

// graceStatus returns "true" if grace period is active, and "false" otherwise.
func graceStatus(w http.ResponseWriter, r *http.Request) {
	value := "false"
	if config.Cfg.Grace() {
		value = "true"
	}
	util.HandleWrite(w.Write([]byte(value)))
	w.WriteHeader(http.StatusOK)
}

// getVersion returns current app version.
func getVersion(w http.ResponseWriter, r *http.Request) {
	util.HandleWrite(w.Write([]byte(config.Info.Version())))
	w.WriteHeader(http.StatusOK)
}

// configReload reloads the config. It then redirects back if Referer field is present.
func configReload(w http.ResponseWriter, r *http.Request) {
	err := config.Cfg.Reload()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		util.HandleWrite(w.Write([]byte(err.Error())))
	}
	if r.Referer() != "" {
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
