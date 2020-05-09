package main

import (
	"./config"
	"./query"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

// Timer function.
func ServerWatch(server config.Server, timer *time.Ticker, fails *int, restarts *int, conn *net.UDPConn) {
	destroy := make(chan struct{})

	for {
		select {
		case <-timer.C:
			// Send A2S_INFO request.
			query.SendRequest(conn)

			// Check for response. If no response, increase fail count. Otherwise, reset fail count to 0.
			if !query.CheckResponse(conn) {
				*fails++
			} else {
				*fails = 0
				*restarts = 0
			}

			fmt.Println(server.IP + ":" + strconv.Itoa(server.Port) + " => " + strconv.Itoa(*fails) + " fails...")

		case <-destroy:
			conn.Close()
			timer.Stop()
			return
		}
	}
}

func main() {
	// Specify config file path.
	configFile := "/etc/pterowatch/pterowatch.conf"

	// Create config struct.
	cfg := config.Config{}

	// Attempt to read config.
	config.ReadConfig(&cfg, configFile)

	// Loop through each container from the config.
	for i := 0; i < len(cfg.Servers); i++ {
		// Specify fails variable.
		var fails int = 0
		var restarts int = 0

		// Get scan time.
		stime := cfg.Servers[i].ScanTime

		if stime < 1 {
			stime = 5
		}

		// Let's create the connection now.
		conn, err := query.CreateConnection(cfg.Servers[i].IP, cfg.Servers[i].Port)

		if err != nil {
			fmt.Println("Error creating UDP connection for " + cfg.Servers[i].IP + ":" + strconv.Itoa(cfg.Servers[i].Port))
			fmt.Println(err)

			return
		}

		// Create repeating timer.
		ticker := time.NewTicker(time.Duration(stime) * time.Second)
		go ServerWatch(cfg.Servers[i], ticker, &fails, &restarts, conn)
	}

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
