package server

import (
	"embed"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"git.a71.su/Andrew71/hibiscus-txt/internal/config"
	"git.a71.su/Andrew71/hibiscus-txt/internal/server/api"
	"git.a71.su/Andrew71/hibiscus-txt/internal/server/routes"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// public contains the static files e.g. CSS, JS.
//
//go:embed public
var public embed.FS

// Serve starts the app's web server.
func Serve() {
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger, middleware.CleanPath, middleware.StripSlashes)
	r.NotFound(routes.NotFound)

	r.Mount("/", routes.UserRouter) // User-facing routes
	r.Mount("/api", api.ApiRouter)  // API routes

	// Static files
	fs := http.FileServer(http.FS(public))
	r.Handle("/public/*", fs)

	slog.Info("🌺 Website working", "port", config.Cfg.Port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Cfg.Port), r))
}
