package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"strconv"
)

const PORT = 7101 // TODO: Obviously don't just declare port here

func Serve() {
	r := chi.NewRouter()
	r.Use(middleware.Logger, middleware.CleanPath, middleware.StripSlashes)
	r.Use(basicAuth) // TODO: ..duh!

	// Home page
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./pages/index.html")
	})

	// API =============
	apiRouter := chi.NewRouter()

	apiRouter.Get("/readme", func(w http.ResponseWriter, r *http.Request) { GetFile("readme", w, r) })
	apiRouter.Get("/log", func(w http.ResponseWriter, r *http.Request) { GetFile("log", w, r) })
	apiRouter.Get("/agenda", func(w http.ResponseWriter, r *http.Request) { GetFile("agenda", w, r) })

	apiRouter.Post("/readme", func(w http.ResponseWriter, r *http.Request) { PostFile("readme", w, r) })
	apiRouter.Post("/agenda", func(w http.ResponseWriter, r *http.Request) { PostFile("agenda", w, r) })

	apiRouter.Get("/day", func(w http.ResponseWriter, r *http.Request) { ListFiles("day", w, r) })
	apiRouter.Get("/day/{day}", func(w http.ResponseWriter, r *http.Request) { GetDay(w, r) })

	apiRouter.Get("/notes", func(w http.ResponseWriter, r *http.Request) { ListFiles("notes", w, r) })
	apiRouter.Get("/notes/{note}", func(w http.ResponseWriter, r *http.Request) { GetNote(w, r) })

	apiRouter.Post("/log", func(w http.ResponseWriter, r *http.Request) { PostLog(w, r) })

	r.Mount("/api", apiRouter)

	r.NotFound(NotFound)

	// Static files
	fs := http.FileServer(http.Dir("public"))
	r.Handle("/public/*", http.StripPrefix("/public/", fs))

	fmt.Println("Website working on port: ", PORT)
	_ = http.ListenAndServe(":"+strconv.Itoa(PORT), r)
}
