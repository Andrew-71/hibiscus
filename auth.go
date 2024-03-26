package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

type failedLogin struct {
	Username  string
	Password  string
	Timestamp time.Time
}

var failedLogins []failedLogin

// NoteLoginFail attempts to counteract bruteforce/spam attacks
func NoteLoginFail(username string, password string, r *http.Request) {
	slog.Warn("failed auth", "username", username, "password", password, "address", r.RemoteAddr)
	NotifyTelegram(fmt.Sprintf("Failed auth attempt in hibiscus:\nusername=%s\npassword=%s\nremote=%s", username, password, r.RemoteAddr))

	attempt := failedLogin{username, password, time.Now()}
	updatedLogins := []failedLogin{attempt}

	for _, attempt := range failedLogins {
		if 100 > time.Now().Sub(attempt.Timestamp).Abs().Seconds() {
			updatedLogins = append(updatedLogins, attempt)
		}
	}

	failedLogins = updatedLogins

	// At least 3 failed attempts in last 100 seconds -> likely bruteforce
	if len(failedLogins) >= 3 {
		Scram()
	}
}

// BasicAuth is a middleware that handles authentication & authorization for the app.
// It uses BasicAuth because I doubt there is a need for something sophisticated in a small hobby project
// Originally taken from https://www.alexedwards.net/blog/basic-authentication-in-go (13.03.2024)
// TODO: why did I have to convert Handler from HandlerFunc?
func BasicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			// Calculate SHA-256 hashes for equal length in ConstantTimeCompare
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUsernameHash := sha256.Sum256([]byte(Cfg.Username))
			expectedPasswordHash := sha256.Sum256([]byte(Cfg.Password))

			usernameMatch := subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1
			passwordMatch := subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1

			if usernameMatch && passwordMatch {
				next.ServeHTTP(w, r)
				return
			} else {
				NoteLoginFail(username, password, r)
			}
		}

		// Unauthorized, inform client that we have auth and return 401
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

// Scram shuts down the service, useful in case of suspected attack
func Scram() {
	slog.Warn("SCRAM triggered, shutting down")
	NotifyTelegram("Hibiscus SCRAM triggered, shutting down")
	os.Exit(0) // TODO: should this be 0 or 1?
}

// NotifyTelegram attempts to send a message to admin through telegram
func NotifyTelegram(msg string) {
	if Cfg.TelegramChat == "" || Cfg.TelegramToken == "" {
		slog.Warn("ignoring telegram request due to lack of credentials")
		return
	}
	client := &http.Client{}
	var data = strings.NewReader("chat_id=" + Cfg.TelegramChat + "&text=" + msg)
	req, err := http.NewRequest("POST", "https://api.telegram.org/bot"+Cfg.TelegramToken+"/sendMessage", data)
	if err != nil {
		slog.Error("failed telegram request", "error", err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("failed telegram request", "error", err)
		return
	}

	if resp.StatusCode != 200 {
		slog.Error("failed telegram request", "status", resp.Status)
	}
}
