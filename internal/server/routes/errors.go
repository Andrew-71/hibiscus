package routes

import (
	"log/slog"
	"net/http"

	"git.a71.su/Andrew71/hibiscus-txt/internal/templates"
)

// NotFound returns a user-friendly 404 error page.
func NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)

	err := templates.Template404.Execute(w, nil)
	if err != nil {
		slog.Error("error rendering error 404 page", "error", err)
		InternalError(w, r)
		return
	}
}

// InternalError returns a user-friendly 500 error page.
func InternalError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(500)

	err := templates.Template500.Execute(w, nil)
	if err != nil { // Well this is awkward
		slog.Error("error rendering error 500 page", "error", err)
		http.Error(w, "500. Something went *very* wrong.", http.StatusInternalServerError)
		return
	}
}
