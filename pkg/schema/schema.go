package schema

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"net/http"
)

// Status - Connection status
type Status string

const (
	// Initializing the connection
	Initializing Status = "Initializing"
	// Connected to the tunnel
	Connected Status = "Connected"
	// MaxDataSize - Maximum data transfer size
	MaxDataSize = 8192
)

// Request - incomming request
type Request struct {
	ConnectionStatus Status
	ID               string
	Rhost            string
	Host             string
	Port             int
	Request          http.Request
}

// Response - Outgoing response
type Response struct {
	ConnectionStatus Status
	ID               string
	ClientAddr       string
	Host             string
	Port             int
	Response         http.Response
}

// MakeRequest - Returns instance for a request
func MakeRequest(
	status Status,
	id string,
	host string,
	port int,
	req http.Request) Request {
	request := Request{}
	request.ConnectionStatus = status
	request.ID = id
	request.Host = host
	request.Port = port
	request.Request = req
	return request
}

// MakeResponse - Creates new Instance for
func MakeResponse(
	status Status,
	id string,
	host string,
	port int,
	res http.Response) Response {
	response := Response{}
	response.ConnectionStatus = status
	response.ID = id
	response.Host = host
	response.Port = port
	response.Response = res
	return response
}

// SendMessage - Sends message via tcp connection
func SendMessage(conn net.Conn, res *Response) {
	messageBuffer := new(bytes.Buffer)
	gob.NewEncoder(messageBuffer).Encode(res)
	conn.Write(messageBuffer.Bytes())
	fmt.Println("messageSent")
}

// ReceiveMessage - Recieves single incomming message tcp connection
func ReceiveMessage(conn net.Conn) (*Request, error) {
	requestBuffer := make([]byte, MaxDataSize)
	_, err := conn.Read(requestBuffer)
	if err != nil {
		fmt.Println(err)
		return &Request{}, err
	}
	bufferRequestCollector := bytes.NewBuffer(requestBuffer)
	request := new(Request)
	gob.NewDecoder(bufferRequestCollector).Decode(request)
	return request, nil
}
