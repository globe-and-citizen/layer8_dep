package main

import (
	"flag"
	"fmt"
	"globe-and-citizen/layer8/proxy/handlers"
	"log"
	"net/http"
)

var (
	port   = flag.Int("port", 5000, "Port to listen on")
	server = flag.String("server", "proxy", "Server type to run")
)

func main() {
	flag.Parse()

	switch *server {
	case "auth": // Errors out for now
		// AuthServer(*port)
		log.Fatal("Authentication server not yet ready. Pass `--server proxy` instead.")
	case "proxy":
		ProxyServer(*port)
	default:
		log.Fatal("Invalid server type. Valid types are: auth, proxy")
	}
}

func ProxyServer(port int) {
	server := http.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "*")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			handlers.Tunnel(w, r)
		}),
	}
	log.Printf("Starting proxy server on port %d...", port)
	log.Fatal(server.ListenAndServe())
}

// func AuthServer(port int) {

// }
