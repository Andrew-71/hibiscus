package main

import (
	"html/template"
	"log/slog"
	"net/http"
)

var infoTemplate = template.Must(template.New("").Funcs(templateFuncs).ParseFS(Pages, "pages/base.html", "pages/info.html"))

type AppInfo struct {
	Version    string
	SourceLink string
}

// Info contains app information.
var Info = AppInfo{
	Version:    "1.1.4",
	SourceLink: "https://git.a71.su/Andrew71/hibiscus",
}

// GetInfo renders the info page.
func GetInfo(w http.ResponseWriter, r *http.Request) {
	err := infoTemplate.ExecuteTemplate(w, "base", Info)
	if err != nil {
		slog.Error("error executing template", "error", err)
		InternalError(w, r)
		return
	}
}

// GetVersionApi returns current app version.
func GetVersionApi(w http.ResponseWriter, r *http.Request) {
	HandleWrite(w.Write([]byte(Info.Version)))
	w.WriteHeader(http.StatusOK)
}
