package main

import (
	"flag"
	"log"
)

// FlagInit processes app flags
func FlagInit() {
	config := flag.String("config", "", "override config file")
	username := flag.String("user", "", "override username")
	password := flag.String("pass", "", "override password")
	port := flag.Int("port", 0, "override port")
	debug := flag.Bool("debug", false, "debug logging")

	flag.Parse()
	if *config != "" {
		ConfigFile = *config
		err := Cfg.Reload()
		if err != nil {
			log.Fatal(err)
		}
	}
	if *username != "" {
		Cfg.Username = *username
	}
	if *password != "" {
		Cfg.Password = *password
	}
	if *port != 0 {
		Cfg.Port = *port
	}
	DebugMode = *debug
}
