package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var ConfigFile = "config/config.txt"

type Config struct {
	Username  string         `config:"username" type:"string" mandatory:"true"`
	Password  string         `config:"password" type:"string" mandatory:"true"`
	Port      int            `config:"port" type:"int" mandatory:"true"`
	Timezone  *time.Location `config:"timezone" type:"location" mandatory:"true"`
	GraceTime time.Duration  `config:"grace_period" type:"duration"`
	Language  string         `config:"language" type:"string" mandatory:"true"`
	Theme     string         `config:"theme" type:"string"`
	Title     string         `config:"title" type:"string"`
	LogToFile bool           `config:"log_to_file" type:"bool"`
	LogFile   string         `config:"log_file" type:"string"`
	Scram     bool           `config:"enable_scram" type:"bool"`

	TelegramToken string `config:"tg_token" type:"string"`
	TelegramChat  string `config:"tg_chat" type:"string"`
	TelegramTopic string `config:"tg_topic" type:"string"`
}

var DefaultConfig = Config{
	Username:  "admin",
	Password:  "admin",
	Port:      7101,
	Timezone:  time.Local,
	GraceTime: 0,
	Language:  "en",
	Theme:     "",
	Title:     "🌺 Hibiscus.txt",
	LogToFile: false,
	LogFile:   "config/log.txt",
	Scram:     false,

	TelegramToken: "",
	TelegramChat:  "",
	TelegramTopic: "",
}

// String returns string representation of modified and mandatory config options.
func (c *Config) String() string {
	output := ""
	v := reflect.ValueOf(*c)
	vDefault := reflect.ValueOf(DefaultConfig)
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		key := typeOfS.Field(i).Tag.Get("config")
		value := v.Field(i).Interface()
		mandatory := typeOfS.Field(i).Tag.Get("mandatory")
		if (mandatory == "true") || (value != vDefault.Field(i).Interface()) { // Only save non-default values
			output += fmt.Sprintf("%s=%v\n", key, value)
		}
	}
	return output
}

// Reload resets, then loads config from the ConfigFile.
// It creates the file with mandatory options if it is missing.
func (c *Config) Reload() error {
	*c = DefaultConfig // Reset config

	if _, err := os.Stat(ConfigFile); errors.Is(err, os.ErrNotExist) {
		err := c.Save()
		if err != nil {
			return err
		}
		return nil
	}
	file, err := os.Open(ConfigFile)
	if err != nil {
		return err
	}

	options := map[string]string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		entry := strings.Split(strings.Trim(scanner.Text(), " \t"), "=")
		if len(entry) != 2 {
			continue
		}
		options[entry[0]] = entry[1]
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}

	timezone := "Local" // Timezone is handled separately because reflection
	refStruct := reflect.ValueOf(*c)
	refElem := reflect.ValueOf(&c).Elem()
	typeOfS := refStruct.Type()
	for i := 0; i < refStruct.NumField(); i++ {
		fieldElem := reflect.Indirect(refElem).Field(i)
		key := typeOfS.Field(i).Tag.Get("config")
		if v, ok := options[key]; ok && fieldElem.CanSet() {
			switch typeOfS.Field(i).Tag.Get("type") {
			case "int":
				{
					numVal, err := strconv.Atoi(v)
					if err == nil {
						fieldElem.SetInt(int64(numVal))
					}
				}
			case "bool":
				fieldElem.SetBool(v == "true")
			case "location":
				timezone = v
			case "duration":
				{
					numVal, err := time.ParseDuration(v)
					if err == nil {
						fieldElem.SetInt(int64(numVal))
					}
				}
			default:
				fieldElem.SetString(v)
			}
		}
	}
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		c.Timezone = time.Local
	} else {
		c.Timezone = loc
	}
	slog.Debug("reloaded config", "config", c)

	return SetLanguage(c.Language) // Load selected language
}

// Read gets raw contents from ConfigFile.
func (c *Config) Read() ([]byte, error) {
	return ReadFile(ConfigFile)
}

// Save writes config's contents to the ConfigFile.
func (c *Config) Save() error {
	return SaveFile(ConfigFile, []byte(c.String()))
}

// PostConfig calls PostEntry for config file, then reloads the config.
func PostConfig(w http.ResponseWriter, r *http.Request) {
	PostEntry(ConfigFile, w, r)
	err := Cfg.Reload()
	if err != nil {
		slog.Error("error reloading config", "error", err)
	}
}

// ConfigReloadApi reloads the config. It then redirects back if Referer field is present.
func ConfigReloadApi(w http.ResponseWriter, r *http.Request) {
	err := Cfg.Reload()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		HandleWrite(w.Write([]byte(err.Error())))
	}
	if r.Referer() != "" {
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// ConfigInit loads config on startup.
func ConfigInit() Config {
	cfg := Config{}
	err := cfg.Reload()
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}
