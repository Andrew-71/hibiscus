package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"strconv"
)

func Serve() {
	r := chi.NewRouter()
	r.Use(middleware.Logger, middleware.CleanPath, middleware.StripSlashes)
	r.Use(BasicAuth) // TODO: is this good enough?
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		http.ServeFile(w, r, "./pages/error/404.html")
	})

	// Home page
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./pages/index.html")
	})

	// API =============
	apiRouter := chi.NewRouter()

	apiRouter.Get("/readme", func(w http.ResponseWriter, r *http.Request) { GetFile("readme", w) })
	apiRouter.Post("/readme", func(w http.ResponseWriter, r *http.Request) { PostFile("readme", w, r) })
	apiRouter.Get("/log", func(w http.ResponseWriter, r *http.Request) { GetFile("log", w) })
	apiRouter.Post("/log", func(w http.ResponseWriter, r *http.Request) { PostLog(w, r) })

	apiRouter.Get("/day", func(w http.ResponseWriter, r *http.Request) { ListFiles("day", w) })
	apiRouter.Get("/day/{day}", func(w http.ResponseWriter, r *http.Request) { GetDay(w, r) })

	apiRouter.Get("/notes", func(w http.ResponseWriter, r *http.Request) { ListFiles("notes", w) })
	apiRouter.Get("/notes/{note}", func(w http.ResponseWriter, r *http.Request) { GetNote(w, r) })
	apiRouter.Post("/notes/{note}", func(w http.ResponseWriter, r *http.Request) { PostNote(w, r) })

	apiRouter.Get("/today", func(w http.ResponseWriter, r *http.Request) { GetToday(w) })
	apiRouter.Post("/today", func(w http.ResponseWriter, r *http.Request) { PostToday(w, r) })

	r.Mount("/api", apiRouter)

	// Static files
	fs := http.FileServer(http.Dir("public"))
	r.Handle("/public/*", http.StripPrefix("/public/", fs))

	fmt.Println("Website working on port: ", Cfg.Port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(Cfg.Port), r))
}
