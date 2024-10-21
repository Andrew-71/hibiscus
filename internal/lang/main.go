package lang

import (
	"embed"
	"encoding/json"
	"log/slog"
)

//go:embed lang
var lang embed.FS
var translations = map[string]string{}

// SetLanguage loads a json file for selected language into the Translations map, with English language as a fallback.
func SetLanguage(language string) error {
	loadLanguage := func(language string) error {
		filename := "lang/" + language + ".json"
		fileContents, err := lang.ReadFile(filename)
		if err != nil {
			slog.Error("error reading language file",
				"error", err,
				"file", filename)
			return err
		}
		return json.Unmarshal(fileContents, &translations)
	}
	translations = map[string]string{} // Clear the map to avoid previous language remaining
	err := loadLanguage("en")          // Load English as fallback
	if err != nil {
		return err
	}
	return loadLanguage(language)
}

// Translate attempts to match an id to a string in current language.
func Translate(id string) string {
	if v, ok := translations[id]; !ok {
		return id
	} else {
		return v
	}
}
