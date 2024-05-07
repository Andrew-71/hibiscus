package main

import (
	"html/template"
	"log/slog"
	"net/http"
)

var infoTemplate = template.Must(template.New("").Funcs(templateFuncs).ParseFiles("./pages/base.html", "./pages/info.html"))

type HibiscusInfo struct {
	Version    string
	SourceLink string
}

// Info contains app information
var Info = HibiscusInfo{
	Version:    "0.2.0",
	SourceLink: "https://git.a71.su/Andrew71/hibiscus",
}

// GetInfo renders the info page
func GetInfo(w http.ResponseWriter, r *http.Request) {
	err := infoTemplate.ExecuteTemplate(w, "base", Info)
	if err != nil {
		slog.Error("error executing template", "error", err)
		InternalError(w, r)
		return
	}
}