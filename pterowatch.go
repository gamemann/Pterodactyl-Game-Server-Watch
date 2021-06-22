package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gamemann/Pterodactyl-Game-Server-Watch/config"
	"github.com/gamemann/Pterodactyl-Game-Server-Watch/pterodactyl"
	"github.com/gamemann/Pterodactyl-Game-Server-Watch/servers"
	"github.com/gamemann/Pterodactyl-Game-Server-Watch/update"
)

func main() {
	// Specify config file path.
	configFile := "/etc/pterowatch/pterowatch.conf"

	// Create config struct.
	cfg := config.Config{}

	// Attempt to read config.
	config.ReadConfig(&cfg, configFile)

	// Level 1 debug.
	if cfg.DebugLevel > 0 {
		fmt.Println("[D1] Found config with API URL => " + cfg.APIURL + ". Token => " + cfg.Token + ". Auto Add Servers => " + strconv.FormatBool(cfg.AddServers) + ". Debug level => " + strconv.Itoa(cfg.DebugLevel))
	}

	// Check if we want to automatically add servers.
	if cfg.AddServers {
		pterodactyl.AddServers(&cfg)
	}

	// Handle all servers (create timers, etc.).
	servers.HandleServers(&cfg, false)

	// Set config file for use later (e.g. updating/reloading).
	cfg.ConfLoc = configFile

	// Initialize updater/reloader.
	update.Init(&cfg)

	// Signal.
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT)

	x := 0

	// Create a loop so the program doesn't exit. Look for signals and if SIGINT, stop the program.
	for x < 1 {
		kill := false
		s := <-sigc

		switch s {
		case os.Interrupt:
			kill = true
		}

		if kill {
			break
		}

		// Sleep every second to avoid unnecessary CPU consumption.
		time.Sleep(time.Duration(1) * time.Second)
	}
}
