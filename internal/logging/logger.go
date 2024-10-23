package logging

import (
	"io"
	"log"
	"log/slog"
	"os"
	"path"
	"time"

	"git.a71.su/Andrew71/hibiscus-txt/internal/config"
	"github.com/go-chi/chi/v5/middleware"
)

var DebugMode = false

// file returns the appropriate filename for log
// (log_dir/hibiscus_YYYY-MM-DD_HH:MM:SS.log)
func file() string {
	return config.Cfg.LogDir + "/hibiscus_" + time.Now().In(config.Cfg.Timezone).Format("2006-01-02_15:04:05") + ".log"
}

// LogInit makes slog output to both os.Stdout and a file if needed, and sets slog.LevelDebug if enabled.
func LogInit() {
	logFile := file()
	var w io.Writer = os.Stdout
	if config.Cfg.LogToFile {
		// Create dir in case it doesn't exist yet to avoid errors
		err := os.MkdirAll(path.Dir(logFile), 0755)
		if err != nil {
			slog.Error("error creating log dir, logging to stdout only", "path", path.Dir(logFile), "error", err)
		} else {
			f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				slog.Error("error opening log file, logging to stdout only", "path", logFile, "error", err)
				return
			}
			// No defer f.Close() because that breaks the MultiWriter
			w = io.MultiWriter(f, os.Stdout)
		}
	}

	// Make slog and chi use intended format
	var opts *slog.HandlerOptions
	if DebugMode {
		opts = &slog.HandlerOptions{Level: slog.LevelDebug}
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(w, opts)))
	middleware.DefaultLogger = middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: log.Default(), NoColor: true})
	slog.Debug("Debug mode enabled") // This string is only shown if debugging
}
