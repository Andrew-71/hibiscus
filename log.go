package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// AppendLog adds the input string to the end of the log file with a timestamp
func appendLog(input string) error {
	t := time.Now().Format("2006-01-02 15:04") // yyyy-mm-dd HH:MM

	f, err := os.OpenFile("./data/log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening/making file")
		return err
	}

	input = strings.Replace(input, "\n", "", -1) // Remove newlines to maintain structure
	if _, err := f.Write([]byte(t + " | " + input + "\n")); err != nil {
		fmt.Println("Error appending to the file")
		return err
	}
	if err := f.Close(); err != nil {
		fmt.Println("Error closing the file")
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
