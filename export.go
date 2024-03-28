package main

import (
	"archive/zip"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

var ExportPath = "data/export.zip"

// Export saves a .zip archive of the data folder to the passed filename
func Export(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		slog.Error("error creating export archive", "error", err)
		return err
	}
	defer file.Close()

	w := zip.NewWriter(file)
	defer w.Close()

	walker := func(path string, info os.FileInfo, err error) error {
		if path == filename || filepath.Ext(path) == ".zip" { //Ignore export file itself and .zip archives
			return nil
		}
		slog.Debug("export crawling", "path", path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		f, err := w.Create(path)
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}

		return nil
	}
	err = filepath.Walk("data/", walker)
	if err != nil {
		slog.Error("error walking files", "error", err)
		return err
	}

	return nil
}

// GetExport returns a .zip archive with contents of the data folder
func GetExport(w http.ResponseWriter, r *http.Request) {
	err := Export(ExportPath)
	if err != nil {
		http.Error(w, "could not export", http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, "data/export.zip")
}
