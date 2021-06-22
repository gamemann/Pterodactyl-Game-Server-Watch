package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Server struct used for each server config.
type Server struct {
	Enable      bool   `json:"enable"`
	IP          string `json:"ip"`
	Port        int    `json:"port"`
	UID         string `json:"uid"`
	ScanTime    int    `json:"scantime"`
	MaxFails    int    `json:"maxfails"`
	MaxRestarts int    `json:"maxrestarts"`
	RestartInt  int    `json:"restartint"`
	ReportOnly  bool   `json:"reportonly"`
	ViaAPI      bool
}

// Misc options.
type Misc struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// Config struct used for the general config.
type Config struct {
	APIURL         string   `json:"apiurl"`
	Token          string   `json:"token"`
	AddServers     bool     `json:"addservers"`
	DebugLevel     int      `json:"debug"`
	ReloadTime     int      `json:"reloadtime"`
	DefEnable      bool     `json:"defenable"`
	DefScanTime    int      `json:"defscantime"`
	DefMaxFails    int      `json:"defmaxfails"`
	DefMaxRestarts int      `json:"defmaxrestarts"`
	DefRestartInt  int      `json:"defrestartint"`
	DefReportOnly  bool     `json:"defreportonly"`
	Servers        []Server `json:"servers"`
	Misc           []Misc   `json:"misc"`
	ConfLoc        string
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
