/*
A Go written library for the RCON Protocol from Valve.

Information to the protocol can be found under:
https://developer.valvesoftware.com/wiki/Source_RCON_Protocol

This is a fork from https://github.com/james4k/rcon with the support for go
modules and with a rework of the original implementation for better readability.
*/
package rcon

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
	"sync"
	"time"
)

const (
	typeAuth          = 3
	typeExecCommand   = 2
	typeResponseValue = 0
	typeAuthResponse  = 2

	fieldPackageSize = 4
	fieldIDSize      = 4
	fieldTypeSize    = 4
	fieldMinBodySize = 1
	fieldEndSize     = 1
)

// The minimum package size contains:
// 4 bytes for the ID field
// 4 bytes for the Type field
// 1 byte minimum for an empty body string
// 1 byte for the empty string at the end
//
// https://developer.valvesoftware.com/wiki/Source_RCON_Protocol#Packet_Size
// The 4 bytes representing the size of the package are not included.
const minPackageSize = fieldIDSize + fieldTypeSize + fieldMinBodySize + fieldEndSize

// maxPackageSize of a request/response package.
// This size does not include the size field.
// https://developer.valvesoftware.com/wiki/Source_RCON_Protocol#Packet_Size
const maxPackageSize = 4096

// RemoteConsole holds the information to communicate withe remote console.
type RemoteConsole struct {
	conn       net.Conn
	readBuff   []byte
	readMutex  sync.Mutex
	queuedBuff []byte
}

var (
	// ErrAuthFailed the authentication against the server failed.
	// This happens if the request id doesn't match the response id.
	ErrAuthFailed = errors.New("rcon: authentication failed")

	// ErrInvalidAuthResponse the response of an authentication request doesn't match the correct type.
	ErrInvalidAuthResponse = errors.New("rcon: invalid response type during auth")

	// ErrUnexpectedFormat the response package is not correctly formatted.
	ErrUnexpectedFormat = errors.New("rcon: unexpected response format")

	// ErrCommandTooLong the command is bigger than the bodyBufferSize.
	ErrCommandTooLong = errors.New("rcon: command too long")

	// ErrResponseTooLong the response package is bigger than the maxPackageSize.
	ErrResponseTooLong = errors.New("rcon: response too long")
)

// Dial establishes a connection with the remote server.
// It can return multiple errors:
// 	- ErrInvalidAuthResponse
// 	- ErrAuthFailed
// 	- and other types of connection errors that are not specified in this package.
func Dial(host, password string, timeout_in_seconds ...int) (*RemoteConsole, error) {
	timeout := 10 * time.Second
	if len(timeout_in_seconds) > 0 {
		timeout = time.Duration(timeout_in_seconds[0]) * time.Second
	}
	conn, err := net.DialTimeout("tcp", host, timeout)
	if err != nil {
		return nil, err
	}

	remoteConsole := &RemoteConsole{
		conn:     conn,
		readBuff: make([]byte, maxPackageSize+fieldPackageSize),
	}
	remoteConsole.auth(password, timeout)
	if err != nil {
		return nil, err
	}

	return remoteConsole, nil
}

// LocalAddr returns the local network address.
func (r *RemoteConsole) LocalAddr() net.Addr {
	return r.conn.LocalAddr()
}

// RemoteAddr returns the remote network address.
func (r *RemoteConsole) RemoteAddr() net.Addr {
	return r.conn.RemoteAddr()
}

// Write a command to the server.
//
// It can return ErrCommandTooLong if the given cmd str is too long.
// Additionally it can return any other connection related errors.
func (r *RemoteConsole) Write(cmd string) (requestID int, err error) {
	requestID = int(newRequestID())
	err = r.writeCmd(int32(requestID), typeExecCommand, cmd)

	return
}

// Read a incoming response from the server.
// If the response doesn't contain the correct ResponseValue it will return a response with a empty string an the request id = 0.
// This is also the case if an error happens even though the error will be returned.
//
// It can return following errors:
//  - ErrResponseTooLong
//  - ErrUnexpectedFormat
//  - or a connection error that isn't typed in this package
func (r *RemoteConsole) Read() (response string, requestID int, err error) {
	var respType int
	var respBytes []byte
	respType, requestID, respBytes, err = r.readResponse(2 * time.Minute)
	if err != nil || respType != typeResponseValue {
		response = ""
		requestID = 0
	} else {
		response = string(respBytes)
	}
	return
}

// Close the connection to the server.
func (r *RemoteConsole) Close() error {
	return r.conn.Close()
}

func newRequestID() int32 {
	return int32((time.Now().UnixNano() / 100000) % 100000)
}

func (r *RemoteConsole) auth(password string, timeout time.Duration) error {
	reqID := newRequestID()
	err := r.writeCmd(reqID, typeAuth, password)
	if err != nil {
		return err
	}

	respType, responseID, _, err := r.readResponse(timeout)
	if err != nil {
		return err
	}

	// if we didn't get an auth response back, try again. it is often a bug
	// with RCON servers that you get an empty response before receiving the
	// auth response.
	if respType != typeAuthResponse {
		respType, responseID, _, err = r.readResponse(timeout)
	}
	if err != nil {
		return err
	}
	if respType != typeAuthResponse {
		return ErrInvalidAuthResponse
	}
	if responseID != int(reqID) {
		return ErrAuthFailed
	}

	return nil
}

func (r *RemoteConsole) writeCmd(reqID, pkgType int32, cmd string) error {
	if len(cmd) > maxPackageSize-minPackageSize {
		return ErrCommandTooLong
	}

	buffer := bytes.NewBuffer(make([]byte, 0, minPackageSize+fieldPackageSize+len(cmd)))

	// packet size
	binary.Write(buffer, binary.LittleEndian, int32(minPackageSize+len(cmd)))

	// request id
	binary.Write(buffer, binary.LittleEndian, int32(reqID))

	// type of the package
	binary.Write(buffer, binary.LittleEndian, int32(pkgType))

	// body
	buffer.WriteString(cmd)

	// double null termination
	binary.Write(buffer, binary.LittleEndian, byte(0))
	binary.Write(buffer, binary.LittleEndian, byte(0))

	r.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err := r.conn.Write(buffer.Bytes())

	return err
}

func (r *RemoteConsole) readResponse(timeout time.Duration) (int, int, []byte, error) {
	r.readMutex.Lock()
	defer r.readMutex.Unlock()

	r.conn.SetReadDeadline(time.Now().Add(timeout))
	var readBytes int
	var err error
	if r.queuedBuff != nil {
		copy(r.readBuff, r.queuedBuff)
		readBytes = len(r.queuedBuff)
		r.queuedBuff = nil
	} else {
		readBytes, err = r.conn.Read(r.readBuff)
		if err != nil {
			return 0, 0, nil, err
		}
	}

	dataSize, readBytes, err := r.readResponsePackageSize(readBytes)
	if err != nil {
		return 0, 0, nil, err
	}

	if dataSize > maxPackageSize {
		return 0, 0, nil, ErrResponseTooLong
	}

	totalPackageSize := dataSize + fieldPackageSize
	readBytes, err = r.readResponsePackage(totalPackageSize, readBytes)
	if err != nil {
		return 0, 0, nil, err
	}

	// The data has to be explicitly selected to prevent copying empty bytes.
	data := r.readBuff[fieldPackageSize:totalPackageSize]

	// Save not package related bytes for the next read.
	if readBytes > totalPackageSize {
		// start of the next buffer was at the end of this packet.
		// save it for the next read.
		// The data has to be explicitly selected to prevent copying empty bytes.
		r.queuedBuff = r.readBuff[totalPackageSize:readBytes]
	}

	return r.readResponseData(data)
}

// readResponsePackageSize wait until first 4 bytes are read to get the package size.
// Takes as param how many bytes are already read. The returned size does not include the size field.
func (r *RemoteConsole) readResponsePackageSize(readBytes int) (int, int, error) {
	for readBytes < fieldPackageSize {
		// need the 4 byte packet size...
		b, err := r.conn.Read(r.readBuff[readBytes:])
		if err != nil {
			return 0, 0, err
		}
		readBytes += b
	}

	var size int32
	b := bytes.NewBuffer(r.readBuff[:fieldPackageSize])
	err := binary.Read(b, binary.LittleEndian, &size)
	if err != nil {
		return 0, 0, err
	}

	if size < minPackageSize {
		return 0, 0, ErrUnexpectedFormat
	}

	return int(size), readBytes, nil
}

// readResponsePackage waits until the whole package is read including the size field.
func (r *RemoteConsole) readResponsePackage(totalPackageSize, readBytes int) (int, error) {
	for totalPackageSize > readBytes {
		b, err := r.conn.Read(r.readBuff[readBytes:])
		if err != nil {
			return readBytes, err
		}
		readBytes += b
	}

	return readBytes, nil
}

func (r *RemoteConsole) readResponseData(data []byte) (int, int, []byte, error) {
	var requestID, responseType int32
	buffer := bytes.NewBuffer(data)

	err := binary.Read(buffer, binary.LittleEndian, &requestID)
	if err != nil {
		return 0, 0, []byte{}, err
	}

	binary.Read(buffer, binary.LittleEndian, &responseType)
	if err != nil {
		return 0, 0, []byte{}, err
	}

	// the rest of the buffer is the body.
	body := buffer.Bytes()
	// remove the to null terminations
	body = body[:len(body)-2]

	return int(responseType), int(requestID), body, nil
}
