package auth

import (
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

// noteLoginFail attempts to log and counteract brute-force attacks.
func noteLoginFail(username string, password string, r *http.Request) {
	slog.Warn("failed auth", "username", username, "password", password, "address", r.RemoteAddr)
	notifyTelegram(fmt.Sprintf(lang.Translate("info.telegram.auth_fail")+":\nusername=%s\npassword=%s\nremote=%s", username, password, r.RemoteAddr))

	attempt := failedLogin{username, password, time.Now()}
	updatedLogins := []failedLogin{attempt}
	for _, attempt := range failedLogins {
		if 100 > time.Since(attempt.Timestamp).Seconds() {
			updatedLogins = append(updatedLogins, attempt)
		}
	}
	failedLogins = updatedLogins

	// At least 3 failed attempts in last 100 seconds -> likely brute-force
	if len(failedLogins) >= 3 && config.Cfg.Scram {
		scram()
	}
}

// notifyTelegram attempts to send a message to the user through Telegram.
func notifyTelegram(msg string) {
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

// scram shuts down the service, useful in case of suspected attack.
func scram() {
	slog.Warn("SCRAM triggered, shutting down")
	notifyTelegram(lang.Translate("info.telegram.scram"))
	os.Exit(0)
}
