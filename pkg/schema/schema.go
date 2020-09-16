package schema

import "net/http"

type Status string

const (
	// Initializing the connection
	Initializing string = "Initializing"
	// Connected to the tunnel
	Connected string = "Connected"
)

// Request - incomming request
type Request struct {
	ID         string
	ClientAddr string
	Rhost      string
	Lhost      string
	Lport      int
	Rport      int
	Request    *http.Request
}

// Response - Outgoing response
type Response struct {
	ID         string
	ClientAddr string
	Rhost      string
	Lhost      string
	Lport      int
	Rport      int
	response   http.ResponseWriter
}

// MessageOut - Message to be send via connection
type MessageOut struct {
	ConnectionStatus Status
	data             *Response
}

// MessageIn - Mesaage to be recieved via connection
type MessageIn struct {
	ConnectionStatus Status
	data             *Response
}
