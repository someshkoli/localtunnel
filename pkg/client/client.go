package client

import (
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/someshkoli/tunnel/pkg/schema"
)

func register(conn net.Conn) error {
	// TODO Recieveing initial data
	fmt.Println("Recieveing initial data")
	request, err := schema.ReceiveMessage(conn)
	if err != nil {
		return err
	}
	// sending initial data
	fmt.Println("registered")

	response := schema.MakeResponse(schema.Initializing, request.ID, request.Host, request.Port, http.Response{})
	schema.SendMessage(conn, &response)
	fmt.Println(request.ID)
	return nil
}

func makeLocalRequest(request *schema.Request) (*http.Response, error) {
	client := http.Client{}
	// replace hostname and port and then send request
	res, err := client.Do(&request.Request)
	return res, err
}

// Client - Client for tunnel connection
type Client struct {
	ID   string
	Host string `json:"Host"`
	Port int    `json:"Port"`
}

// NewClient - Returns instance of new client
func NewClient(host string, port int) *Client {
	client := Client{
		Host: host,
		Port: port,
	}
	return &client
}

// ConnectAndListen - connect to the host server
func (client *Client) ConnectAndListen() {
	var done chan bool
	conn, err := net.Dial("tcp", client.Host+":"+strconv.Itoa(client.Port))
	if err != nil {
		fmt.Println(err)
		return
	}
	// registers the localtunnel
	register(conn)

	// handling local tunnel requests
	go handleConnection(conn, done)
	<-done
}

func handleConnection(conn net.Conn, done chan bool) {
	defer conn.Close()
	for {
		// read incomming request
		request, err := schema.ReceiveMessage(conn)
		if err != nil {
			fmt.Println(err)
			done <- true
		}

		// calling local client
		res, err := makeLocalRequest(request)

		// writing response back to the host
		response := schema.MakeResponse(schema.Connected, request.ID, request.Host, request.Port, *res)
		schema.SendMessage(conn, &response)

		// handle close connection state here
	}
	done <- true
}
