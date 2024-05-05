# 🌺 Hibiscus.txt

Simple plaintext diary.

This project is *very* opinionated and minimal, and is designed primarily for my usage. 
As a result, I can't guarantee that it's either secure or stable.

## Features:
* Each day, you get a new text file. You have until the end of that very day to finalise it.
* You can save named notes to document milestones, big events, or just a nice game you played this month
* You can easily export entire `data` dir in a `.zip` archive for backups

* Everything is plain(text) and simple. No databases, encryption, OAuth, or anything fancy. Even the password is plainte- *wait is this a feature?*
* Docker support (in fact, that's probably the best way to run this)
* Optional Telegram notifications for failed login attempts

## Technical details
You can read a relevant entry in my blog [here](https://a71.su/notes/hibiscus/). 
It provides some useful information and context for why this app exists in the first place.
There is also a [TODO.md](./TODO.md) file that shows what I will (or *may*) change in the future.
This repository is [mirrored to GitHub](https://github.com/Andrew-71/hibiscus) in case my server goes down.

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
Below are defaults of config.txt. The settings are defined in newline separated key=value pairs.
Please don't include the bash-style "comments" in your actual config, 
they are provided purely for demonstration only and **will break the config if present**.
```
username=admin  # Your username
password=admin  # Your password
port=7101  # What port to run on (probably leave on 7101 if using docker)
timezone=Local  # IANA Time zone database identifier ("UTC", Local", "Europe/Moscow" etc.), Defaults to Local if can't parse.
grace_period=0s  # Time after a new day begins, but before the switch to next day's file. e.g. 2h30m - files will change at 2:30am
language=en  # ISO-639 language code (currently supported - en, ru)
log_to_file=false  # Whether to write logs to a file
log_file=config/log.txt  # Where to store the log file if it is enabled
enable_scram=false  # Whether the app should shut down if there are 3 or more failed login attempts within 100 seconds

# Not present by default, set only if you want to be notified of any failed login attempts over telegram
tg_token=tgtoken
tg_chat=chatid
```

### Docker deployment:
Due to project's simplicity ~~and me rarely using them~~ there are no image tags, I just use `latest` and push to it.
The [package](https://git.a71.su/Andrew71/hibiscus/packages) provided in this repository is for `linux/amd64` architecture,
and there is a [Dockerfile](./Dockerfile) in case you want to compile for something rarer (like a Pi).
This repo contains the [compose.yml](./compose.yml) that I personally use.

### Executable flags
If you for some reason decide to run plain executable instead of docker, it supports following flags:
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
You can access the API at `/api/<method>`. They are protected by same HTTP Basic Auth as "normal" site.
```
GET  /today        - get file contents for today
POST /today        - save request body into today's file
GET  /day          - get JSON list of all daily entries
GET  /day/<name>   - get file contents for a specific day

GET  /notes        - get JSON list of all named notes
GET  /notes/<name> - get file contents for a specific note
POST /notes/<name> - save request body into a named note

GET  /export       - get .zip archive of entire data directory

GET  /readme       - get file contents for readme.txt in data dir's root
POST /readme       - save request body into readme.txt
```