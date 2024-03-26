package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"log/slog"
	"net/http"
	"strconv"
)

func Serve() {
	r := chi.NewRouter()
	r.Use(middleware.Logger, middleware.CleanPath, middleware.StripSlashes)
	r.Use(BasicAuth) // Is this good enough? Sure hope so
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		http.ServeFile(w, r, "./pages/error/404.html")
	})

	// Home page
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./pages/index.html")
	})
	r.Get("/days", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./pages/days.html")
	})

	// API =============
	apiRouter := chi.NewRouter()

	apiRouter.Get("/readme", func(w http.ResponseWriter, r *http.Request) { GetFile("readme", w) })
	apiRouter.Post("/readme", func(w http.ResponseWriter, r *http.Request) { PostFile("readme", w, r) })
	apiRouter.Get("/log", func(w http.ResponseWriter, r *http.Request) { GetFile("log", w) })
	apiRouter.Post("/log", PostLog)

	apiRouter.Get("/day", func(w http.ResponseWriter, r *http.Request) { ListFiles("day", w) })
	apiRouter.Get("/day/{day}", GetDay)

	apiRouter.Get("/notes", func(w http.ResponseWriter, r *http.Request) { ListFiles("notes", w) })
	apiRouter.Get("/notes/{note}", GetNote)
	apiRouter.Post("/notes/{note}", PostNote)

	apiRouter.Get("/today", GetToday)
	apiRouter.Post("/today", PostToday)

	apiRouter.Get("/export", GetExport)

	r.Mount("/api", apiRouter)

	// Static files
	fs := http.FileServer(http.Dir("public"))
	r.Handle("/public/*", http.StripPrefix("/public/", fs))

	slog.Info("🌺 Website working", "port", Cfg.Port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(Cfg.Port), r))
}
