package routes

import (
	"net/http"

	"git.a71.su/Andrew71/hibiscus-txt/internal/config"
	"git.a71.su/Andrew71/hibiscus-txt/internal/files"
	"git.a71.su/Andrew71/hibiscus-txt/internal/lang"
	"git.a71.su/Andrew71/hibiscus-txt/internal/server/auth"
	"github.com/go-chi/chi/v5"
)

var UserRouter *chi.Mux

func init() {
	UserRouter = chi.NewRouter()
	UserRouter.Use(auth.BasicAuth)
	UserRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		getEntry(w, r, lang.Translate("title.today"), files.DataFile("day/"+config.Cfg.TodayDate()), true)
	})
	UserRouter.Post("/", func(w http.ResponseWriter, r *http.Request) {
		postEntry(files.DataFile("day/"+config.Cfg.TodayDate()), w, r)
	})
	UserRouter.Get("/day", getDays)
	UserRouter.Get("/day/{day}", getDay)
	UserRouter.Get("/notes", getNotes)
	UserRouter.Get("/notes/{note}", getNote)
	UserRouter.Post("/notes/{note}", postNote)
	UserRouter.Get("/info", getInfo)
	UserRouter.Get("/readme", func(w http.ResponseWriter, r *http.Request) {
		getEntry(w, r, "readme.txt", files.DataFile("readme"), true)
	})
	UserRouter.Post("/readme", func(w http.ResponseWriter, r *http.Request) { postEntry(files.DataFile("readme"), w, r) })
	UserRouter.Get("/config", func(w http.ResponseWriter, r *http.Request) {
		getEntry(w, r, "config.txt", config.ConfigFile, true)
	})
	UserRouter.Post("/config", postConfig)
}
