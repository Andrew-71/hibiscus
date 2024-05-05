# Changelog
This file keeps track of changes in more human-readable fashion

# 5 May 2024
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
* Bug fixes
  * "Today" redirect from days list no longer uses UTC
  * Date JS script no longer uses UTC
  * The API no longer uses UTC for today
  * `/public/` files is no longer behind auth
  * Removed Sigmund® Corp. integration :)