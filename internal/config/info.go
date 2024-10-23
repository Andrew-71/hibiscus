package config

type AppInfo struct {
	version    string
	source string
}

// Info contains app information.
var Info = AppInfo{
	version:    "2.0.0",
	source: "https://git.a71.su/Andrew71/hibiscus",
}

// Version returns the current app version
func (i AppInfo) Version() string {
	return i.version
}

// Source returns app's git repository
func (i AppInfo) Source() string {
	return i.source
}