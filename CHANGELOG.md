# Changelog
This file keeps track of changes in a human-readable fashion

## v1.1.4
* Fixed HTML `lang` tag
* Theme CSS link is now only present if non-default is set
* Improved template consistency (backend)

## v1.1.3
This release mostly consists of backend improvements
* List items no longer replace hyphens with spaces for consistency
* Telegram message for SCRAM is now translatable
* Ensured HTML escape in list descriptions
* Refactored many methods, improved comments
## v1.1.2
This release contains a few bug fixes
* Real IPs are now logged (By Rithas K.)
* CSS now has `box-sizing: border-box` to fix textarea in some cases (By Rithas K.)
* Done some minor code housekeeping 
## v1.1.1
This release is mostly a technicality, with a move over to GitHub (`ghcr.io/andrew-71/hibiscus`) for packages due to DockerHub's prior anti-Russian actions making old "CI/CD" unsustainable.
## v1.1.0
* You can now specify the Telegram *topic* to send notification to via `tg_topic` config key (By Rithas K.)
* The Telegram message is now partially translated
* Fixed CSS `margin` and `text-align` inherited from my website

## v1.0.0
This release includes several **breaking** changes
* Made a new favicon
* English is now used as a fallback language, making incomplete translations somewhat usable
* Added a link to the bottom of the page in day list, for when you need to get to footer but been using the app for months...
* `pages`, `public` and `i18n` directories now use embed.FS
  * Running plain executable is now a somewhat valid option.
  * A problem with this is that languages and themes are now harder to add. I will think about what to do about that, maybe some kind of `custom.css` file.
  I might also be open to GitHub pull requests adding **some** languages (German and French could be nice starting points, I have friends studying them)
* Added a new "sans" theme
  * Light blue accent colour
  * Comic Sans MS for *everything*
  * sorry

## v0.6.0
* Replaced config reload with edit in info (api method still available, config reloads on save)
* Bug fixes
  * Filenames are now sanitized when writing files
  * "Tomorrow" in days list is now also displayed if Timezone is changed and grace period is off
  * Frontend date display now uses configured timezone
### v0.6.1
* Fixed date display when using `Local` timezone

## v0.5.0
* Added a JS prompt to create new note
  * "Sanitization" for this method is basic and assumes a well-meaning user
  * Old instructions appear if JS is disabled
* Bug fixes
  * Non-latin note names are now rendered correctly
  * Config reload now sets removed values to defaults

## v0.4.0
* Customisation changes
  * Added `title` option to config
    * Controls the text in the header, "🌺 Hibiscus.txt" by default
  * Added a nice `lavender` theme :)
  * No longer ensuring config.Theme ends up inside `/public`, unsure what to do in that regard
* Technical changes to config
  * Now only *some* default values are saved to file when creating initial config.txt
  * Spaces in config options are now supported (basically just for `title`)

## v0.3.0
* Added themes
  * Picked theme is set by `theme` key in config. Default is ...`default`
  * Themes are defined in `/public/themes/<name>.css` and modify colours (or, theoretically, do more)
  * Current pre-made themes are `default`, `gruvbox` and `high-contrast`

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