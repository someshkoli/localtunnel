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

func send(conn net.Conn, res *schema.Response) {
	messageBuffer := new(bytes.Buffer)
	gob.NewEncoder(messageBuffer).Encode(res)
	conn.Write(messageBuffer.Bytes())
	fmt.Println("messageSent")
}

func recieve(conn net.Conn) (*schema.Request, error) {
	requestBuffer := make([]byte, schema.MaxDataSize)
	_, err := conn.Read(requestBuffer)
	if err != nil {
		fmt.Println(err)
		return &schema.Request{}, err
	}
	bufferRequestCollector := bytes.NewBuffer(requestBuffer)
	request := new(schema.Request)
	gob.NewDecoder(bufferRequestCollector).Decode(request)
	return request, nil
}

func register(conn net.Conn) error {
	// TODO Recieveing initial data
	fmt.Println("Recieveing initial data")
	request, err := recieve(conn)
	if err != nil {
		return err
	}
	// sending initial data
	fmt.Println("registered")

	response := schema.MakeResponse(schema.Initializing, request.ID, request.Host, request.Port, http.Response{})
	send(conn, &response)
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
	host string `json:"host"`
	port int    `json:"host"`
}

// NewClient - Returns instance of new client
func NewClient(host string, port int) *Client {
	client := Client{
		host: host,
		port: port,
	}
	return &client
}

// ConnectAndListen - connect to the host server
func (client *Client) ConnectAndListen() {
	var done chan bool
	conn, err := net.Dial("tcp", client.host+":"+strconv.Itoa(client.port))
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
		request, err := recieve(conn)
		if err != nil {
			fmt.Println(err)
			done <- true
		}

		// calling local client
		res, err := makeLocalRequest(request)

		// writing response back to the host
		response := schema.MakeResponse(schema.Connected, request.ID, request.Host, request.Port, *res)
		send(conn, &response)
	}
	done <- true
}
