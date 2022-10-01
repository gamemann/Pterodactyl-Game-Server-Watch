package misc

import (
	"strconv"
	"strings"

	"github.com/gamemann/Pterodactyl-Game-Server-Watch/pkg/config"
)

func FormatContents(app string, formatstr *string, fails int, restarts int, srv *config.Server, mentionstr string) {
	*formatstr = strings.ReplaceAll(*formatstr, "{IP}", srv.IP)
	*formatstr = strings.ReplaceAll(*formatstr, "{PORT}", strconv.Itoa(srv.Port))
	*formatstr = strings.ReplaceAll(*formatstr, "{FAILS}", strconv.Itoa(fails))
	*formatstr = strings.ReplaceAll(*formatstr, "{RESTARTS}", strconv.Itoa(restarts))
	*formatstr = strings.ReplaceAll(*formatstr, "{MAXFAILS}", strconv.Itoa(srv.MaxFails))
	*formatstr = strings.ReplaceAll(*formatstr, "{MAXRESTARTS}", strconv.Itoa(srv.MaxRestarts))
	*formatstr = strings.ReplaceAll(*formatstr, "{UID}", srv.UID)
	*formatstr = strings.ReplaceAll(*formatstr, "{SCANTIME}", strconv.Itoa(srv.ScanTime))
	*formatstr = strings.ReplaceAll(*formatstr, "{RESTARTINT}", strconv.Itoa(srv.RestartInt))
	*formatstr = strings.ReplaceAll(*formatstr, "{NAME}", srv.Name)
	*formatstr = strings.ReplaceAll(*formatstr, "{MENTIONS}", mentionstr)
}
