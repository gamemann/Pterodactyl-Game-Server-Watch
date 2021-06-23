package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Server struct used for each server config.
type Server struct {
	Name        string `json:"name"`
	Enable      bool   `json:"enable"`
	IP          string `json:"ip"`
	Port        int    `json:"port"`
	UID         string `json:"uid"`
	ScanTime    int    `json:"scantime"`
	MaxFails    int    `json:"maxfails"`
	MaxRestarts int    `json:"maxrestarts"`
	RestartInt  int    `json:"restartint"`
	ReportOnly  bool   `json:"reportonly"`
	A2STimeout  int    `json:"a2stimeout"`
	Mentions    string `json:"mentions"`
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
	AppToken       string   `json:"apptoken"`
	AddServers     bool     `json:"addservers"`
	DebugLevel     int      `json:"debug"`
	ReloadTime     int      `json:"reloadtime"`
	DefEnable      bool     `json:"defenable"`
	DefScanTime    int      `json:"defscantime"`
	DefMaxFails    int      `json:"defmaxfails"`
	DefMaxRestarts int      `json:"defmaxrestarts"`
	DefRestartInt  int      `json:"defrestartint"`
	DefReportOnly  bool     `json:"defreportonly"`
	DefA2STimeout  int      `json:"defa2stimeout"`
	DefMentions    string   `json:"defmentions"`
	Servers        []Server `json:"servers"`
	Misc           []Misc   `json:"misc"`
	ConfLoc        string
}

// Reads a config file based off of the file name (string) and returns a Config struct.
func (cfg *Config) ReadConfig(filename string) bool {
	file, err := os.Open(filename)

	if err != nil {
		fmt.Println("[ERR] Cannot open config file.")
		fmt.Println(err)

		return false
	}

	stat, _ := file.Stat()

	data := make([]byte, stat.Size())

	_, err = file.Read(data)

	if err != nil {
		fmt.Println("[ERR] Cannot read config file.")
		fmt.Println(err)

		return false
	}

	err = json.Unmarshal([]byte(data), cfg)

	file.Close()

	if err != nil {
		fmt.Println("[ERR] Cannot parse JSON Data.")
		fmt.Println(err)

		return false
	}

	return true
}

// Sets config's default values.
func (cfg *Config) SetDefaults() {
	// Set config defaults.
	cfg.AddServers = false
	cfg.DebugLevel = 0
	cfg.ReloadTime = 500

	cfg.DefEnable = true
	cfg.DefScanTime = 5
	cfg.DefMaxFails = 10
	cfg.DefMaxRestarts = 2
	cfg.DefRestartInt = 120
	cfg.DefReportOnly = false
	cfg.DefA2STimeout = 1
}
