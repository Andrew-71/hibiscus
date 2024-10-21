package config

import (
	"log/slog"
	"time"
)

// FIXME: This probably shouldn't be part of config package

// Grace returns whether the grace period is active.
// NOTE: Grace period has minute precision
func (c Config) Grace() bool {
	t := time.Now().In(c.Timezone)
	active := (60*t.Hour() + t.Minute()) < int(c.GraceTime.Minutes())
	if active {
		slog.Debug("grace period active",
			"time", 60*t.Hour()+t.Minute(),
			"grace", c.GraceTime.Minutes())
	}
	return active
}

// TodayDate returns today's formatted date. It accounts for Config.GraceTime.
func (c Config) TodayDate() string {
	dateFormatted := time.Now().In(c.Timezone).Format(time.DateOnly)
	if c.Grace() {
		dateFormatted = time.Now().In(c.Timezone).AddDate(0, 0, -1).Format(time.DateOnly)
	}
	slog.Debug("today", "time", time.Now().In(c.Timezone).Format(time.DateTime))
	return dateFormatted
}
