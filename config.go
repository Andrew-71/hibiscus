package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var ConfigFile = "config/config.txt"

type Config struct {
	Username  string         `config:"username"`
	Password  string         `config:"password"`
	Port      int            `config:"port"`
	Timezone  *time.Location `config:"timezone"`
	LogToFile bool           `config:"log_to_file"`
	Scram     bool           `config:"enable_scram"`

	TelegramToken string `config:"tg_token"`
	TelegramChat  string `config:"tg_chat"`
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
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		entry := strings.Split(strings.Replace(scanner.Text(), " ", "", -1), "=")
		if len(entry) != 2 {
			continue
		}
		key := entry[0]
		value := entry[1]
		if key == "username" {
			c.Username = value
		} else if key == "password" {
			c.Password = value
		} else if key == "port" {
			numVal, err := strconv.Atoi(value)
			if err == nil {
				c.Port = numVal
			}
		} else if key == "timezone" {
			loc, err := time.LoadLocation(value)
			if err != nil {
				c.Timezone = time.UTC
			} else {
				c.Timezone = loc
			}
		} else if key == "tg_token" {
			c.TelegramToken = value
		} else if key == "tg_chat" {
			c.TelegramChat = value
		} else if key == "enable_scram" {
			if value == "true" {
				c.Scram = true
			} else {
				c.Scram = false
			}
		} else if key == "log_to_file" {
			if value == "true" {
				c.LogToFile = true
			} else {
				c.LogToFile = false
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

// ConfigInit loads config on startup
func ConfigInit() Config {
	cfg := Config{Port: 7101, Username: "admin", Password: "admin", Timezone: time.UTC} // Default values are declared here, I guess
	err := cfg.Reload()
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}
