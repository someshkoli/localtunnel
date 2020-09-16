package main

import (
	"fmt"

	"github.com/meta-boy/atsumaru.git/pkg/server"
)

func main() {
	server := server.NewServer()
	fmt.Println("listening to port 8000")
	if server.ListenAndServe() != nil {
		fmt.Println("error running server")
	}
}
