package query

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/gamemann/Pterodactyl-Game-Server-Watch/config"
)

// The A2S_INFO request.
var query = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x54, 0x53, 0x6F, 0x75, 0x72, 0x63, 0x65, 0x20, 0x45, 0x6E, 0x67, 0x69, 0x6E, 0x65, 0x20, 0x51, 0x75, 0x65, 0x72, 0x79, 0x00}

// Creates a UDP connection using the host and port.
func CreateConnection(host string, port int) (*net.UDPConn, error) {
	var UDPC *net.UDPConn

	// Combine host and port.
	fullHost := host + ":" + strconv.Itoa(port)

	UDPAddr, err := net.ResolveUDPAddr("udp", fullHost)

	if err != nil {
		return UDPC, err
	}

	// Attempt to open a UDP connection.
	UDPC, err = net.DialUDP("udp", nil, UDPAddr)

	if err != nil {
		fmt.Println(err)

		return UDPC, errors.New("NoConnection")
	}

	return UDPC, nil
}

// Sends an A2S_INFO request to the host and port specified in the arguments.
func SendRequest(conn *net.UDPConn) {
	conn.Write(query)
}

// Checks for A2S_INFO response. Returns true if it receives a response. Returns false otherwise.
func CheckResponse(conn *net.UDPConn, srv config.Server) bool {
	buffer := make([]byte, 1024)

	// Set read timeout.
	conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(srv.A2STimeout)))

	_, _, err := conn.ReadFromUDP(buffer)

	if err != nil {
		return false
	}

	return true
}
