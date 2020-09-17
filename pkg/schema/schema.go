package schema

import (
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
	Lhost            string
	Lport            int
	Rport            int
	Request          http.Request
}

// Response - Outgoing response
type Response struct {
	ConnectionStatus Status
	ID               string
	ClientAddr       string
	Rhost            string
	Lhost            string
	Lport            int
	Rport            int
	Response         http.Response
}

// MakeRequest - Returns instance for a request
func MakeRequest(
	status Status,
	id string,
	rhost string,
	rport int,
	lhost string,
	lport int,
	req http.Request) Request {
	request := Request{}
	request.ConnectionStatus = status
	request.ID = id
	request.Rhost = rhost
	request.Lhost = lhost
	request.Rport = rport
	request.Lport = lport
	request.Request = req
	return request
}

// MakeResponse - Creates new Instance for
func MakeResponse(
	status Status,
	id string,
	rhost string,
	rport int,
	lhost string,
	lport int,
	res http.Response) Response {
	response := Response{}
	response.ConnectionStatus = status
	response.ID = id
	response.Lhost = lhost
	response.Lport = lport
	response.Rhost = rhost
	response.Rport = rport
	response.Response = res
	return response
}
