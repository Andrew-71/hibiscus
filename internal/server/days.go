package server

import (
	"html/template"
	"net/http"
	"strings"
	"time"

	"git.a71.su/Andrew71/hibiscus-txt/internal/config"
	"git.a71.su/Andrew71/hibiscus-txt/internal/files"
	"git.a71.su/Andrew71/hibiscus-txt/internal/lang"
	"github.com/go-chi/chi/v5"
)

// GetDays calls GetEntries for previous days' entries.
func GetDays(w http.ResponseWriter, r *http.Request) {
	description := template.HTML(
		"<a href=\"#footer\">" + template.HTMLEscapeString(lang.Translate("prompt.days")) + "</a>")
	GetEntries(w, r, lang.Translate("title.days"), description, "day", func(files []string) []Entry {
		var filesFormatted []Entry
		for i := range files {
			v := files[len(files)-1-i] // This is suboptimal, but reverse order is better here
			dayString := v
			t, err := time.Parse(time.DateOnly, v)
			if err == nil {
				dayString = t.Format("02 Jan 2006")
			}

			// Fancy text for today and tomorrow
			// This looks bad, but strings.Title is deprecated, and I'm not importing a golang.org/x package for this...
			// (chances we ever run into tomorrow are really low)
			if v == config.Cfg.TodayDate() {
				dayString = lang.Translate("link.today")
				dayString = strings.ToTitle(string([]rune(dayString)[0])) + string([]rune(dayString)[1:])
			} else if v > config.Cfg.TodayDate() {
				dayString = lang.Translate("link.tomorrow")
				dayString = strings.ToTitle(string([]rune(dayString)[0])) + string([]rune(dayString)[1:])
			}
			filesFormatted = append(filesFormatted, Entry{Title: dayString, Link: "day/" + v})
		}
		return filesFormatted
	})
}

// GetDay calls GetEntry for a day entry.
func GetDay(w http.ResponseWriter, r *http.Request) {
	dayString := chi.URLParam(r, "day")
	if dayString == "" {
		w.WriteHeader(http.StatusBadRequest)
		HandleWrite(w.Write([]byte("day not specified")))
		return
	}
	if dayString == config.Cfg.TodayDate() { // Today can still be edited
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	title := dayString
	t, err := time.Parse(time.DateOnly, dayString)
	if err == nil { // This is low priority so silently fail
		title = t.Format("02 Jan 2006")
	}

	GetEntry(w, r, title, files.DataFile("day/"+dayString), false)
}
