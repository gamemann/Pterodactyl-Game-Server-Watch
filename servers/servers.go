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

type Tuple struct {
	IP   string
	Port int
	UID  string
}

type Stats struct {
	Fails    *int
	Restarts *int
	NextScan *int64
}

type TickerHolder struct {
	Info      Tuple
	Ticker    *time.Ticker
	Conn      *net.UDPConn
	ScanTime  int
	Destroyer *chan bool
	Idx       *int
	Stats     Stats
}

var tickers []TickerHolder

func RemoveTicker(t *[]TickerHolder, idx int) {
	copy((*t)[idx:], (*t)[idx+1:])
	*t = (*t)[:len(*t)-1]
}

func RemoveServer(cfg *config.Config, idx int) {
	copy(cfg.Servers[idx:], cfg.Servers[idx+1:])
	cfg.Servers = cfg.Servers[:len(cfg.Servers)-1]
}

// Timer function.
func ServerWatch(srv *config.Server, timer *time.Ticker, fails *int, restarts *int, nextscan *int64, conn *net.UDPConn, cfg *config.Config, destroy *chan bool) {
	for {
		select {
		case <-timer.C:
			// If the UDP connection or server is nil, break the timer.
			if conn == nil || srv == nil {
				*destroy <- true

				break
			}

			// Check if server is enabled.
			if !srv.Enable {
				continue
			}

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
			if !query.CheckResponse(conn, *srv) {
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

					events.OnServerDown(cfg, srv, *fails, *restarts)
				}
			} else {
				// Reset everything.
				*fails = 0
				*restarts = 0
				*nextscan = 0
			}

		case <-*destroy:
			// Close UDP connection and check.
			err := conn.Close()

			if err != nil {
				fmt.Println("[ERR] Failed to close UDP connection.")
				fmt.Println(err)
			}

			// Stop timer/ticker.
			timer.Stop()

			// Stop function.
			return
		}
	}
}

func HandleServers(cfg *config.Config, update bool) {
	stats := make(map[Tuple]Stats)

	// Retrieve current server stats before removing tickers
	for _, srvticker := range tickers {
		for _, srv := range cfg.Servers {
			// Create tuple.
			var srvt Tuple
			srvt.IP = srv.IP
			srvt.Port = srv.Port
			srvt.UID = srv.UID

			if cfg.DebugLevel > 4 {
				fmt.Println("[D5] HandleServers :: Comparing " + srvt.IP + ":" + strconv.Itoa(srvt.Port) + ":" + srvt.UID + " == " + srv.IP + ":" + strconv.Itoa(srv.Port) + ":" + srv.UID + ".")
			}

			if srvt == srvticker.Info {
				if cfg.DebugLevel > 3 {
					fmt.Println("[D4] HandleServers :: Found match on " + srvt.IP + ":" + strconv.Itoa(srvt.Port) + ":" + srvt.UID + ".")
				}

				// Fill in stats.
				stats[srvt] = Stats{
					Fails:    srvticker.Stats.Fails,
					Restarts: srvticker.Stats.Restarts,
					NextScan: srvticker.Stats.NextScan,
				}

			}
		}

		// Destroy ticker.
		*srvticker.Destroyer <- true
	}

	// Remove servers that should be deleted.
	for i, srv := range cfg.Servers {
		if srv.Delete {
			if cfg.DebugLevel > 1 {
				fmt.Println("[D2] Found server that should be deleted UID => " + srv.UID + ". Name => " + srv.Name + ". IP => " + srv.IP + ". Port => " + strconv.Itoa(srv.Port) + ".")
			}

			RemoveServer(cfg, i)
		}
	}

	tickers = []TickerHolder{}

	// Loop through each container from the config.
	for i, srv := range cfg.Servers {
		// If we're not enabled, ignore.
		if !srv.Enable {
			continue
		}

		// Create tuple.
		var srvt Tuple
		srvt.IP = srv.IP
		srvt.Port = srv.Port
		srvt.UID = srv.UID

		// Specify server-specific variables
		var fails int = 0
		var restarts int = 0
		var nextscan int64 = 0

		// Replace stats with old ticker's stats.
		if stat, ok := stats[srvt]; ok {
			fails = *stat.Fails
			restarts = *stat.Restarts
			nextscan = *stat.NextScan
		}

		if cfg.DebugLevel > 0 && !update {
			fmt.Println("[D1] Adding server " + srv.IP + ":" + strconv.Itoa(srv.Port) + " with UID " + srv.UID + ". Auto Add => " + strconv.FormatBool(srv.ViaAPI) + ". Scan time => " + strconv.Itoa(srv.ScanTime) + ". Max Fails => " + strconv.Itoa(srv.MaxFails) + ". Max Restarts => " + strconv.Itoa(srv.MaxRestarts) + ". Restart Interval => " + strconv.Itoa(srv.RestartInt) + ". Report Only => " + strconv.FormatBool(srv.ReportOnly) + ". Enabled => " + strconv.FormatBool(srv.Enable) + ". Name => " + srv.Name + ". A2S Timeout => " + strconv.Itoa(srv.A2STimeout) + ". Mentions => " + srv.Mentions + ".")
		}

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

		// Create destroyer channel.
		destroyer := make(chan bool)

		// Create repeating timer.
		ticker := time.NewTicker(time.Duration(stime) * time.Second)
		go ServerWatch(&cfg.Servers[i], ticker, &fails, &restarts, &nextscan, conn, cfg, &destroyer)

		// Add ticker to global list.
		var newticker TickerHolder
		newticker.Info = srvt
		newticker.Ticker = ticker
		newticker.Conn = conn
		newticker.ScanTime = stime
		newticker.Destroyer = &destroyer
		newticker.Stats.Fails = &fails
		newticker.Stats.Restarts = &restarts
		newticker.Stats.NextScan = &nextscan

		tickers = append(tickers, newticker)
	}
}
