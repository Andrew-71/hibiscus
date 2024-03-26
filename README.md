# 🌺 Hibiscus.txt

Simple plaintext diary.

This project is *very* opinionated and minimal, and is designed primarily for my usage. 
As a result of this, it is also neither secure nor stable.

## Features:
* Each day, you get a text file. You have until 23:59 of that very day to finalise it.
* At any moment, you can append a single line to log.txt
* You can save named notes to document milestones, big events, or just a nice game you played this month
* There is also a readme.txt file (just like this one, except you get to write it!)
* You can easily export everything in a zip file for backups

* Everything is plain(text) and simple. No databases, encryption, OAuth, or anything fancy. Even the password is plainte- *wait is this a feature?*
* Docker support (in fact, that's probably the best way to run this)
* Optional Telegram notifications for failed login attempts

## Data format:
```
data
+-- day
|   +-- yyyy-mm-dd.txt (ISO 8601)
|       ...
+-- notes
|   +-- note1.txt
|   +-- note2.txt
|       ...
+-- log.txt
+-- readme.txt

config
+-- config.txt
```