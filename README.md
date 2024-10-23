# 🌺 Hibiscus.txt

Simple plaintext diary.

This project is *very* opinionated and minimal, and is designed primarily for my usage. 
As a result, neither security nor stability are guaranteed.

## Features:

* Each day, you get a new text file. You have until the end of that very day to finalise it.
* You can save named notes to document milestones, big events, or just a nice game you played this month
* You can easily export the files in a `.zip` archive for backups

* Everything is plain(text) and simple.
No databases, encryption, OAuth, or anything fancy. Even the password is plainte- *wait is this a feature?*
* [Docker support](#docker-deployment)
* Optional Telegram notifications for failed login attempts

## Technical details

[CHANGELOG.md](./CHANGELOG.md) provides a good overview of updates, and [TODO.md](./TODO.md) file shows my plans for the future.

You can read a relevant entry in my blog [here](https://a71.su/notes/hibiscus/).
It provides some useful information and context for why this app exists in the first place.
This repository is [self-hosted by me](https://git.a71.su/Andrew71/hibiscus),
but [mirrored to GitHub](https://github.com/Andrew-71/hibiscus).

### Data format:

```
data
+-- day
|   +-- yyyy-mm-dd.txt (ISO 8601)
|       ...
+-- notes
|   +-- note1.txt
|   +-- note2.txt
|       ...
+-- readme.txt

config
+-- config.txt
```
Deleting notes is done by clearing contents and clicking "Save" - the app deletes empty files when saving.

### Config options:

Below are the available configuration options and their defaults. 
The settings are defined as newline separated `key=value` pairs in the config file.
If you do not provide an option, the default will be used.
Please don't include the bash-style "comments" in your actual config, 
they are provided purely for demonstration and **will break the config if present**.
```
username=admin  # Your username
password=admin  # Your password
port=7101  # What port to run on (probably leave on 7101 if using docker)
timezone=Local  # IANA Time zone database identifier ("UTC", Local", "Europe/Moscow" etc.), Defaults to Local if can't parse.
grace_period=0s  # Time after a new day begins, but before the switch to next day's file. e.g. 3h26m - files will change at 3:26am
language=en  # ISO-639 language code (available - en, ru)
theme=""  # Picked theme (available - default (if left empty), high-contrast, lavender, gruvbox, sans)
title=🌺 Hibiscus.txt  # The text in the header
log_to_file=false  # Whether to write logs to a file
log_file=config/log.txt  # Where to store the log file if it is enabled
enable_scram=false  # Whether the app should shut down if there are 3 or more failed login attempts within 100 seconds

# Not present by default, set only if you want to be notified of any failed login attempts over Telegram
# Values correspond to API flags, see https://core.telegram.org/bots/api#sendmessage
tg_token=your_telegram_token
tg_chat=chat_id
tg_topic=message_thread_id
```

### Docker deployment:

The Docker images are hosted via GitHub over at `ghcr.io/andrew-71/hibiscus:<tag>`, 
built from the [Dockerfile](./Dockerfile).
This repo contains the [compose.yml](./compose.yml) that I personally use.
*Note: an extremely outdated self-hosted [package](https://git.a71.su/Andrew71/hibiscus/packages) will be provided for some time.*

### Executable flags

If you decide to use plain executable instead of docker, it supports the following flags:
```
-config string
    override config file location
-user string
    override username
-pass string
    override password
-port int
    override port
-debug
    show debug log
```

### API methods

You can access the API at `/api/<method>`. It is protected by same HTTP Basic Auth as "normal" routes.
```
GET  /today        - get file contents for today
POST /today        - save request body into today's file
GET  /day          - get JSON list of all daily entries
GET  /day/<name>   - get file contents for a specific day

GET  /notes        - get JSON list of all named notes
GET  /notes/<name> - get file contents for a specific note
POST /notes/<name> - save request body into a named note
GET  /readme       - get file contents for readme.txt in data dir's root
POST /readme       - save request body into readme.txt

GET  /export       - get .zip archive of entire data directory
GET  /grace        - "true" if grace period is active, otherwise "false"
GET  /version      - get app's version
GET  /reload       - reload app config
```