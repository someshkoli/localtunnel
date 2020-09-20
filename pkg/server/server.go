package server

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/someshkoli/tunnel/pkg/server/tunnel"
)

// Server - Server struct
type Server struct {
	// port number for the proxy server
	ProxyPort int `json:"ProxyPort"`
	// port number for the tunneling tcp server
	TunnelPort int `json:"TunnelPort"`
	tunnel     *tunnel.Tunnel
}

// NewServer - returns new Server instance
func NewServer(proxyPort int, tunnelPort int) *Server {
	return &Server{
		ProxyPort:  proxyPort,
		TunnelPort: tunnelPort,
	}
}

// StartServer - Starts localtunnel server
func (server *Server) StartServer() error {
	// starting tcp server
	tunnel, err := tunnel.MakeTunnel(server.TunnelPort)
	if err != nil {
		fmt.Println(err)
		return err
	}
	server.tunnel = &tunnel
	go server.tunnel.Listen()

	// http server to handle client request
	mux := mux.NewRouter()
	mux.HandleFunc("{*}", server.temphandle)
	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(server.ProxyPort),
		Handler: mux,
	}
	if err := srv.ListenAndServe(); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (server *Server) temphandle(res http.ResponseWriter, req *http.Request) {
	newRequest, id, ok := server.getModifiedRequest(req)

}

// Returns modified http request with remote id and error if any
func (server *Server) getModifiedRequest(request *http.Request) (*http.Request, string, bool) {
	oldURL := request.URL
	path := oldURL.Path
	hostName := oldURL.Hostname()
	hostnameSlice := strings.Split(hostName, ".")
	ID := hostnameSlice[0]

	newURL := &url.URL{}
	newURL.Path = path
}
