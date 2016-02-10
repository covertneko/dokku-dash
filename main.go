package main

import (
	"log"
	"fmt"
	"path"
	"os"
	"io/ioutil"
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

func Index(w http.ResponseWriter, req *http.Request) {
	req, err := http.NewRequest("GET", os.Getenv("DOKKU_API_SOCKET"), nil)
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

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading HTTP response body: ", err)
		return
	}

	fmt.Printf("Response body: %q\n", string(contents))

	fmt.Fprintf(w, string(contents))
}

func main() {
	http.HandleFunc("/", Index)
	log.Fatal(http.ListenAndServe(":" + os.Getenv("PORT"), nil))
}
