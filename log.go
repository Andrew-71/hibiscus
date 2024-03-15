package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func AppendLog(input string) error {
	t := time.Now().Format("02-01-2006 15:04") // dd-mm-yyyy HH:MM

	f, err := os.OpenFile("./data/log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening/making file")
		return err
	}

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
	err = AppendLog(string(body))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error appending to log"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("appended to log"))
}
