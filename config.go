package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var ConfigFile = "config/config.txt"

type Config struct {
	Username string
	Password string
	Port     int
}

func (c *Config) Save() error {
	output := fmt.Sprintf("port=%d\nusername=%s\npassword=%s", c.Port, c.Username, c.Password)

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
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func ConfigInit() Config {
	cfg := Config{Port: 7101, Username: "admin", Password: "admin"} // Default values are declared here, I guess
	err := cfg.Reload()
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}
