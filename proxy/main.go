package main

import (
	"context"
	"flag"
	"fmt"
	"globe-and-citizen/layer8/proxy/config"
	"globe-and-citizen/layer8/proxy/handlers"
	"globe-and-citizen/layer8/proxy/internals/repository"
	"globe-and-citizen/layer8/proxy/internals/usecases"
	"log"
	"net/http"
	"strings"
)

var (
	port   = flag.Int("port", 5000, "Port to listen on")
	server = flag.String("server", "auth", "Server type to run")
)

func main() {
	config.InitDB()

	flag.Parse()

	switch *server {
	case "auth":
		AuthServer(*port)
	case "proxy":
		ProxyServer(*port)
	default:
		log.Fatal("Invalid server type. Valid types are: auth, proxy")
	}
}

func AuthServer(port int) {
	log.Printf("Starting auth server on port %d...", port)

	// intialize a postgres connection and a usecase
	repo, err := repository.CreateRepository("postgres")
	if err != nil {
		log.Fatal(err)
	}
	usecase := &usecases.UseCase{Repo: repo}

	// for testing purposes, we'll create a test client for the example service provider
	_, err = usecase.AddTestClient()
	if err != nil {
		log.Fatal(err)
	}

	server := http.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// allow from all origins
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "*")

			// add the usecase to context
			r = r.WithContext(context.WithValue(r.Context(), "usecase", usecase))
			// routing
			switch path := r.URL.Path; {
			// case path == "" || path == "/":
			// 	handlers.Welcome(w, r)
			case path == "/login":
				handlers.Login(w, r)
			case path == "/authorize":
				handlers.Authorize(w, r)
			case path == "/error":
				handlers.Error(w, r)
			case path == "/api/oauth":
				handlers.OAuthToken(w, r)
			case path == "/api/user":
				handlers.UserInfo(w, r)
			case strings.HasPrefix(path, "/assets"):
				// serve static files
				http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))).ServeHTTP(w, r)
			default:
				http.Error(w, "Invalid path", http.StatusNotFound)
			}
			log.Printf("%d %s\t%s", http.StatusOK, r.Method, r.URL.Path)
		}),
	}
	log.Fatal(server.ListenAndServe())
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
