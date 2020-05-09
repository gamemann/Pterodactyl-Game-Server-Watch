package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Server struct used for each server config.
type Server struct {
	IP       string `json:"IP"`
	Port     int    `json:"port"`
	UID      string `json:"uid"`
	ScanTime int    `json:"scantime"`
}

// Config struct used for the general config.
type Config struct {
	APIURL  string   `json:"apiURL"`
	Token   string   `json:"token"`
	Servers []Server `json:"servers"`
}

// Reads a config file based off of the file name (string) and returns a Config struct.
func ReadConfig(cfg *Config, filename string) bool {

	file, err := os.Open(filename)

	if err != nil {
		fmt.Println("Error opening config file.")

		return false
	}

	stat, _ := file.Stat()

	data := make([]byte, stat.Size())

	_, err = file.Read(data)

	if err != nil {
		fmt.Println("Error reading config file.")

		return false
	}

	err = json.Unmarshal([]byte(data), cfg)

	if err != nil {
		fmt.Println("Error parsing JSON Data.")

		return false
	}

	return true
}
