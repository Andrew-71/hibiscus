# Changelog
This file keeps track of changes in more human-readable fashion

## v0.2.0
* Added config reload
  * Can be reloaded in info page
  * Can be reloaded with new `reload` api method (be aware of the redirect if referer is present)

## 7 May 2024 - v0.1.0
* Began move towards [semantic versioning](https://semver.org/)
  * Current version is now 0.1.0
  * Added `version` api method
  * Versioned docker images may be provided in the future
  * Added version to footer
* Added info page
  * Accessed by clicking version number in footer
  * Allows you to edit readme.txt
  * Provides UI link for `export` api method
  * Can be expanded with other functionality in the future (see [TODO](./TODO.md))
* Bug fixes
  * Fixed export function failing

## 6 May 2024
* Grace period is now non-inclusive (so `4h` means the switch will happen right at `4:00`, not `4:01`)
* Added API method to check if grace period is active
* Made changes to date display on frontend
  * The date is now updated every minute, instead of every hour.
  * Now using the API to dynamically update the grace indicator
  * I kinda dislike this change, since it complicates the structure a *bit*.
  But I think it's fine.

## 5 May 2024
* Added this changelog
* Added grace period (as per suggestions)
  * Set in config like `grace_period=3h26m` (via Go's `time.ParseDuration`)
  * Defines how long the app should wait before switching to next day.
  The example provided means you will still be editing yesterday's file at 3am, but at 3:27am it will be a new day
  * When in effect, the header will show "(grace period active)" next to actual date
  * Fun fact: if you suddenly increase the setting mid-day, you can have a "Tomorrow" in previous days list! :)
  * This feature came with some free minor bug-fixes 'cause I had to re-check time management code.
  Now let's hope it doesn't introduce any new ones! :D
* Began adding PWA manifest.json
  * This will allow for (slightly) better mobile experience
  * Still missing an icon, will be likely installable once I make one and test the app :D
  * Known issue: making notes is impossible in the PWA, since you can't navigate to arbitrary page. 
  I might leave it as a WONTFIX or try to find a workaround
* Date is now shown in local language if possible (in case you add your own or use Russian)
* Added API reference to README
* Bug fixes
  * "Today" redirect from days list no longer uses UTC
  * Date JS script no longer uses UTC
  * The API no longer uses UTC for today
  * `/public/` files is no longer behind auth
  * Removed Sigmund® Corp. integration :)