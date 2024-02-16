package query

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/gamemann/Pterodactyl-Game-Server-Watch/pkg/config"

	"github.com/gamemann/Pterodactyl-Game-Server-Watch/internal/rcon"
)

// Creates a RCON connection using the host and port.
func CreateConnection(hostname string, port int, srv *config.Server) (*rcon.RemoteConsole, error) {
	request_timeout := srv.A2STimeout
	password := srv.RconPassword

	// dump.P(srv)

	hostPortArgs := []string{}
	hostPortArgs = append(hostPortArgs, hostname)
	hostPortArgs = append(hostPortArgs, strconv.Itoa(port))
	hostPort := strings.Join(hostPortArgs, ":")

	remoteConsole, err := rcon.Dial(hostPort, password, request_timeout)
	if err != nil {
		// log.Fatal("Failed to connect to RCON server", err)
		fmt.Println(err)

		return remoteConsole, errors.New("NoConnection")
	}
	// defer remoteConsole.Close()

	return remoteConsole, nil

}

func CloseConnect(conn *rcon.RemoteConsole) error {
	if conn != nil {
		return conn.Close()
	}
	return nil
}

// Sends an RCON request to the host and port specified in the arguments.
func SendRequest(conn *rcon.RemoteConsole) (int, error) {
	preparedCmd := "KeepAlive"
	preparedCmd = "ListPlayers"
	reqId, err := conn.Write(preparedCmd)

	return reqId, err
}

// Checks for RCON response. Returns true if it receives a response. Returns false otherwise.
func CheckResponse(conn *rcon.RemoteConsole, reqId int, srv config.Server, cfg *config.Config) bool {

	resp, respReqId, err := conn.Read()

	// println(conn, respReqId, err.Error())

	if err != nil {
		if err == io.EOF {
			// return nil, err.Error()
			return false
		}
		_, _ = fmt.Fprintln(os.Stderr, "Failed to read command:", err.Error())
		return false
	}

	if reqId != respReqId {
		_, _ = fmt.Fprintln(os.Stdout, "Weird. This response is for another request.", reqId, respReqId, strings.TrimSpace(resp))
	}

	if cfg.DebugLevel > 3 {
		fmt.Println("[D4][" + srv.IP + ":" + strconv.Itoa(srv.Port) + "] RCON received (" + strings.TrimSpace(resp) + ").")
	}

	if err != nil {
		return false
	}

	return true
}
