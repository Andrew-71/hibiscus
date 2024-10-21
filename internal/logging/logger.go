package logging

import (
	"io"
	"log"
	"log/slog"
	"os"

	"git.a71.su/Andrew71/hibiscus-txt/internal/config"
	"github.com/go-chi/chi/v5/middleware"
)

var DebugMode = false

// LogInit makes slog output to both os.Stdout and a file if needed, and sets slog.LevelDebug if enabled.
func LogInit() {
	var w io.Writer
	if config.Cfg.LogToFile {
		f, err := os.OpenFile(config.Cfg.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			slog.Error("error opening log file, logging to stdout only", "path", config.Cfg.LogFile, "error", err)
			return
		}
		// No defer f.Close() because that breaks the MultiWriter
		w = io.MultiWriter(f, os.Stdout)
	} else {
		w = os.Stdout
	}

	// Make slog and chi use intended format
	var opts *slog.HandlerOptions
	if DebugMode {
		opts = &slog.HandlerOptions{Level: slog.LevelDebug}
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(w, opts)))
	middleware.DefaultLogger = middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: log.Default(), NoColor: true})
}
