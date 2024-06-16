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

// DataFile modifies file path to ensure it's a .txt inside the data folder.
func DataFile(filename string) string {
	return "data/" + path.Clean(filename) + ".txt"
}

// ReadFile returns contents of a file.
func ReadFile(filename string) ([]byte, error) {
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

// SaveFile Writes contents to a file.
func SaveFile(filename string, contents []byte) error {
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

// ListFiles returns slice of filenames in a directory without extensions or path.
// NOTE: What if I ever want to list non-text files or those outside data directory?
func ListFiles(directory string) ([]string, error) {
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

// GraceActive returns whether the grace period (Cfg.GraceTime) is active. Grace period has minute precision
func GraceActive() bool {
	t := time.Now().In(Cfg.Timezone)
	active := (60*t.Hour() + t.Minute()) < int(Cfg.GraceTime.Minutes())
	if active {
		slog.Debug("grace period active",
			"time", 60*t.Hour()+t.Minute(),
			"grace", Cfg.GraceTime.Minutes())
	}
	return active
}

// TodayDate returns today's formatted date. It accounts for Config.GraceTime.
func TodayDate() string {
	dateFormatted := time.Now().In(Cfg.Timezone).Format(time.DateOnly)
	if GraceActive() {
		dateFormatted = time.Now().In(Cfg.Timezone).AddDate(0, 0, -1).Format(time.DateOnly)
	}
	slog.Debug("today", "time", time.Now().In(Cfg.Timezone).Format(time.DateTime))
	return dateFormatted
}
