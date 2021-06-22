package pterodactyl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gamemann/Pterodactyl-Game-Server-Watch/config"
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
	// Build endpoint.
	urlstr := cfg.APIURL + "/" + "api/client"

	// Setup HTTP GET request.
	client := &http.Client{Timeout: time.Second * 5}
	req, _ := http.NewRequest("GET", urlstr, nil)

	// Set authorization header.
	req.Header.Set("Authorization", "Bearer "+cfg.Token)

	// Accept only JSON.
	req.Header.Set("Accept", "application/json")

	// Perform HTTP request and check for errors.
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)

		return false
	}

	// Close body at the end.
	defer resp.Body.Close()

	// Read body.
	body, err := ioutil.ReadAll(resp.Body)

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

	// Loop through each data item (server).
	for _, j := range dataobj.(map[string]interface{})["data"].([]interface{}) {
		item := j.(map[string]interface{})

		// Make sure we have a server object.
		if item["object"] == "server" {
			attr := item["attributes"].(map[string]interface{})

			// Build new server structure.
			var sta config.Server

			// Set UID (in this case, identifier) and default values.
			sta.Enable = true
			sta.ViaAPI = true
			sta.ReportOnly = false
			sta.UID = attr["identifier"].(string)
			sta.ScanTime = cfg.DefScanTime
			sta.MaxFails = cfg.DefMaxFails
			sta.MaxRestarts = cfg.DefMaxRestarts
			sta.RestartInt = cfg.DefRestartInt
			sta.ReportOnly = cfg.DefReportOnly

			// Retrieve default IP/port.
			for _, i := range attr["relationships"].(map[string]interface{})["allocations"].(map[string]interface{})["data"].([]interface{}) {
				if i.(map[string]interface{})["object"].(string) != "allocation" {
					continue
				}

				alloc := i.(map[string]interface{})["attributes"].(map[string]interface{})

				if alloc["is_default"].(bool) {
					sta.IP = alloc["ip"].(string)
					sta.Port = int(alloc["port"].(float64))
				}
			}

			// Look for overrides.
			if attr["relationships"].(map[string]interface{})["variables"].(map[string]interface{})["data"] != nil {
				for _, i := range attr["relationships"].(map[string]interface{})["variables"].(map[string]interface{})["data"].([]interface{}) {
					if i.(map[string]interface{})["object"].(string) != "egg_variable" {
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

					// Check for report only override.
					if vari["env_variable"].(string) == "PTEROWATCH_REPORTONLY" {
						reportonly, _ := strconv.Atoi(val)

						if reportonly > 0 {
							sta.ReportOnly = true
						}
					}

					// Check for disable override.
					if vari["env_variable"].(string) == "PTEROWATCH_DISABLE" {
						disable, _ := strconv.Atoi(val)

						if disable > 0 {
							sta.Enable = false
						}
					}
				}
			}

			// Append to servers slice.
			cfg.Servers = append(cfg.Servers, sta)
		}
	}

	return true
}

// Checks the status of a Pterodactyl server. Returns true if on and false if off.
// DOES NOT INCLUDE IN "STARTING" MODE.
func CheckStatus(cfg *config.Config, uid string) bool {
	// Build endpoint.
	urlstr := cfg.APIURL + "/" + "api/client/servers/" + uid + "/resources"

	// Setup HTTP GET request.
	client := &http.Client{Timeout: time.Second * 5}
	req, _ := http.NewRequest("GET", urlstr, nil)

	// Set authorization header.
	req.Header.Set("Authorization", "Bearer "+cfg.Token)

	// Set data to JSON.
	req.Header.Set("Content-Type", "application/json")

	// Accept JSON.
	req.Header.Set("Accept", "application/json")

	// Perform HTTP request and check for errors.
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)

		return false
	}

	// Close body at the end.
	defer resp.Body.Close()

	// Read body.
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)

		return false
	}

	// Create utilization struct.
	var util Utilization

	// Parse JSON.
	json.Unmarshal([]byte(string(body)), &util)

	// Check if the server's state isn't on. If not, return false.
	if util.Attributes.State != "running" {
		return false
	}

	// Otherwise, return true meaning the container is online.
	return true
}

// Kills the specified server.
func KillServer(cfg *config.Config, uid string) {
	// Build endpoint.
	urlstr := cfg.APIURL + "/" + "api/client/servers/" + uid + "/" + "power"

	// Setup form data.
	var formdata = []byte(`{"signal": "kill"}`)

	// Setup HTTP GET request.
	client := &http.Client{Timeout: time.Second * 5}
	req, _ := http.NewRequest("POST", urlstr, bytes.NewBuffer(formdata))

	// Set authorization header.
	req.Header.Set("Authorization", "Bearer "+cfg.Token)

	// Set data to JSON.
	req.Header.Set("Content-Type", "application/json")

	// Accept JSON.
	req.Header.Set("Accept", "application/json")

	// Perform HTTP request and check for errors.
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
	}

	// Close body at the end.
	resp.Body.Close()
}

// Starts the specified server.
func StartServer(cfg *config.Config, uid string) {
	// Build endpoint.
	urlstr := cfg.APIURL + "/" + "api/client/servers/" + uid + "/" + "power"

	// Setup form data.
	var formdata = []byte(`{"signal": "start"}`)

	// Setup HTTP GET request.
	client := &http.Client{Timeout: time.Second * 5}
	req, _ := http.NewRequest("POST", urlstr, bytes.NewBuffer(formdata))

	// Set authorization header.
	req.Header.Set("Authorization", "Bearer "+cfg.Token)

	// Set data to JSON.
	req.Header.Set("Content-Type", "application/json")

	// Accept JSON.
	req.Header.Set("Accept", "application/json")

	// Perform HTTP request and check for errors.
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
	}

	// Close body at the end.
	resp.Body.Close()
}
