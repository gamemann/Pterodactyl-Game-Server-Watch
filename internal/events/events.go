package events

import (
	"github.com/gamemann/Pterodactyl-Game-Server-Watch/internal/misc"
	"github.com/gamemann/Pterodactyl-Game-Server-Watch/pkg/config"
)

func OnServerDown(cfg *config.Config, srv *config.Server, fails int, restarts int) {
	// Handle Misc options.
	misc.HandleMisc(cfg, srv, fails, restarts)
}
