package files

import (
	"bytes"
	"errors"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// DataFile modifies file path to ensure it's a .txt inside the data folder.
func DataFile(filename string) string {
	return "data/" + path.Clean(filename) + ".txt"
}

// Read returns contents of a file.
func Read(filename string) ([]byte, error) {
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		slog.Error("error reading file", "error", err, "file", filename)
		return nil, err
	}
	return fileContents, nil
}

// Save Writes contents to a file.
func Save(filename string, contents []byte) error {
	contents = bytes.TrimSpace(contents)
	if len(contents) == 0 { // Delete empty files
		err := os.Remove(filename)
		slog.Error("error deleting empty file", "error", err, "file", filename)
		return err
	}
	err := os.MkdirAll(path.Dir(filename), 0755) // Create dir in case it doesn't exist yet to avoid errors
	if err != nil {
		slog.Error("error creating directory", "error", err, "file", filename)
		return err
	}
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		slog.Error("error opening/creating file", "error", err, "file", filename)
		return err
	}
	if _, err := f.Write(contents); err != nil {
		slog.Error("error writing to file", "error", err, "file", filename)
		return err
	}
	return nil
}

// List returns slice of filenames in a directory without extensions or path.
// NOTE: What if I ever want to list non-text files or those outside data directory?
func List(directory string) ([]string, error) {
	filenames, err := filepath.Glob("data/" + path.Clean(directory) + "/*.txt")
	if err != nil {
		return nil, err
	}
	for i, file := range filenames {
		file, _ := strings.CutSuffix(filepath.Base(file), filepath.Ext(file))
		filenames[i] = file
	}
	return filenames, nil
}
