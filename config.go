package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var ConfigFile = "config/config.txt"

type Config struct {
	Username  string         `config:"username" type:"string"`
	Password  string         `config:"password" type:"string"`
	Port      int            `config:"port" type:"int"`
	Timezone  *time.Location `config:"timezone" type:"location"`
	GraceTime time.Duration  `config:"grace_period" type:"duration"`
	Language  string         `config:"language" type:"string"`
	LogToFile bool           `config:"log_to_file" type:"bool"`
	LogFile   string         `config:"log_file" type:"string"`
	Scram     bool           `config:"enable_scram" type:"bool"`

	TelegramToken string `config:"tg_token" type:"string"`
	TelegramChat  string `config:"tg_chat" type:"string"`
}

func (c *Config) Save() error {
	output := ""

	v := reflect.ValueOf(*c)
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		key := typeOfS.Field(i).Tag.Get("config")
		value := v.Field(i).Interface()
		if value != "" {
			output += fmt.Sprintf("%s=%v\n", key, value)
		}
	}

	f, err := os.OpenFile(ConfigFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	if _, err := f.Write([]byte(output)); err != nil {
		return err
	}
	return nil
}

func (c *Config) Reload() error {
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
		entry := strings.Split(strings.Replace(scanner.Text(), " ", "", -1), "=")
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
				{
					if v == "true" {
						fieldElem.SetBool(true)
					} else {
						fieldElem.SetBool(false)
					}
				}
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

	return LoadLanguage(c.Language) // Load selected language
}

// ConfigInit loads config on startup
// Some defaults are declared here
func ConfigInit() Config {
	cfg := Config{
		Port:      7101,
		Username:  "admin",
		Password:  "admin",
		Timezone:  time.Local,
		Language:  "en",
		LogFile:   "config/log.txt",
		GraceTime: 0,
	}
	err := cfg.Reload()
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}
