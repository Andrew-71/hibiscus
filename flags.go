package main

import (
	"flag"
)

func FlagInit() {
	username := flag.String("user", "", "override username")
	password := flag.String("pass", "", "override password")
	port := flag.Int("port", 0, "override port")

	flag.Parse()
	if *username != "" {
		Cfg.Username = *username
	}
	if *password != "" {
		Cfg.Password = *password
	}
	if *port != 0 {
		Cfg.Port = *port
	}
}
