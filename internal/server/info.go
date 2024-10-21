package server

import (
	"log/slog"
	"net/http"

	"git.a71.su/Andrew71/hibiscus-txt/internal/config"
	"git.a71.su/Andrew71/hibiscus-txt/internal/templates"
)

// GetInfo renders the info page.
func GetInfo(w http.ResponseWriter, r *http.Request) {
	err := templates.Info.ExecuteTemplate(w, "base", config.Info)
	if err != nil {
		slog.Error("error executing template", "error", err)
		InternalError(w, r)
		return
	}
}

// GetVersionApi returns current app version.
func GetVersionApi(w http.ResponseWriter, r *http.Request) {
	HandleWrite(w.Write([]byte(config.Info.Version)))
	w.WriteHeader(http.StatusOK)
}
