package pterodactyl

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gamemann/Pterodactyl-Game-Server-Watch/pkg/config"
	pteroapi "github.com/gamemann/Rust-Auto-Wipe/pkg/pterodactyl"
)

// Attributes struct from /api/client/servers/xxxx/resources.
type Attributes struct {
	State string `json:"current_state"`
}

// Utilization struct from /api/client/servers/xxxx/resources.
type Utilization struct {
	Attributes Attributes `json:"attributes"`
}

// Retrieves all servers/containers from Pterodactyl API and add them to the config.
func AddServers(cfg *config.Config) bool {
	// Retrieve max page count.
	pagecount := 1
	maxpages := 1
	total := 0
	done := false

	for done != true {
		body, _, err := pteroapi.SendAPIRequest(cfg.APIURL, cfg.AppToken, "GET", "application/servers?page="+strconv.Itoa(pagecount)+"&include=allocations,variables", nil)

		if err != nil {
			fmt.Println(err)

			return false
		}

		// Create data interface.
		var dataobj interface{}

		// Parse JSON.
		err = json.Unmarshal([]byte(string(body)), &dataobj)

		if err != nil {
			fmt.Println(err)

			return false
		}

		// Look for object item before anything.
		if dataobj.(map[string]interface{})["object"] == nil {
			fmt.Println("[ERR] 'object' item not found when listing all servers.")

			fmt.Println(string(body))

			return false
		}

		// Retrieve max page count and total count.
		maxpages = int(dataobj.(map[string]interface{})["meta"].(map[string]interface{})["pagination"].(map[string]interface{})["total_pages"].(float64))
		total = int(dataobj.(map[string]interface{})["meta"].(map[string]interface{})["pagination"].(map[string]interface{})["total"].(float64))

		// Loop through each data item (server).
		for _, j := range dataobj.(map[string]interface{})["data"].([]interface{}) {
			item := j.(map[string]interface{})

			// Make sure we have a server object.
			if item["object"] == "server" {
				attr := item["attributes"].(map[string]interface{})

				// Build new server structure.
				var sta config.Server

				// Set UID (in this case, identifier) and default values.
				sta.ViaAPI = true
				sta.UID = attr["identifier"].(string)
				sta.Name = attr["name"].(string)

				sta.Enable = cfg.DefEnable
				sta.ScanTime = cfg.DefScanTime
				sta.MaxFails = cfg.DefMaxFails
				sta.MaxRestarts = cfg.DefMaxRestarts
				sta.RestartInt = cfg.DefRestartInt
				sta.ReportOnly = cfg.DefReportOnly
				sta.A2STimeout = cfg.DefA2STimeout
				sta.RconPassword = cfg.DefRconPassword
				sta.Mentions = cfg.DefMentions

				if attr["relationships"] == nil {
					fmt.Println("[ERR] Server has invalid relationships.")

					continue
				}

				// Retrieve default IP/port.
				for _, i := range attr["relationships"].(map[string]interface{})["allocations"].(map[string]interface{})["data"].([]interface{}) {
					if i.(map[string]interface{})["object"].(string) != "allocation" {
						continue
					}

					alloc := i.(map[string]interface{})["attributes"].(map[string]interface{})

					if alloc["assigned"].(bool) {
						sta.IP = alloc["ip"].(string)
						sta.Port = int(alloc["port"].(float64))
					}
				}

				// Look for overrides.
				if attr["relationships"].(map[string]interface{})["variables"].(map[string]interface{})["data"] != nil {
					for _, i := range attr["relationships"].(map[string]interface{})["variables"].(map[string]interface{})["data"].([]interface{}) {
						if i.(map[string]interface{})["object"].(string) != "server_variable" {
							continue
						}

						vari := i.(map[string]interface{})["attributes"].(map[string]interface{})

						// Check if we have a value.
						if vari["server_value"] == nil {
							continue
						}

						val := vari["server_value"].(string)

						// Override variables should always be at least one byte in length.
						if len(val) < 1 {
							continue
						}

						// Check for IP override.
						if vari["env_variable"].(string) == "PTEROWATCH_IP" {
							sta.IP = val
						}

						// Check for port override.
						if vari["env_variable"].(string) == "PTEROWATCH_PORT" {
							sta.Port, _ = strconv.Atoi(val)
						}

						// Check for scan override.
						if vari["env_variable"].(string) == "PTEROWATCH_SCANTIME" {
							sta.ScanTime, _ = strconv.Atoi(val)
						}

						// Check for max fails override.
						if vari["env_variable"].(string) == "PTEROWATCH_MAXFAILS" {
							sta.MaxFails, _ = strconv.Atoi(val)
						}

						// Check for max restarts override.
						if vari["env_variable"].(string) == "PTEROWATCH_MAXRESTARTS" {
							sta.MaxRestarts, _ = strconv.Atoi(val)
						}

						// Check for restart interval override.
						if vari["env_variable"].(string) == "PTEROWATCH_RESTARTINT" {
							sta.RestartInt, _ = strconv.Atoi(val)
						}

						// Check for A2S_INFO timeout override.
						if vari["env_variable"].(string) == "PTEROWATCH_A2STIMEOUT" {
							sta.A2STimeout, _ = strconv.Atoi(val)
						}

						// Check for RCON_PASSWORD override.
						if vari["env_variable"].(string) == "PTEROWATCH_RCONPASSWORD" {
							sta.RconPassword = val
						}

						// Check for mentions override.
						if vari["env_variable"].(string) == "PTEROWATCH_MENTIONS" {
							sta.Mentions = val
						}

						// Check for report only override.
						if vari["env_variable"].(string) == "PTEROWATCH_REPORTONLY" {
							reportonly, _ := strconv.Atoi(val)

							if reportonly > 0 {
								sta.ReportOnly = true
							} else {
								sta.ReportOnly = false
							}
						}

						// Check for disable override.
						if vari["env_variable"].(string) == "PTEROWATCH_DISABLE" {
							disable, _ := strconv.Atoi(val)

							if disable > 0 {
								sta.Enable = false
							} else {
								sta.Enable = true
							}
						}
					}
				}

				// Append to servers slice.
				cfg.Servers = append(cfg.Servers, sta)
			}
		}

		// Check page count.
		if pagecount >= maxpages {
			done = true

			break
		}

		pagecount++
	}

	// Level 2 debug.
	if cfg.DebugLevel > 1 {
		fmt.Println("[D2] Found " + strconv.Itoa(total) + " servers from API (" + strconv.Itoa(maxpages) + " page(s)).")
	}

	return true
}

// Checks the status of a Pterodactyl server. Returns true if on and false if off.
// DOES NOT INCLUDE IN "STARTING" MODE.
func CheckStatus(cfg *config.Config, uid string) bool {
	body, _, err := pteroapi.SendAPIRequest(cfg.APIURL, cfg.AppToken, "GET", "client/servers/"+uid+"/resources", nil)

	if err != nil {
		fmt.Println(err)

		return false
	}

	// Create utilization struct.
	var util Utilization

	// Parse JSON.
	json.Unmarshal([]byte(string(body)), &util)

	if cfg.DebugLevel > 3 {
		fmt.Println("[D6] CheckStatus ", string(body))
	}

	// Check if the server's state isn't on. If not, return false.
	if util.Attributes.State != "running" {
		return false
	}

	// Otherwise, return true meaning the container is online.
	return true
}

// Kills the specified server.
func KillServer(cfg *config.Config, uid string) bool {
	form_data := make(map[string]interface{})
	form_data["signal"] = "kill"

	_, _, err := pteroapi.SendAPIRequest(cfg.APIURL, cfg.AppToken, "POST", "client/servers/"+uid+"/"+"power", form_data)

	if err != nil {
		fmt.Println(err)

		return false
	}

	return true
}

// Stops the specified server.
func StopServer(cfg *config.Config, uid string) bool {
	form_data := make(map[string]interface{})
	form_data["signal"] = "stop"

	_, _, err := pteroapi.SendAPIRequest(cfg.APIURL, cfg.AppToken, "POST", "client/servers/"+uid+"/"+"power", form_data)

	if err != nil {
		fmt.Println(err)

		return false
	}

	return true
}

// Starts the specified server.
func StartServer(cfg *config.Config, uid string) bool {
	form_data := make(map[string]interface{})
	form_data["signal"] = "start"

	_, _, err := pteroapi.SendAPIRequest(cfg.APIURL, cfg.AppToken, "POST", "client/servers/"+uid+"/"+"power", form_data)

	if err != nil {
		fmt.Println(err)

		return false
	}

	return true
}
