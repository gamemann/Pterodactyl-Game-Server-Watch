package update

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gamemann/Pterodactyl-Game-Server-Watch/config"
	"github.com/gamemann/Pterodactyl-Game-Server-Watch/misc"
	"github.com/gamemann/Pterodactyl-Game-Server-Watch/pterodactyl"
	"github.com/gamemann/Pterodactyl-Game-Server-Watch/servers"
)

type Tuple struct {
	IP   string
	Port int
	UID  string
}

var updateticker *time.Ticker

func AddNewServers(newcfg *config.Config, cfg *config.Config) {
	// Loop through all new servers.
	for i, newsrv := range newcfg.Servers {
		if cfg.DebugLevel > 2 {
			fmt.Println("[D3] Looking for " + newsrv.IP + ":" + strconv.Itoa(newsrv.Port) + ":" + newsrv.UID + " (" + strconv.Itoa(i) + ") inside of old configuration.")
		}

		toadd := true

		// Now loop through old servers.
		for j, oldsrv := range cfg.Servers {
			// Create new server tuple.
			var nt Tuple
			nt.IP = newsrv.IP
			nt.Port = newsrv.Port
			nt.UID = newsrv.UID

			// Create old server tuple.
			var ot Tuple
			ot.IP = oldsrv.IP
			ot.Port = oldsrv.Port
			ot.UID = oldsrv.UID

			if cfg.DebugLevel > 2 {
				fmt.Println("[D3] Comparing " + nt.IP + ":" + strconv.Itoa(nt.Port) + ":" + nt.UID + " == " + ot.IP + ":" + strconv.Itoa(ot.Port) + ":" + ot.UID + " (" + strconv.Itoa(j) + ").")
			}

			// Now compare.
			if nt == ot {
				// We don't have to insert this server into the slice.
				toadd = false

				if cfg.DebugLevel > 3 {
					fmt.Println("[D4] Found matching server on Add Server check. Applying new configuration. Enabled: " + strconv.FormatBool(oldsrv.Enable) + " => " + strconv.FormatBool(newsrv.Enable) + ". Max fails: " + strconv.Itoa(oldsrv.MaxFails) + " => " + strconv.Itoa(newsrv.MaxFails) + ". Max Restarts: " + strconv.Itoa(oldsrv.MaxRestarts) + " => " + strconv.Itoa(newsrv.MaxRestarts) + ". Restart Int: " + strconv.Itoa(oldsrv.RestartInt) + " => " + strconv.Itoa(newsrv.RestartInt) + ". Scan Time: " + strconv.Itoa(oldsrv.ScanTime) + " => " + strconv.Itoa(newsrv.ScanTime) + ". Report Only: " + strconv.FormatBool(oldsrv.ReportOnly) + " => " + strconv.FormatBool(newsrv.ReportOnly) + ".")
				}

				// Update specific configuration.
				cfg.Servers[i].Enable = newsrv.Enable
				cfg.Servers[i].MaxFails = newsrv.MaxFails
				cfg.Servers[i].MaxRestarts = newsrv.MaxRestarts
				cfg.Servers[i].RestartInt = newsrv.RestartInt
				cfg.Servers[i].ScanTime = newsrv.ScanTime
				cfg.Servers[i].ReportOnly = newsrv.ReportOnly
			}
		}

		// If we're not inside of the current configuration, add the server.
		if toadd {
			if cfg.DebugLevel > 0 {
				fmt.Println("[D1] Adding server from update " + newsrv.IP + ":" + strconv.Itoa(newsrv.Port) + " with UID " + newsrv.UID + ". Auto Add => " + strconv.FormatBool(newsrv.ViaAPI) + ". Scan time => " + strconv.Itoa(newsrv.ScanTime) + ". Max Fails => " + strconv.Itoa(newsrv.MaxFails) + ". Max Restarts => " + strconv.Itoa(newsrv.MaxRestarts) + ". Restart Interval => " + strconv.Itoa(newsrv.RestartInt) + ". Enabled => " + strconv.FormatBool(newsrv.Enable) + ".")
			}

			cfg.Servers = append(cfg.Servers, newsrv)
		}
	}
}

func DelOldServers(newcfg *config.Config, cfg *config.Config) {
	// Loop through all old servers.
	for i, oldsrv := range cfg.Servers {
		if cfg.DebugLevel > 2 {
			fmt.Println("[D3] Looking for " + oldsrv.IP + ":" + strconv.Itoa(oldsrv.Port) + ":" + oldsrv.UID + " (" + strconv.Itoa(i) + ") inside of new configuration.")
		}

		todel := true

		// Now loop through new servers.
		for j, newsrv := range newcfg.Servers {
			// Create old server tuple.
			var ot Tuple
			ot.IP = oldsrv.IP
			ot.Port = oldsrv.Port
			ot.UID = oldsrv.UID

			// Create new server tuple.
			var nt Tuple
			nt.IP = newsrv.IP
			nt.Port = newsrv.Port
			nt.UID = newsrv.UID

			if cfg.DebugLevel > 2 {
				fmt.Println("[D3] Comparing " + ot.IP + ":" + strconv.Itoa(ot.Port) + ":" + ot.UID + " == " + nt.IP + ":" + strconv.Itoa(nt.Port) + ":" + nt.UID + " (" + strconv.Itoa(j) + ").")
			}

			// Now compare.
			if nt == ot {
				todel = false
			}
		}

		// If we're not inside of the new configuration, delete the server.
		if todel {
			if cfg.DebugLevel > 0 {
				fmt.Println("[D1] Deleting server from update " + oldsrv.IP + ":" + strconv.Itoa(oldsrv.Port) + " with UID " + oldsrv.UID + ".")
			}

			misc.RemoveIndex(cfg, i)
		}
	}
}

func ReloadServers(timer *time.Ticker, cfg *config.Config) {
	destroy := make(chan struct{})

	for {
		select {
		case <-timer.C:
			// First, we'll want to read the new config.
			newcfg := config.Config{}

			// Set default values.
			newcfg.AddServers = false
			newcfg.DebugLevel = 0
			newcfg.ReloadTime = 500

			newcfg.DefEnable = true
			newcfg.DefScanTime = 5
			newcfg.DefMaxFails = 10
			newcfg.DefMaxRestarts = 2
			newcfg.DefRestartInt = 120
			newcfg.DefReportOnly = false

			success := config.ReadConfig(&newcfg, cfg.ConfLoc)

			if !success {
				continue
			}

			if newcfg.AddServers {
				pterodactyl.AddServers(&newcfg)
			}

			// Assign new values.
			cfg.APIURL = newcfg.APIURL
			cfg.Token = newcfg.Token
			cfg.DebugLevel = newcfg.DebugLevel
			cfg.AddServers = newcfg.AddServers

			cfg.DefEnable = newcfg.DefEnable
			cfg.DefScanTime = newcfg.DefScanTime
			cfg.DefMaxFails = newcfg.DefMaxFails
			cfg.DefMaxRestarts = newcfg.DefMaxRestarts
			cfg.DefRestartInt = newcfg.DefRestartInt
			cfg.DefReportOnly = newcfg.DefReportOnly

			// If reload time is different, recreate reload timer.
			if cfg.ReloadTime != newcfg.ReloadTime {
				if updateticker != nil {
					updateticker.Stop()
				}

				if cfg.DebugLevel > 2 {
					fmt.Println("[D3] Recreating update timer due to updated reload time (" + strconv.Itoa(cfg.ReloadTime) + " => " + strconv.Itoa(newcfg.ReloadTime) + ").")
				}

				// Create repeating timer.
				updateticker = time.NewTicker(time.Duration(newcfg.ReloadTime) * time.Second)
				go ReloadServers(updateticker, cfg)
			}

			cfg.ReloadTime = newcfg.ReloadTime

			// Level 2 debug message.
			if cfg.DebugLevel > 1 {
				fmt.Println("[D2] Updating server list.")
			}

			// Add new servers.
			AddNewServers(&newcfg, cfg)

			// Remove servers that are not a part of new configuration.
			DelOldServers(&newcfg, cfg)

			// Now rehandle servers.
			servers.HandleServers(cfg, true)

		case <-destroy:
			timer.Stop()

			return
		}
	}
}

func Init(cfg *config.Config) {
	if cfg.ReloadTime < 1 {
		return
	}

	if cfg.DebugLevel > 0 {
		fmt.Println("[D1] Setting up reload timer for every " + strconv.Itoa(cfg.ReloadTime) + " seconds.")
	}

	// Create repeating timer.
	updateticker = time.NewTicker(time.Duration(cfg.ReloadTime) * time.Second)
	go ReloadServers(updateticker, cfg)
}
