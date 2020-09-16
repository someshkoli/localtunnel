package proxy

import (
	"net/http"

	"github.com/gorilla/mux"
)

func getRootMux() *mux.Router {
	mux := mux.NewRouter()
	mux.HandleFunc("", temphandle)
	mux.HandleFunc("/{*}", temphandle)
	return mux
}

// NewServer returns new server instance
func NewServer() *http.Server {
	server := &http.Server{
		Addr:    ":8000",
		Handler: getRootMux(),
	}
	return server
}
