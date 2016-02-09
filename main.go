package main

import (
	"log"
	"fmt"
	"path"
	"os"
	"net/http"
)

// Construct and send an HTTP request
func Send(req *http.Request) (response *http.Response, err error) {
	var client *http.Client

	// Convert unix socket request to HTTP request
	if req.URL.Scheme == "unix" {
		var socketPath string
		var reqPath string
		socketPath, reqPath, err = LocateSocket(path.Join(req.URL.Host, req.URL.Path))

		if err != nil {
			log.Println(err)
			return
		}

		req.URL.Path = reqPath
		req.URL.Host = ""
		req.URL.Scheme = "http"

		client = &http.Client{Transport: SocketTransport{Path: socketPath}}
	} else {
		client = &http.Client{}
	}

	response, err = client.Do(req)
	return
}

func Index(w http.ResponseWriter, req *http.Request)  {
	req, err := http.NewRequest("GET", "unix:///tmp/dokku-api.sock/", nil)
	if err != nil {
		log.Fatal("Could not construct HTTP request: ", err)
		return
	}

	resp, err := Send(req)
	if err != nil {
		log.Fatal("Could not send HTTP request: ", err)
		return
	}

	defer resp.Body.Close()

	var contents []byte
	resp.Body.Read(contents)

	fmt.Fprintf(w, string(contents))
}

func main() {
	http.HandleFunc("/", Index)
	log.Fatal(http.ListenAndServe(":" + os.Getenv("PORT"), nil))
}
