package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/gamemann/Pterodactyl-Game-Server-Watch/internal/pterodactyl"
	"github.com/gamemann/Pterodactyl-Game-Server-Watch/internal/servers"
	"github.com/gamemann/Pterodactyl-Game-Server-Watch/internal/update"
	"github.com/gamemann/Pterodactyl-Game-Server-Watch/pkg/config"
)

func main() {
	// Look for 'cfg' flag in command line arguments (default path: /etc/pterowatch/pterowatch.conf).
	configFile := flag.String("cfg", "/etc/pterowatch/pterowatch.conf", "The path to the Pterowatch config file.")
	flag.Parse()

	// Create config struct.
	cfg := config.Config{}

	// Set config defaults.
	cfg.SetDefaults()

	// Attempt to read config.
	err := cfg.ReadConfig(*configFile)

	// If we have no config, create the file with the defaults.
	if err != nil {
		// If there's an error and it contains "no such file", try to create the file with defaults.
		if strings.Contains(err.Error(), "no such file") {
			err = cfg.WriteDefaultsToFile(*configFile)

			if err != nil {
				fmt.Println("Failed to open config file and cannot create file.")
				fmt.Println(err)

				os.Exit(1)
			}
		}

		fmt.Println("WARNING - No config file found. Created config file at " + *configFile + " with defaults.")
	} else {
		// Check if we want to automatically add servers.
		if cfg.AddServers {
			pterodactyl.AddServers(&cfg)
		}
	}

	// Level 1 debug.
	if cfg.DebugLevel > 0 {
		fmt.Println("[D1] Found config with API URL => " + cfg.APIURL + ". Token => " + cfg.Token + ". App Token => " + cfg.AppToken + ". Auto Add Servers => " + strconv.FormatBool(cfg.AddServers) + ". Debug level => " + strconv.Itoa(cfg.DebugLevel) + ". Reload time => " + strconv.Itoa(cfg.ReloadTime))
	}

	// Level 2 debug.
	if cfg.DebugLevel > 1 {
		fmt.Println("[D2] Config default server values. Enable => " + strconv.FormatBool(cfg.DefEnable) + ". Scan time => " + strconv.Itoa(cfg.DefScanTime) + ". Max Fails => " + strconv.Itoa(cfg.DefMaxFails) + ". Max Restarts => " + strconv.Itoa(cfg.DefMaxRestarts) + ". Restart Interval => " + strconv.Itoa(cfg.DefRestartInt) + ". Report Only => " + strconv.FormatBool(cfg.DefReportOnly) + ". A2S Timeout => " + strconv.Itoa(cfg.DefA2STimeout) + ". RCON Password => " + cfg.DefRconPassword + ". Mentions => " + cfg.DefMentions + ".")
	}

	// Handle all servers (create timers, etc.).
	servers.HandleServers(&cfg, false)

	// Set config file for use later (e.g. updating/reloading).
	cfg.ConfLoc = *configFile

	// Initialize updater/reloader.
	update.Init(&cfg)

	// Signal.
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	<-sigc
}
