package app

import (
	"git.a71.su/Andrew71/hibiscus-txt/internal/logging"
	"git.a71.su/Andrew71/hibiscus-txt/internal/server"
)

func Execute() {
	FlagInit()
	logging.LogInit()
	server.Serve()
}
