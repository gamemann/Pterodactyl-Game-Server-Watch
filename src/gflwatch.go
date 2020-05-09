package main

import (
	"./config"
	"fmt"
)

func main() {
	configFile := "/etc/gflwatch/gflwatch.conf"

	cfg := config.Config{}

	config.ReadConfig(&cfg, configFile)

	fmt.Println(cfg)
}
