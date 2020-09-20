package tunnel

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/someshkoli/tunnel/pkg/schema"
)

func getUID() string {
	return uuid.New().String()
}

// Tunnel - instance for tcp tunnel with remote and host
type Tunnel struct {
	Nat     Nat
	Listner net.Listener
}

// MakeTunnel - returns an instance of a tunnel
func MakeTunnel(port int) (Tunnel, error) {
	tunnel := Tunnel{}
	tunnel.Nat = Nat{}

	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return Tunnel{}, err
	}
	tunnel.Listner = ln
	return tunnel, nil
}

// Listen - Starts accepting connections
func (tunnel *Tunnel) Listen() error {
	defer tunnel.Listner.Close()
	for {
		conn, err := tunnel.Listner.Accept()
		if err != nil {
			panic(err)
		}
		// initializing first request to register the tunnel
		// sending uid with initilizing status
		message := schema.MakeRequest(schema.Initializing, getUID(), "", 0, http.Request{})
		tempBuffer := new(bytes.Buffer)
		gob.NewEncoder(tempBuffer).Encode(message)
		conn.Write(tempBuffer.Bytes())

		// receiving initial data about localhost
		bufferResponse := make([]byte, schema.MaxDataSize)
		_, err = conn.Read(bufferResponse)
		if err != nil {
			panic(err)
		}
		bufferResponseCollector := bytes.NewBuffer(bufferResponse)
		initialResponse := new(schema.Response)
		gob.NewDecoder(bufferResponseCollector).Decode(initialResponse)

		// creating new connection and registering nat for this uid
		connection := MakeConnection(initialResponse.Host, initialResponse.Port)
		connection.Conn = conn
		tunnel.Nat[initialResponse.ID] = &connection
		// running this connection through routine
		go connectionHandler(&connection)
		connection.Status = Closed
	}
}

func connectionHandler(connection *Connection) {
	for {
		select {
		case data := <-connection.send:
			fmt.Println("send this data")
			connection.Write(data)
			connection.receive <- connection.Read()
		case <-connection.stop:
			connection.Conn.Close()
			return
		}
	}
}
