package proxy

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func temphandle(res http.ResponseWriter, req *http.Request) {
	url := req.URL
	path := url.Path
	fmt.Println(path)
	res.Write([]byte(path))
}

// NewServer returns new server instance
func NewServer(port int) *http.Server {
	// http server to handle client request
	mux := mux.NewRouter()
	mux.HandleFunc("{*}", temphandle)
	server := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: mux,
	}
	return server
}
