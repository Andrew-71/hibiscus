package api

import (
	"net/http"

	"git.a71.su/Andrew71/hibiscus-txt/internal/config"
	"git.a71.su/Andrew71/hibiscus-txt/internal/files"
	"git.a71.su/Andrew71/hibiscus-txt/internal/server/auth"
	"github.com/go-chi/chi/v5"
)

var ApiRouter *chi.Mux

func init() {
	ApiRouter = chi.NewRouter()
	ApiRouter.Use(auth.BasicAuth)
	ApiRouter.Get("/readme", func(w http.ResponseWriter, r *http.Request) { getFile("readme", w) })
	ApiRouter.Post("/readme", func(w http.ResponseWriter, r *http.Request) { postFile("readme", w, r) })
	ApiRouter.Get("/day", func(w http.ResponseWriter, r *http.Request) { fileList("day", w) })
	ApiRouter.Get("/day/{day}", getDay)
	ApiRouter.Get("/notes", func(w http.ResponseWriter, r *http.Request) { fileList("notes", w) })
	ApiRouter.Get("/notes/{note}", getNote)
	ApiRouter.Post("/notes/{note}", postNote)
	ApiRouter.Get("/today", func(w http.ResponseWriter, r *http.Request) {
		getFile(files.DataFile("day/"+config.Cfg.TodayDate()), w)
	})
	ApiRouter.Post("/today", func(w http.ResponseWriter, r *http.Request) {
		postFile(files.DataFile("day/"+config.Cfg.TodayDate()), w, r)
	})
	ApiRouter.Get("/export", files.GetExport)
	ApiRouter.Get("/grace", graceStatus)
	ApiRouter.Get("/version", getVersion)
	ApiRouter.Get("/reload", configReload)
}
