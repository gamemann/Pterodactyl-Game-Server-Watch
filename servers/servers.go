package servers

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/gamemann/Pterodactyl-Game-Server-Watch/config"
	"github.com/gamemann/Pterodactyl-Game-Server-Watch/events"
	"github.com/gamemann/Pterodactyl-Game-Server-Watch/pterodactyl"
	"github.com/gamemann/Pterodactyl-Game-Server-Watch/query"
)

var tickers []*time.Ticker

// Timer function.
func ServerWatch(srvidx int, timer *time.Ticker, fails *int, restarts *int, nextscan *int64, conn *net.UDPConn, cfg *config.Config) {
	destroy := make(chan struct{})

	for {
		select {
		case <-timer.C:
			srv := cfg.Servers[srvidx]

			// Check if container status is 'on'.
			if !pterodactyl.CheckStatus(cfg, srv.UID) {
				continue
			}

			// Send A2S_INFO request.
			query.SendRequest(conn)

			if cfg.DebugLevel > 2 {
				fmt.Println("[D3][" + srv.IP + ":" + strconv.Itoa(srv.Port) + "] A2S_INFO sent (" + srv.Name + ").")
			}

			// Check for response. If no response, increase fail count. Otherwise, reset fail count to 0.
			if !query.CheckResponse(conn, srv) {
				// Increase fail count.
				*fails++
				if cfg.DebugLevel > 1 {
					fmt.Println("[D2][" + srv.IP + ":" + strconv.Itoa(srv.Port) + "] Fails => " + strconv.Itoa(*fails))
				}

				// Check to see if we want to restart the server.
				if *fails >= srv.MaxFails && *restarts < srv.MaxRestarts && *nextscan < time.Now().Unix() {
					// Check if we want to restart the container.
					if !srv.ReportOnly {
						// Attempt to kill container.
						pterodactyl.KillServer(cfg, srv.UID)

						// Now attempt to start it again.
						pterodactyl.StartServer(cfg, srv.UID)
					}

					// Increment restarts count.
					*restarts++

					// Set next scan time and ensure the restart interval is at least 1.
					restartint := srv.RestartInt

					if restartint < 1 {
						restartint = 120
					}

					// Get new scan time.
					*nextscan = time.Now().Unix() + int64(restartint)

					// Debug.
					if cfg.DebugLevel > 0 {
						fmt.Println("[D1][" + srv.IP + ":" + strconv.Itoa(srv.Port) + "] Server found down. Report Only => " + strconv.FormatBool(srv.ReportOnly) + ". Fail Count => " + strconv.Itoa(*fails) + ". Restart Count => " + strconv.Itoa(*restarts) + " (" + srv.Name + ").")
					}

					events.OnServerDown(cfg, srvidx, *fails, *restarts)
				}
			} else {
				// Reset everything.
				*fails = 0
				*restarts = 0
				*nextscan = 0
			}

		case <-destroy:
			conn.Close()
			timer.Stop()
			return
		}
	}
}

func HandleServers(cfg *config.Config, update bool) {
	// Before anything, destroy any existing tickers.
	for _, tick := range tickers {
		if tick != nil {
			tick.Stop()
		}
	}

	// Empty tickets list.
	tickers = []*time.Ticker{}

	// Loop through each container from the config.
	for i, srv := range cfg.Servers {
		if cfg.DebugLevel > 0 && !update {
			fmt.Println("[D1] Adding server " + srv.IP + ":" + strconv.Itoa(srv.Port) + " with UID " + srv.UID + ". Auto Add => " + strconv.FormatBool(srv.ViaAPI) + ". Scan time => " + strconv.Itoa(srv.ScanTime) + ". Max Fails => " + strconv.Itoa(srv.MaxFails) + ". Max Restarts => " + strconv.Itoa(srv.MaxRestarts) + ". Restart Interval => " + strconv.Itoa(srv.RestartInt) + ". Report Only => " + strconv.FormatBool(srv.ReportOnly) + ". Enabled => " + strconv.FormatBool(srv.Enable) + ". Name => " + srv.Name + ". A2S Timeout => " + strconv.Itoa(srv.A2STimeout) + ". Mentions => " + srv.Mentions + ".")
		}

		// Check if server is enabled for scanning.
		if !srv.Enable {
			continue
		}

		// Specify server-specific variable.s
		var fails int = 0
		var restarts int = 0
		var nextscan int64 = 0

		// Get scan time.
		stime := srv.ScanTime

		if stime < 1 {
			stime = 5
		}

		// Let's create the connection now.
		conn, err := query.CreateConnection(srv.IP, srv.Port)

		if err != nil {
			fmt.Println("Error creating UDP connection for " + srv.IP + ":" + strconv.Itoa(srv.Port) + " ( " + srv.Name + ").")
			fmt.Println(err)

			continue
		}

		if cfg.DebugLevel > 3 {
			fmt.Println("[D4] Creating timer for " + srv.IP + ":" + strconv.Itoa(srv.Port) + ":" + srv.UID + " (" + srv.Name + ").")
		}

		// Create repeating timer.
		ticker := time.NewTicker(time.Duration(stime) * time.Second)
		go ServerWatch(i, ticker, &fails, &restarts, &nextscan, conn, cfg)

		// Add ticker to tickers variable.
		tickers = append(tickers, ticker)
	}
}
