package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"log/slog"
	"net/http"
	"strconv"
)

// Serve starts the app's web server
func Serve() {
	r := chi.NewRouter()
	r.Use(middleware.Logger, middleware.CleanPath, middleware.StripSlashes)
	r.Use(BasicAuth) // Is this good enough? Sure hope so
	r.NotFound(NotFound)

	// Routes ==========
	r.Get("/", GetToday)
	r.Post("/", PostToday)
	r.Get("/day", GetDays)
	r.Get("/day/{day}", GetDay)
	r.Get("/notes", GetNotes)
	r.Get("/notes/{note}", GetNote)
	r.Post("/notes/{note}", PostNote)

	// API =============
	apiRouter := chi.NewRouter()
	apiRouter.Get("/readme", func(w http.ResponseWriter, r *http.Request) { GetFile("readme", w) })
	apiRouter.Post("/readme", func(w http.ResponseWriter, r *http.Request) { PostFile("readme", w, r) })
	apiRouter.Get("/day", func(w http.ResponseWriter, r *http.Request) { GetFileList("day", w) })
	apiRouter.Get("/day/{day}", GetDayApi)
	apiRouter.Get("/notes", func(w http.ResponseWriter, r *http.Request) { GetFileList("notes", w) })
	apiRouter.Get("/notes/{note}", GetNoteApi)
	apiRouter.Post("/notes/{note}", PostNoteApi)
	apiRouter.Get("/today", GetTodayApi)
	apiRouter.Post("/today", PostTodayApi)
	apiRouter.Get("/export", GetExport)
	r.Mount("/api", apiRouter)

	// Static files
	fs := http.FileServer(http.Dir("public"))
	r.Handle("/public/*", http.StripPrefix("/public/", fs))

	slog.Info("🌺 Website working", "port", Cfg.Port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(Cfg.Port), r))
}
