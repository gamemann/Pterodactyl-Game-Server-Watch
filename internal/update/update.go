package update

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gamemann/Pterodactyl-Game-Server-Watch/internal/pterodactyl"
	"github.com/gamemann/Pterodactyl-Game-Server-Watch/internal/servers"
	"github.com/gamemann/Pterodactyl-Game-Server-Watch/pkg/config"
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
		if cfg.DebugLevel > 3 {
			fmt.Println("[D4] Looking for " + newsrv.IP + ":" + strconv.Itoa(newsrv.Port) + ":" + newsrv.UID + " (" + strconv.Itoa(i) + ") inside of old configuration.")
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

			if cfg.DebugLevel > 4 {
				fmt.Println("[D5] Comparing " + nt.IP + ":" + strconv.Itoa(nt.Port) + ":" + nt.UID + " == " + ot.IP + ":" + strconv.Itoa(ot.Port) + ":" + ot.UID + " (" + strconv.Itoa(j) + ").")
			}

			// Now compare.
			if nt == ot {
				// We don't have to insert this server into the slice.
				toadd = false

				if cfg.DebugLevel > 2 {
					fmt.Println("[D3] Found matching server (" + newsrv.IP + ":" + strconv.Itoa(newsrv.Port) + ":" + newsrv.UID + ") on Add Server check. Applying new configuration. Name => " + newsrv.Name + ". Enabled: " + strconv.FormatBool(oldsrv.Enable) + " => " + strconv.FormatBool(newsrv.Enable) + ". Max fails: " + strconv.Itoa(oldsrv.MaxFails) + " => " + strconv.Itoa(newsrv.MaxFails) + ". Max Restarts: " + strconv.Itoa(oldsrv.MaxRestarts) + " => " + strconv.Itoa(newsrv.MaxRestarts) + ". Restart Int: " + strconv.Itoa(oldsrv.RestartInt) + " => " + strconv.Itoa(newsrv.RestartInt) + ". Scan Time: " + strconv.Itoa(oldsrv.ScanTime) + " => " + strconv.Itoa(newsrv.ScanTime) + ". Report Only: " + strconv.FormatBool(oldsrv.ReportOnly) + " => " + strconv.FormatBool(newsrv.ReportOnly) + ". A2S Timeout: " + strconv.Itoa(oldsrv.A2STimeout) + " => " + strconv.Itoa(newsrv.A2STimeout) + ". Mentions: " + oldsrv.Mentions + " => " + newsrv.Mentions + ".")
				}

				// Update specific configuration.
				cfg.Servers[j].Enable = newsrv.Enable
				cfg.Servers[j].MaxFails = newsrv.MaxFails
				cfg.Servers[j].MaxRestarts = newsrv.MaxRestarts
				cfg.Servers[j].RestartInt = newsrv.RestartInt
				cfg.Servers[j].ScanTime = newsrv.ScanTime
				cfg.Servers[j].ReportOnly = newsrv.ReportOnly
				cfg.Servers[j].A2STimeout = newsrv.A2STimeout
				cfg.Servers[j].RconPassword = newsrv.RconPassword
				cfg.Servers[j].Mentions = newsrv.Mentions
			}
		}

		// If we're not inside of the current configuration, add the server.
		if toadd {
			if cfg.DebugLevel > 1 {
				fmt.Println("[D2] Adding server from update " + newsrv.IP + ":" + strconv.Itoa(newsrv.Port) + " with UID " + newsrv.UID + ". Name => " + newsrv.Name + ". Auto Add => " + strconv.FormatBool(newsrv.ViaAPI) + ". Scan time => " + strconv.Itoa(newsrv.ScanTime) + ". Max Fails => " + strconv.Itoa(newsrv.MaxFails) + ". Max Restarts => " + strconv.Itoa(newsrv.MaxRestarts) + ". Restart Interval => " + strconv.Itoa(newsrv.RestartInt) + ". Enabled => " + strconv.FormatBool(newsrv.Enable) + ". A2S Timeout => " + strconv.Itoa(newsrv.A2STimeout) + ". RCON Password => " + newsrv.RconPassword + ". Mentions => " + newsrv.Mentions + ".")
			}

			cfg.Servers = append(cfg.Servers, newsrv)
		}
	}
}

func DelOldServers(newcfg *config.Config, cfg *config.Config) {
	// Loop through all old servers.
	for i, oldsrv := range cfg.Servers {
		if cfg.DebugLevel > 3 {
			fmt.Println("[D4] Looking for " + oldsrv.IP + ":" + strconv.Itoa(oldsrv.Port) + ":" + oldsrv.UID + " (" + strconv.Itoa(i) + ") inside of new configuration.")
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

			if cfg.DebugLevel > 4 {
				fmt.Println("[D5] Comparing " + ot.IP + ":" + strconv.Itoa(ot.Port) + ":" + ot.UID + " == " + nt.IP + ":" + strconv.Itoa(nt.Port) + ":" + nt.UID + " (" + strconv.Itoa(j) + ").")
			}

			// Now compare.
			if nt == ot {
				todel = false
			}
		}

		// If we're not inside of the new configuration, delete the server.
		if todel {
			if cfg.DebugLevel > 1 {
				fmt.Println("[D2] Deleting server from update " + oldsrv.IP + ":" + strconv.Itoa(oldsrv.Port) + " with UID " + oldsrv.UID + ". Name => " + oldsrv.Name + ".")
			}

			// Set Delete to true so we'll delete the server, close the connection, etc. on the next scan.
			cfg.Servers[i].Delete = true
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
			newcfg.SetDefaults()

			err := newcfg.ReadConfig(cfg.ConfLoc)

			if err != nil {
				fmt.Println(err)

				continue
			}

			if newcfg.AddServers {
				cont := pterodactyl.AddServers(&newcfg)

				if !cont {
					fmt.Println("[ERR] Not updating server list due to error.")

					continue
				}
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
