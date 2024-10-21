package server

import (
	"log/slog"
	"net/http"

	"git.a71.su/Andrew71/hibiscus-txt/internal/config"
)

// PostConfig calls PostEntry for config file, then reloads the config.
func PostConfig(w http.ResponseWriter, r *http.Request) {
	PostEntry(config.ConfigFile, w, r)
	err := config.Cfg.Reload()
	if err != nil {
		slog.Error("error reloading config", "error", err)
	}
}

// ConfigReloadApi reloads the config. It then redirects back if Referer field is present.
func ConfigReloadApi(w http.ResponseWriter, r *http.Request) {
	err := config.Cfg.Reload()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		HandleWrite(w.Write([]byte(err.Error())))
	}
	if r.Referer() != "" {
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
