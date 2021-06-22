package misc

import (
	"fmt"
	"strconv"

	"github.com/gamemann/Pterodactyl-Game-Server-Watch/config"
)

func HandleMisc(cfg *config.Config, srvidx int, fails int, restarts int) {
	srv := cfg.Servers[srvidx]

	// Look for Misc options.
	if len(cfg.Misc) > 0 {
		for i, v := range cfg.Misc {
			// Level 2 debug.
			if cfg.DebugLevel > 1 {
				fmt.Println("[D2] Loading misc option #" + strconv.Itoa(i) + " with type " + v.Type + ".")
			}

			// Handle web hooks.
			if v.Type == "webhook" {
				// Set defaults.
				contentpre := "**SERVER DOWN**\n**Name** => {NAME}\n- **IP** => {IP}:{PORT}\n- **Fail Count** => {FAILS}/{MAXFAILS}\n**Restart Count** => {RESTARTS}/{MAXRESTARTS}\n\nScanning again in *{RESTARTINT}* seconds..."
				username := "Pterowatch"
				avatarurl := ""
				allowedmentions := false
				app := "discord"

				// Look for app override.
				if v.Data.(map[string]interface{})["app"] != nil {
					app = v.Data.(map[string]interface{})["app"].(string)
				}

				// Check for webhook URL.
				if v.Data.(map[string]interface{})["url"] == nil {
					fmt.Println("[ERR] Web hook ID #" + strconv.Itoa(i) + " has no webhook URL.")

					continue
				}

				url := v.Data.(map[string]interface{})["url"].(string)

				// Look for contents override.
				if v.Data.(map[string]interface{})["contents"] != nil {
					contentpre = v.Data.(map[string]interface{})["contents"].(string)
				}

				// Look for username override.
				if v.Data.(map[string]interface{})["username"] != nil {
					username = v.Data.(map[string]interface{})["username"].(string)
				}

				// Look for avatar URL override.
				if v.Data.(map[string]interface{})["avatarurl"] != nil {
					avatarurl = v.Data.(map[string]interface{})["avatarurl"].(string)
				}

				// Look for allowed mentions override.
				if v.Data.(map[string]interface{})["allowedmentions"] != nil {
					allowedmentions = v.Data.(map[string]interface{})["avatarurl"].(bool)
				}

				// Replace variables in strings.
				contents := contentpre
				FormatContents(&contents, fails, restarts, &srv)

				// Level 3 debug.
				if cfg.DebugLevel > 2 {
					fmt.Println("[D3] Loaded web hook with App => " + app + ". URL => " + url + ". Contents => " + contents + ". Username => " + username + ". Avatar URL => " + avatarurl + ". Allowed Mentions => " + strconv.FormatBool(allowedmentions) + ".")
				}

				// Submit web hook.
				if app == "slack" {
					SlackWebHook(url, contents)
				} else {
					DiscordWebHook(url, contents, username, avatarurl, allowedmentions)
				}
			}
		}
	}
}
