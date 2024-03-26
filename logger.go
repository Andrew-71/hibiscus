package main

import (
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"log"
	"log/slog"
	"os"
)

var LogFile = "config/log.txt"

// LogInit makes slog output to both stdout and a file
func LogInit() {
	f, err := os.OpenFile(LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil))) // Fallback to stdout
		slog.Error("error opening log file, logging to stdout", "path", LogFile, "error", err)
		return
	}
	// No defer f.Close() because that breaks the MultiWriter
	w := io.MultiWriter(f, os.Stdout)
	slog.SetDefault(slog.New(slog.NewTextHandler(w, nil)))

	// Make chi log to file too
	middleware.DefaultLogger = middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: log.Default(), NoColor: true})
}
