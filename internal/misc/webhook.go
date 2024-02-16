package misc

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/gamemann/Pterodactyl-Game-Server-Watch/pkg/config"
	"github.com/gookit/goutil/dump"
)

type AllowMentions struct {
	Roles bool `json:"roles"`
	Users bool `json:"users"`
}

type DiscordWH struct {
	Contents  string `json:"content"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatarurl"`
}

type SlackWH struct {
	Text string `json:"text"`
}
type NtfyWH struct {
	Topic    string   `json:"topic"`
	Title    string   `json:"title"`
	Message  string   `json:"message"`
	Priority int      `json:"priority"`
	Tags     []string `json:"tags"`
}

func DiscordWebHook(url string, contents string, username string, avatarurl string, allowmentions AllowMentions, srv *config.Server) bool {
	var data DiscordWH

	// Build out JSON/form data for Discord web hook.
	data.Contents = contents
	data.Username = username
	data.AvatarURL = avatarurl

	// Convert interface to JSON data string.
	datajson, err := json.Marshal(data)

	if err != nil {
		fmt.Println(err)

		return false
	}

	// Setup HTTP POST request with form data.
	client := &http.Client{Timeout: time.Second * 5}
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(datajson))

	// Set content type to JSON.
	req.Header.Set("Content-Type", "application/json")

	// Perform HTTP request and check for errors.
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)

		return false
	}

	resp.Body.Close()

	return true
}

func SlackWebHook(url string, contents string) bool {
	var data SlackWH

	// Build out JSON/form data for Slack web hook.
	data.Text = contents

	// Convert interface to JSON data string.
	datajson, err := json.Marshal(data)

	if err != nil {
		fmt.Println(err)

		return false
	}

	// Setup HTTP POST request with form data.
	client := &http.Client{Timeout: time.Second * 5}
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(datajson))

	// Set content type to JSON.
	req.Header.Set("Content-Type", "application/json")

	// Perform HTTP request and check for errors.
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)

		return false
	}

	resp.Body.Close()

	return true
}

func NtfyWebHook(url string, contents string, topic string, username string, password string) bool {
	var data NtfyWH

	// Build out JSON/form data for Slack web hook.
	data.Message = contents
	data.Topic = topic

	dump.P(data)

	// Convert interface to JSON data string.
	datajson, err := json.Marshal(data)

	if err != nil {
		fmt.Println(err)

		return false
	}

	// Setup HTTP POST request with form data.
	client := &http.Client{Timeout: time.Second * 5}
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(datajson))

	// Set content type to JSON.
	req.Header.Set("Content-Type", "application/json")

	// Authenticate
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))
	req.Header.Set("Authorization", authHeader)

	// Perform HTTP request and check for errors.
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)

		return false
	}

	defer resp.Body.Close()

	b, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(b))

	return true
}
