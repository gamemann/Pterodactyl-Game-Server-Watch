package pterodactyl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Attributes struct from /api/client/servers/xxxx/utilization.
type Attributes struct {
	State string `json:"state"`
}

// Utilization struct from /api/client/servers/xxxx/utilization.
type Utilization struct {
	Attributes Attributes `json:"attributes"`
}

// Checks the status of a Pterodactyl server. Returns true if on and false if off.
func CheckStatus(apiURL string, apiToken string, uid string) bool {
	// Build endpoint.
	urlstr := apiURL + "/" + "api/client/servers/" + uid + "/" + "utilization"

	// Setup HTTP GET request.
	client := &http.Client{Timeout: time.Second * 5}
	req, _ := http.NewRequest("GET", urlstr, nil)

	// Set authorization header.
	req.Header.Set("Authorization", "Bearer "+apiToken)

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
	if util.Attributes.State != "on" {
		return false
	}

	// Otherwise, return true meaning the container is online.
	return true
}

// Kills the specified server.
func KillServer(apiURL string, apiToken string, uid string) {
	// Build endpoint.
	urlstr := apiURL + "/" + "api/client/servers/" + uid + "/" + "power"

	// Setup form data.
	var formdata = []byte(`{"signal": "kill"}`)

	// Setup HTTP GET request.
	client := &http.Client{Timeout: time.Second * 5}
	req, _ := http.NewRequest("POST", urlstr, bytes.NewBuffer(formdata))

	// Set authorization header.
	req.Header.Set("Authorization", "Bearer "+apiToken)

	// Set data type to JSON.
	req.Header.Set("Content-Type", "application/json")

	// Perform HTTP request and check for errors.
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
	}

	// Close body at the end.
	resp.Body.Close()
}

// Starts the specified server.
func StartServer(apiURL string, apiToken string, uid string) {
	// Build endpoint.
	urlstr := apiURL + "/" + "api/client/servers/" + uid + "/" + "power"

	// Setup form data.
	var formdata = []byte(`{"signal": "start"}`)

	// Setup HTTP GET request.
	client := &http.Client{Timeout: time.Second * 5}
	req, _ := http.NewRequest("POST", urlstr, bytes.NewBuffer(formdata))

	// Set authorization header.
	req.Header.Set("Authorization", "Bearer "+apiToken)

	// Set data type to JSON.
	req.Header.Set("Content-Type", "application/json")

	// Perform HTTP request and check for errors.
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
	}

	// Close body at the end.
	resp.Body.Close()
}
