package handlers

import (
	"io"
	"net"
	"net/http"
)

// Tunnel is a handler that redirects TCP streams to the target host
func Tunnel(w http.ResponseWriter, r *http.Request) {
	// get the target host
	target := r.URL.Host
	// connect to the target host
	targetConn, err := net.Dial("tcp", target)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	// send a 200 status code to the client
	w.WriteHeader(http.StatusOK)
	// flush the response
	w.(http.Flusher).Flush()
	// copy data from the client to the target host
	go func() {
		defer targetConn.Close()
		defer r.Body.Close()
		io.Copy(targetConn, r.Body)
	}()
	// copy data from the target host to the client
	defer targetConn.Close()
	defer r.Body.Close()
	io.Copy(w, targetConn)
}
