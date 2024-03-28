package main

import (
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"log"
	"log/slog"
	"os"
)

var LogFile = "config/log.txt"

// LogInit makes slog output to both stdout and a file if needed
func LogInit() {
	var w io.Writer
	if Cfg.LogToFile {
		f, err := os.OpenFile(LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			slog.Error("error opening log file, logging to stdout", "path", LogFile, "error", err)
			return
		}
		// No defer f.Close() because that breaks the MultiWriter
		w = io.MultiWriter(f, os.Stdout)
	} else {
		w = os.Stdout
	}

	// Make slog and chi use intended format
	slog.SetDefault(slog.New(slog.NewTextHandler(w, nil)))
	middleware.DefaultLogger = middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: log.Default(), NoColor: true})
}
