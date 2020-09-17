package tunnel

import (
	"bytes"
	"encoding/gob"
	"net"

	"github.com/google/uuid"
	"github.com/someshkoli/tunnel/pkg/schema"
)

// ConnectionStatus - runnning status of the tunnel connection
type ConnectionStatus string

const (
	// Open - Refers to open tunnel connection
	Open ConnectionStatus = "Open"
	// Closed - Refers to closed tunnel connection
	Closed ConnectionStatus = "Closed"
)

// Connection - instance for new tunnel connection
type Connection struct {
	ID      string   `json:"ID"`
	Lhost   string   `json:"Lhost"`
	Lport   int      `json:"Lport"`
	Rhost   string   `json:"Rhost"`
	Rport   int      `json:"Rport"`
	Conn    net.Conn `json:"Conn"`
	send    chan *schema.Request
	receive chan schema.Response
	stop    chan bool
	Status  ConnectionStatus
}

// MakeConnection - Returns new instance for the connection
func MakeConnection(lhost string, lport int, rhost string, rport int) Connection {
	connection := Connection{
		ID:    uuid.New().String(),
		Rhost: rhost,
		Rport: rport,
		Lhost: lhost,
		Lport: lport,
	}
	connection.Status = Open
	connection.send = make(chan *schema.Request)
	connection.receive = make(chan schema.Response)
	connection.stop = make(chan bool)
	return connection
}

// Write - Writes data to the connection
func (connection *Connection) Write(message *schema.Request) error {
	messageBuffer := new(bytes.Buffer)
	gob.NewEncoder(messageBuffer).Encode(message)
	_, err := connection.Conn.Write(messageBuffer.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func (connection *Connection) Read() schema.Response {
	bufferResponse := make([]byte, schema.MaxDataSize)
	_, _ = connection.Conn.Read(bufferResponse)
	bufferResposeCollector := bytes.NewBuffer(bufferResponse)
	tempResponse := new(schema.Response)
	gob.NewDecoder(bufferResposeCollector).Decode(tempResponse)
	return *tempResponse
}
