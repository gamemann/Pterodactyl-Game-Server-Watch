package misc

import (
	"github.com/gamemann/Pterodactyl-Game-Server-Watch/config"
)

func RemoveIndex(cfg *config.Config, idx int) {
	copy(cfg.Servers[idx:], cfg.Servers[idx+1:])
	cfg.Servers = cfg.Servers[:len(cfg.Servers)-1]
}
