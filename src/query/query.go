package query

import (
	"errors"
	"fmt"
	"net"
)

// The A2S_INFO request.
var query = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x54, 0x53, 0x6F, 0x75, 0x72, 0x63, 0x65, 0x20, 0x45, 0x6E, 0x67, 0x69, 0x6E, 0x65, 0x20, 0x51, 0x75, 0x65, 0x72, 0x79, 0x00}

// Creates a UDP connection using the host and port.
func CreateConnection(host string, port int) (*net.UDPConn, error) {
	var UDPC *net.UDPConn

	// Combine host and port.
	fullHost := host + ":" + string(port)

	// Create dialer with one second timeout.
	d := net.Dialer{Timeout: 1}

	// Attempt to open a UDP connection.
	conn, err := d.Dial("udp4", fullHost)

	if err != nil {
		fmt.Println(err)

		return UDPC, errors.New("NoConnection")
	}

	UDPC, _ = conn.(*net.UDPConn)

	return UDPC, nil
}

// Sends an A2S_INFO request to the host and port specified in the arguments.
func SendRequest(conn *net.UDPConn) {
	conn.Write(query)
}

// Checks for A2S_INFO response. Returns true if it receives a response. Returns false otherwise.
func CheckResponse(conn *net.UDPConn) bool {
	buffer := make([]byte, 256)

	_, _, err := conn.ReadFromUDP(buffer)

	if err != nil {
		return false
	}

	return true
}
