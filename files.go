package main

import (
	"bytes"
	"errors"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// ReadFile returns raw contents of a file
func ReadFile(filename string) ([]byte, error) {
	filename = "data/" + path.Clean(filename) + ".txt" // Does this sanitize the path?

	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	fileContents, err := os.ReadFile(filename)
	if err != nil {
		slog.Error("error reading file",
			"error", err,
			"file", filename)
		return nil, err
	}

	return fileContents, nil
}

// SaveFile Writes request's contents to a file
func SaveFile(filename string, contents []byte) error {
	filename = "data/" + filename + ".txt"
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		slog.Error("error opening/making file",
			"error", err,
			"file", filename)
		return err
	}
	if _, err := f.Write(bytes.TrimSpace(contents)); err != nil {
		slog.Error("error writing to file",
			"error", err,
			"file", filename)
		return err
	}
	return nil
}

// ListFiles returns slice of filenames in a directory without extensions or path
func ListFiles(directory string) ([]string, error) {
	filenames, err := filepath.Glob("data/" + directory + "/*.txt")
	if err != nil {
		return nil, err
	}
	for i, file := range filenames {
		file, _ := strings.CutSuffix(filepath.Base(file), filepath.Ext(file))
		filenames[i] = file
	}
	return filenames, nil
}

// ReadToday runs ReadFile with today's date as filename
func ReadToday() ([]byte, error) {
	return ReadFile("day/" + time.Now().In(Cfg.Timezone).Format(time.DateOnly))
}

// SaveToday runs SaveFile with today's date as filename
func SaveToday(contents []byte) error {
	return SaveFile("day/"+time.Now().In(Cfg.Timezone).Format(time.DateOnly), contents)
}
