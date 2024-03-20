package main

import (
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

// AppendLog adds the input string to the end of the log file with a timestamp
func appendLog(input string) error {
	t := time.Now().Format("2006-01-02 15:04") // yyyy-mm-dd HH:MM
	filename := "data/log.txt"

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		slog.Error("error opening/making file",
			"error", err,
			"file", filename)
		return err
	}

	input = strings.Replace(input, "\n", "", -1) // Remove newlines to maintain structure
	if _, err := f.Write([]byte(t + " | " + input + "\n")); err != nil {
		slog.Error("error appending to file",
			"error", err,
			"file", filename)
		return err
	}
	if err := f.Close(); err != nil {
		slog.Error("error closing file",
			"error", err,
			"file", filename)
		return err
	}
	return nil
}

func PostLog(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error reading body"))
		return
	}
	err = appendLog(string(body))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error appending to log"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("appended to log"))
}
