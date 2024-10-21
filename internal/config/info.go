package config

type AppInfo struct {
	Version    string
	SourceLink string
}

// Info contains app information.
var Info = AppInfo{
	Version:    "2.0.0",
	SourceLink: "https://git.a71.su/Andrew71/hibiscus",
}