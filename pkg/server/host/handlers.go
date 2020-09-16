package proxy

import (
	"fmt"
	"net/http"
)

func temphandle(res http.ResponseWriter, req *http.Request) {
	url := req.URL
	path := url.Path
	fmt.Println(path)
	res.Write([]byte(path))
}
