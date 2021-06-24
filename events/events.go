package events

import (
	"github.com/gamemann/Pterodactyl-Game-Server-Watch/config"
	"github.com/gamemann/Pterodactyl-Game-Server-Watch/misc"
)

func OnServerDown(cfg *config.Config, srv *config.Server, fails int, restarts int) {
	// Handle Misc options.
	misc.HandleMisc(cfg, srv, fails, restarts)
}
