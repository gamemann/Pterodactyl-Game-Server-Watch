package config

import (
	"encoding/json"
	"os"
)

// Reads a config file based off of the file name (string) and returns a Config struct.
func (cfg *Config) ReadConfig(path string) error {
	file, err := os.Open(path)

	if err != nil {
		return err
	}

	defer file.Close()

	stat, _ := file.Stat()

	data := make([]byte, stat.Size())

	_, err = file.Read(data)

	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(data), cfg)

	return err
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
