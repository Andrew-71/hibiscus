package main

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Username string
	Password string
	Port     int
}

func CreateConfig(config Config) {

}

func LoadConfig() (Config, error) {
	filename := "config/config.txt"

	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		CreateConfig(Config{})
		return Config{}, err
	}

	cfg := Config{Port: 7101}

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
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
			cfg.Username = value
		} else if key == "password" {
			cfg.Password = value
		} else if key == "port" {
			numVal, err := strconv.Atoi(value)
			if err == nil {
				cfg.Port = numVal
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return cfg, nil
}
