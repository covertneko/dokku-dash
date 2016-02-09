// Utilities for making HTTP requests over Unix sockets
// Mostly copied from this dead PR:
// https://github.com/apatil/napping-unixsocket/blob/master/unix_socket.go
package main

import (
	"net"
	"net/http"
	"net/http/httputil"
	"path"
	"os"
	"fmt"
)

// Transport for HTTP requests over sockets
type SocketTransport struct { Path string }

// Roundtripper for unix socket requests
func (t SocketTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	dial, err := net.Dial("unix", t.Path)
	if err != nil {
		return nil, err
	}

	conn := httputil.NewClientConn(dial, nil)
	defer conn.Close()
	return conn.Do(req)
}

// Helper to test if a path identifies a unix socket
func isSocket(path string) bool {
	fi, err := os.Lstat(path)
	if err != nil {
		return false
	}

	return fi.Mode() & os.ModeType == os.ModeSocket
}

// Split a path into a socket path and request path.
// Returns an error if the path does not identify a socket.
func LocateSocket(rawPath string) (string, string, error) {
	p := rawPath
	// Ensure path is absolute
	if p[0] != '/' {
		p = "/" + p
	}

	req := ""
	req_ := ""
	for p != "" {
		// Remove trailing slash from path, if any
		if l := len(p) - 1; l >= 0 && p[l] == '/' {
			p = p[:l]
		}

		if isSocket(p) {
			return p, "/" + req, nil
		}

		// Path is not a socket. Prepend path node to request, set p to
		// remaining path.
		p, req_ = path.Split(p)
		req = path.Join(req_, req)
	}

	return "", "", fmt.Errorf("%q does not identify a socket", rawPath)
}
