package main

import (
	"embed"
	"encoding/json"
	"log/slog"
)

//go:embed i18n
var I18n embed.FS
var Translations = map[string]string{}

// SetLanguage loads a json file for selected language into the Translations map, with English language as a fallback.
func SetLanguage(language string) error {
	loadLanguage := func(language string) error {
		filename := "i18n/" + language + ".json"
		fileContents, err := I18n.ReadFile(filename)
		if err != nil {
			slog.Error("error reading language file",
				"error", err,
				"file", filename)
			return err
		}
		return json.Unmarshal(fileContents, &Translations)
	}
	Translations = map[string]string{} // Clear the map to avoid previous language remaining
	err := loadLanguage("en")          // Load English as fallback
	if err != nil {
		return err
	}
	return loadLanguage(language)
}

// TranslatableText attempts to match an id to a string in current language.
func TranslatableText(id string) string {
	if v, ok := Translations[id]; !ok {
		return id
	} else {
		return v
	}
}
