package routes

import (
	"log/slog"
	"net/http"

	"git.a71.su/Andrew71/hibiscus-txt/internal/config"
	"git.a71.su/Andrew71/hibiscus-txt/internal/templates"
)

// postConfig calls postEntry for config file, then reloads the config.
func postConfig(w http.ResponseWriter, r *http.Request) {
	postEntry(config.ConfigFile, w, r)
	err := config.Cfg.Reload()
	if err != nil {
		slog.Error("error reloading config", "error", err)
	}
}

// getInfo renders the info page.
func getInfo(w http.ResponseWriter, r *http.Request) {
	err := templates.Info.ExecuteTemplate(w, "base", config.Info)
	if err != nil {
		slog.Error("error executing template", "error", err)
		InternalError(w, r)
		return
	}
}
