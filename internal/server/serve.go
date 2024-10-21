package server

import (
	"embed"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"git.a71.su/Andrew71/hibiscus-txt/internal/config"
	"git.a71.su/Andrew71/hibiscus-txt/internal/files"
	"git.a71.su/Andrew71/hibiscus-txt/internal/lang"
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
	r.NotFound(NotFound)

	// Routes ==========
	userRouter := chi.NewRouter()
	userRouter.Use(BasicAuth)
	userRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		GetEntry(w, r, lang.Translate("title.today"), files.DataFile("day/"+config.Cfg.TodayDate()), true)
	})
	userRouter.Post("/", func(w http.ResponseWriter, r *http.Request) { PostEntry(files.DataFile("day/"+config.Cfg.TodayDate()), w, r) })
	userRouter.Get("/day", GetDays)
	userRouter.Get("/day/{day}", GetDay)
	userRouter.Get("/notes", GetNotes)
	userRouter.Get("/notes/{note}", GetNote)
	userRouter.Post("/notes/{note}", PostNote)
	userRouter.Get("/info", GetInfo)
	userRouter.Get("/readme", func(w http.ResponseWriter, r *http.Request) {
		GetEntry(w, r, "readme.txt", files.DataFile("readme"), true)
	})
	userRouter.Post("/readme", func(w http.ResponseWriter, r *http.Request) { PostEntry(files.DataFile("readme"), w, r) })
	userRouter.Get("/config", func(w http.ResponseWriter, r *http.Request) { GetEntry(w, r, "config.txt", config.ConfigFile, true) })
	userRouter.Post("/config", PostConfig)
	r.Mount("/", userRouter)

	// API =============
	apiRouter := chi.NewRouter()
	apiRouter.Use(BasicAuth)
	apiRouter.Get("/readme", func(w http.ResponseWriter, r *http.Request) { GetFileApi("readme", w) })
	apiRouter.Post("/readme", func(w http.ResponseWriter, r *http.Request) { PostFileApi("readme", w, r) })
	apiRouter.Get("/day", func(w http.ResponseWriter, r *http.Request) { GetFileList("day", w) })
	apiRouter.Get("/day/{day}", GetDayApi)
	apiRouter.Get("/notes", func(w http.ResponseWriter, r *http.Request) { GetFileList("notes", w) })
	apiRouter.Get("/notes/{note}", GetNoteApi)
	apiRouter.Post("/notes/{note}", PostNoteApi)
	apiRouter.Get("/today", func(w http.ResponseWriter, r *http.Request) {
		GetFileApi(files.DataFile("day/"+config.Cfg.TodayDate()), w)
	})
	apiRouter.Post("/today", func(w http.ResponseWriter, r *http.Request) {
		PostEntry(files.DataFile("day/"+config.Cfg.TodayDate()), w, r)
	})
	apiRouter.Get("/export", files.GetExport)
	apiRouter.Get("/grace", GraceActiveApi)
	apiRouter.Get("/version", GetVersionApi)
	apiRouter.Get("/reload", ConfigReloadApi)
	r.Mount("/api", apiRouter)

	// Static files
	fs := http.FileServer(http.FS(public))
	r.Handle("/public/*", fs)

	slog.Info("🌺 Website working", "port", config.Cfg.Port)
	slog.Debug("Debug mode enabled")
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Cfg.Port), r))
}
