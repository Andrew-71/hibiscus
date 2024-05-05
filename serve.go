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
	r.NotFound(NotFound)

	// Routes ==========
	userRouter := chi.NewRouter()
	userRouter.Use(BasicAuth)
	userRouter.Get("/", GetToday)
	userRouter.Post("/", PostToday)
	userRouter.Get("/day", GetDays)
	userRouter.Get("/day/{day}", GetDay)
	userRouter.Get("/notes", GetNotes)
	userRouter.Get("/notes/{note}", GetNote)
	userRouter.Post("/notes/{note}", PostNote)
	r.Mount("/", userRouter)

	// API =============
	apiRouter := chi.NewRouter()
	apiRouter.Use(BasicAuth)
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
	slog.Debug("Debug mode enabled")
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(Cfg.Port), r))
}
