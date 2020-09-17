package client

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/someshkoli/tunnel/pkg/schema"
)

func register(conn net.Conn) {
	// TODO do registration processs herre
}

func makeLocalRequest(request schema.Request) (*http.Response, error) {
	client := http.Client{}
	// replace hostname and port and then send request
	res, err := client.Do(&request.Request)
	return res, err
}

// Client - Client for tunnel connection
type Client struct {
	ID    string
	Lport int    `json:"Lport"`
	Rport int    `json:"Rport"`
	Lhost string `json:"Lhost"`
	Rhost string `json:"Rhost"`
}

// NewClient - Returns instance of new client
func NewClient(Lhost string, Lport int, Rhost string, Rport int) *Client {
	client := Client{
		Lhost: Lhost,
		Lport: Lport,
		Rhost: Rhost,
		Rport: Rport,
	}
	return &client
}

// ConnectAndListen - connect to the host server
func (client *Client) ConnectAndListen() {
	conn, err := net.Dial("tcp", "host:"+strconv.Itoa(client.Rport))
	if err != nil {
		fmt.Println("Couldnot connect to the server")
	}
	// registers the localtunnel
	register(conn)

	// handling local tunnel requests
	go handleConnection(conn)
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		// read incomming request
		requestBuffer := make([]byte, schema.MaxDataSize)
		_, err := conn.Read(requestBuffer)
		if err != nil {
			break
		}
		bufferRequestCollector := bytes.NewBuffer(requestBuffer)
		request := new(schema.Request)
		gob.NewDecoder(bufferRequestCollector).Decode(request)

		// calling local client
		res, err := makeLocalRequest(Request)

		// writing response back to the host
		response := schema.MakeResponse(schema.Connected, request.ID, request.Rhost, request.Rport, request.Lhost, request.Lport, *res)
		messageBuffer := new(bytes.Buffer)
		gob.NewEncoder(messageBuffer).Encode(response)
		conn.Write(messageBuffer.Bytes())
	}
}
