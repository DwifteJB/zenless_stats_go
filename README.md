# Zenless Discord Widget!!! (stats!!!)

WARNING!! this uses a new discord preview feature, i havent seen anyone get banned but this is "unofficial" at the moment of writing, so just beware. also your hoyolab cookies could be compromised by any  viruses, but they will probably only look at your browser anyway!

## how to use!

* first, follow this [guide](https://gist.github.com/aamiaa/7cdd590e3949cd654758bc90bcb4710b) here to create your widget BUT!!! use the script [create-widget.js](./create-widget.js) to create the proper fields!! if you ignored that, then use [fix-widget.js](./fix-widget.js), be sure to note down the Application ID within "General Information" as well as the bot token within the CURL / Invoke script, or just generate a new one in Bot!
* secondly follow [cookie_extract.md](./cookie_extract.md) to get your hoyolab cookie 
* download or compile the app. on first run it will prompt for everything it needs: your zenless UID, region (a numbered selector), discord user id, widget bot token, widget bot id, hoyolab id & hoyolab cookie. it saves them to `config.json` next to the binary, so you only enter them once.
* thats it! run it with no arguments to update once, run `zenless_stats cron` to run as a service that updates every hour, or `zenless_stats cron "*/5 * * * *"` to pick your own schedule (standard 5 field cron).

every run writes what it fetched, the exact payload it sent, and the discord response code to `json/results.json`. if the widget isnt updating, check that file: `discord_status: 204` means discord accepted the update, so the problem is on the widget/display side rather than your credentials.

## widget not updating? bind the fields

by default `create-widget.js` builds the widget with static placeholder text (`value_type: "custom_string"`, "text 1 here"). those fields never look at the data this app pushes, so nothing changes. each widget field has to be set to `value_type: "data"` with `value` equal to the dynamic name below (in the dev portal editor: Value Type = "User Data", Data Field = the name).

| widget field | value name | label name |
| --- | --- | --- |
| title | `nickname` | - |
| stat 1 | `IL` | `IL_str` (Interknot Level) |
| stat 2 | `ach` | `ach_str` (Achievements) |
| stat 3 | `SBT` | `SBT_str` (Simulated Battle Trial) |
| stat 4 | `ENER` | `ENER_str` (Energy) |
| stat 5 | `days` | `days_str` (Days Active) |
| stat 6 | `polychromes` | `polychromes_str` (Monthly Polychromes) |

`create-widget.js` now sets these bindings for a fresh app. if your app already exists (you already have a token), dont re-run it (that makes a duplicate app) - instead put your application id into `fix-widget.js` and paste it into the discord dev portal console to rebind the existing widget.

if you need to find your hoyolab ID, you can find it by going to [https://www.hoyolab.com/accountCenter/postList](https://www.hoyolab.com/accountCenter/postList) & copying the ?id=[THIS NUMBER!!!]

## why?

this is cool, i found the code out from [pinkblossom3's repo](https://github.com/PinkBlossom3/Zenless-Stats) and i think the code quality could be a bit... better.
