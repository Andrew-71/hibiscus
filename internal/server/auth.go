package server

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"git.a71.su/Andrew71/hibiscus-txt/internal/config"
	"git.a71.su/Andrew71/hibiscus-txt/internal/lang"
)

type failedLogin struct {
	Username  string
	Password  string
	Timestamp time.Time
}

var failedLogins []failedLogin

// NoteLoginFail attempts to log and counteract bruteforce attacks.
func NoteLoginFail(username string, password string, r *http.Request) {
	slog.Warn("failed auth", "username", username, "password", password, "address", r.RemoteAddr)
	NotifyTelegram(fmt.Sprintf(lang.Translate("info.telegram.auth_fail")+":\nusername=%s\npassword=%s\nremote=%s", username, password, r.RemoteAddr))

	attempt := failedLogin{username, password, time.Now()}
	updatedLogins := []failedLogin{attempt}
	for _, attempt := range failedLogins {
		if 100 > time.Since(attempt.Timestamp).Seconds() {
			updatedLogins = append(updatedLogins, attempt)
		}
	}
	failedLogins = updatedLogins

	// At least 3 failed attempts in last 100 seconds -> likely bruteforce
	if len(failedLogins) >= 3 && config.Cfg.Scram {
		Scram()
	}
}

// BasicAuth is a middleware that handles authentication & authorization for the app.
// It uses BasicAuth because I doubt there is a need for something sophisticated in a small hobby project.
// Originally taken from Alex Edwards's https://www.alexedwards.net/blog/basic-authentication-in-go, MIT Licensed (13.03.2024).
func BasicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			// Calculate SHA-256 hashes for equal length in ConstantTimeCompare
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUsernameHash := sha256.Sum256([]byte(config.Cfg.Username))
			expectedPasswordHash := sha256.Sum256([]byte(config.Cfg.Password))

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

// Scram shuts down the service, useful in case of suspected attack.
func Scram() {
	slog.Warn("SCRAM triggered, shutting down")
	NotifyTelegram(lang.Translate("info.telegram.scram"))
	os.Exit(0)
}

// NotifyTelegram attempts to send a message to admin through Telegram.
func NotifyTelegram(msg string) {
	if config.Cfg.TelegramChat == "" || config.Cfg.TelegramToken == "" {
		slog.Debug("ignoring telegram request due to lack of credentials")
		return
	}
	client := &http.Client{}
	data := "chat_id=" + config.Cfg.TelegramChat + "&text=" + msg
	if config.Cfg.TelegramTopic != "" {
		data += "&message_thread_id=" + config.Cfg.TelegramTopic
	}
	req, err := http.NewRequest("POST", "https://api.telegram.org/bot"+config.Cfg.TelegramToken+"/sendMessage", strings.NewReader(data))
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
