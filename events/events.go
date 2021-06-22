package events

import (
	"github.com/gamemann/Pterodactyl-Game-Server-Watch/config"
	"github.com/gamemann/Pterodactyl-Game-Server-Watch/misc"
)

func OnServerDown(cfg *config.Config, srvidx int, fails int, restarts int) {
	// Handle Misc options.
	misc.HandleMisc(cfg, srvidx, fails, restarts)
}
