package app

import (
	"flag"
	"log"

	"git.a71.su/Andrew71/hibiscus-txt/internal/config"
	"git.a71.su/Andrew71/hibiscus-txt/internal/logging"
)

// FlagInit processes app flags.
func FlagInit() {
	conf := flag.String("config", "", "override config file")
	username := flag.String("user", "", "override username")
	password := flag.String("pass", "", "override password")
	port := flag.Int("port", 0, "override port")
	debug := flag.Bool("debug", false, "debug logging")

	flag.Parse()
	if *conf != "" {
		config.ConfigFile = *conf
		err := config.Cfg.Reload()
		if err != nil {
			log.Fatal(err)
		}
	}
	if *username != "" {
		config.Cfg.Username = *username
	}
	if *password != "" {
		config.Cfg.Password = *password
	}
	if *port != 0 {
		config.Cfg.Port = *port
	}
	logging.DebugMode = *debug
}
