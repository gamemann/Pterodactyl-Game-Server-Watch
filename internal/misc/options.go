package misc

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"

	"github.com/gamemann/Pterodactyl-Game-Server-Watch/pkg/config"
)

func HandleMisc(cfg *config.Config, srv *config.Server, fails int, restarts int) {
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
				contentpre := "**SERVER DOWN**\n- **Name** => {NAME}\n- **IP** => {IP}:{PORT}\n- **Fail Count** => {FAILS}/{MAXFAILS}\n- **Restart Count** => {RESTARTS}/{MAXRESTARTS}\n\nScanning again in *{RESTARTINT}* seconds..."
				username := "Pterowatch"
				password := ""
				avatarurl := ""
				topic := ""
				allowedmentions := AllowMentions{
					Roles: false,
					Users: false,
				}
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

				// Look for password override.
				if v.Data.(map[string]interface{})["password"] != nil {
					password = v.Data.(map[string]interface{})["password"].(string)
				}

				// Look for avatar URL override.
				if v.Data.(map[string]interface{})["avatarurl"] != nil {
					avatarurl = v.Data.(map[string]interface{})["avatarurl"].(string)
				}

				// Look for Ntfy Topic override.
				if v.Data.(map[string]interface{})["topic"] != nil {
					topic = v.Data.(map[string]interface{})["topic"].(string)
				}

				// Look for allowed mentions override.
				if v.Data.(map[string]interface{})["mentions"] != nil {
					mentdata := v.Data.(map[string]interface{})["mentions"].(map[string]interface{})

					roles := false
					users := false

					if mentdata["roles"] != nil {
						roles = mentdata["roles"].(bool)
					}

					if mentdata["users"] != nil {
						users = mentdata["users"].(bool)
					}

					allowedmentions.Roles = roles
					allowedmentions.Users = users
				}

				// Handle mentions.
				mentionstr := ""

				if (allowedmentions.Roles || allowedmentions.Users) && len(srv.Mentions) > 0 {
					if cfg.DebugLevel > 1 {
						fmt.Println("[D2] Parsing mention data for " + srv.UID + " ( " + srv.Name + ").")
					}

					var mentdata interface{}

					err := json.Unmarshal([]byte(srv.Mentions), &mentdata)

					if cfg.DebugLevel > 3 {
						fmt.Println("[D4] Mention JSON => " + srv.Mentions + ".")
					}

					if err != nil {
						fmt.Println("[ERR] Failed to parse JSON mention data for server " + srv.UID + " (" + srv.Name + ").")
						fmt.Println(err)

						goto skipment
					}

					if mentdata.(map[string]interface{})["data"] == nil {
						fmt.Println("[ERR] Mentions string missing data list for " + srv.UID + " (" + srv.Name + ").")

						goto skipment
					}

					if _, ok := mentdata.(map[string]interface{})["data"].([]interface{}); !ok {
						fmt.Println("[ERR] Mentions string's data item not a list for " + srv.UID + " (" + srv.Name + ").")

						goto skipment
					}

					len := len(mentdata.(map[string]interface{})["data"].([]interface{}))

					// Loop through each item.
					for i, m := range mentdata.(map[string]interface{})["data"].([]interface{}) {
						item := m.(map[string]interface{})

						// Check to ensure we have both elements/items.
						if item["role"] == nil || item["id"] == nil {
							continue
						}

						// Check types.
						if _, ok := item["role"].(bool); !ok {
							fmt.Println("[ERR] Mentions string's role field is not a boolean for " + srv.UID + " (" + srv.Name + ").")

							continue
						}

						if _, ok := item["id"].(string); !ok {
							fmt.Println("[Err] Mentions string's id field is not a string for " + srv.UID + " (" + srv.Name + ").")

							continue
						}

						// For security, we want to parse the values as big integers (float64 is also too small for IDs).
						id := big.Int{}
						id.SetString(item["id"].(string), 10)

						// Check for role.
						if item["role"].(bool) && allowedmentions.Roles {
							mentionstr += "<@&" + id.String() + ">"
						}

						// Check for user.
						if !item["role"].(bool) && allowedmentions.Users {
							mentionstr += "<@" + id.String() + ">"
						}

						// Check to see if we need a comma.
						if i != (len - 1) {
							mentionstr += ", "
						}
					}

					if cfg.DebugLevel > 3 {
						fmt.Println("[D4] Mention string => " + mentionstr + " for " + srv.UID + " (" + srv.Name + ").")
					}
				}

			skipment:

				// Replace variables in strings.
				contents := contentpre
				FormatContents(app, &contents, fails, restarts, srv, mentionstr)

				// Level 3 debug.
				if cfg.DebugLevel > 2 {
					fmt.Println("[D3] Loaded web hook with App => " + app + ". URL => " + url + ". Contents => " + contents + ". Username => " + username + ". Avatar URL => " + avatarurl + ". Mentions => Roles:" + strconv.FormatBool(allowedmentions.Roles) + "; Users:" + strconv.FormatBool(allowedmentions.Users) + ". Topic:" + topic + ".")
				}

				// Submit web hook.
				if app == "ntfy" {
					NtfyWebHook(url, contents, topic, username, password)
				} else if app == "slack" {
					SlackWebHook(url, contents)
				} else {
					DiscordWebHook(url, contents, username, avatarurl, allowedmentions, srv)
				}
			}
		}
	}
}
