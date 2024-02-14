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

/*
func Execute(hostPort string, password string, command ...string) (*string, error) {
	request_timeout := 2
	remoteConsole, err := rcon.Dial(hostPort, password, request_timeout)
	if err != nil {
		// log.Fatal("Failed to connect to RCON server", err)
		return nil, err
	}
	defer remoteConsole.Close()

	preparedCmd := strings.Join(command, " ")
	reqId, err := remoteConsole.Write(preparedCmd)

	resp, respReqId, err := remoteConsole.Read()
	if err != nil {
		if err == io.EOF {
			// return nil, err.Error()
			return nil, err
		}
		_, _ = fmt.Fprintln(os.Stderr, "Failed to read command:", err.Error())
		return nil, err
	}

	if reqId != respReqId {
		_, _ = fmt.Fprintln(os.Stdout, "Weird. This response is for another request.")
	}

	return &resp, err
}
*/

// Creates a UDP connection using the host and port.
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

// // The A2S_INFO request.
// var query = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x54, 0x53, 0x6F, 0x75, 0x72, 0x63, 0x65, 0x20, 0x45, 0x6E, 0x67, 0x69, 0x6E, 0x65, 0x20, 0x51, 0x75, 0x65, 0x72, 0x79, 0x00}

// // Creates a UDP connection using the host and port.
// func CreateConnection(host string, port int) (*net.UDPConn, error) {

// 	var UDPC *net.UDPConn

// 	// Combine host and port.
// 	fullHost := host + ":" + strconv.Itoa(port)

// 	UDPAddr, err := net.ResolveUDPAddr("udp", fullHost)

// 	if err != nil {
// 		return UDPC, err
// 	}

// 	// Attempt to open a UDP connection.
// 	UDPC, err = net.DialUDP("udp", nil, UDPAddr)

// 	if err != nil {
// 		fmt.Println(err)

// 		return UDPC, errors.New("NoConnection")
// 	}

// 	return UDPC, nil
// }

// // Sends an A2S_INFO request to the host and port specified in the arguments.
// func SendRequest(conn *net.UDPConn) {
// 	conn.Write(query)
// }

// // Checks for A2S_INFO response. Returns true if it receives a response. Returns false otherwise.
// func CheckResponse(conn *net.UDPConn, srv config.Server) bool {
// 	buffer := make([]byte, 1024)

// 	// Set read timeout.
// 	conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(srv.A2STimeout)))

// 	_, _, err := conn.ReadFromUDP(buffer)

// 	if err != nil {
// 		return false
// 	}

// 	return true
// }

// Sends an A2S_INFO request to the host and port specified in the arguments.
func SendRequest(conn *rcon.RemoteConsole) (int, error) {
	preparedCmd := "KeepAlive"
	preparedCmd = "ListPlayers"
	reqId, err := conn.Write(preparedCmd)

	return reqId, err
}

// Checks for A2S_INFO response. Returns true if it receives a response. Returns false otherwise.
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

	println(reqId, respReqId, resp)

	if reqId != respReqId {
		_, _ = fmt.Fprintln(os.Stdout, "Weird. This response is for another request.")
	}

	if cfg.DebugLevel > 3 {
		fmt.Println("[D4][" + srv.IP + ":" + strconv.Itoa(srv.Port) + "] RCON received (" + strings.TrimSpace(resp) + ").")
	}

	if err != nil {
		return false
	}

	return true
}
