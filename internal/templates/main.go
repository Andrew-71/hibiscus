package templates

import (
	"embed"
	"html/template"
	"log/slog"

	"git.a71.su/Andrew71/hibiscus-txt/internal/config"
	"git.a71.su/Andrew71/hibiscus-txt/internal/lang"
)

// pages contains the HTML templates used by the app.
//
//go:embed pages
var pages embed.FS

// EmbeddedPage returns contents of a file in Pages while "handling" potential errors.
func EmbeddedPage(name string) []byte {
	data, err := pages.ReadFile(name)
	if err != nil {
		slog.Error("error reading embedded file", "err", err)
		return []byte("")
	}
	return data
}

var templateFuncs = map[string]interface{}{
	"translate": lang.Translate,
	"info":             func() config.AppInfo { return config.Info },
	"config":           func() config.Config { return config.Cfg },
}
var Edit = template.Must(template.New("").Funcs(templateFuncs).ParseFS(pages, "pages/base.html", "pages/edit.html"))
var View = template.Must(template.New("").Funcs(templateFuncs).ParseFS(pages, "pages/base.html", "pages/entry.html"))
var List = template.Must(template.New("").Funcs(templateFuncs).ParseFS(pages, "pages/base.html", "pages/list.html"))

var Info = template.Must(template.New("").Funcs(templateFuncs).ParseFS(pages, "pages/base.html", "pages/info.html"))

var Template404 = template.Must(template.New("404").Funcs(templateFuncs).ParseFS(pages, "pages/error/404.html"))
var Template500 = template.Must(template.New("500").Funcs(templateFuncs).ParseFS(pages, "pages/error/500.html"))
