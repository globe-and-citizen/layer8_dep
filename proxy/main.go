package main

import (
	"context"
	"fmt"
	"globe-and-citizen/layer8/proxy/config"
	"globe-and-citizen/layer8/proxy/handlers"
	"globe-and-citizen/layer8/proxy/internals/repository"
	"globe-and-citizen/layer8/proxy/internals/usecases"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func main() {

	config.InitDB()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	authServerPort := os.Getenv("AUTH_SERVER_PORT")

	authServerPortInt, err := strconv.Atoi(authServerPort)
	if err != nil {
		log.Fatal(err)
	}

	// proxyServerPort := os.Getenv("PROXY_SERVER_PORT")

	// proxyServerPortInt, err := strconv.Atoi(proxyServerPort)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// go ProxyServer(proxyServerPortInt)

	AuthServer(authServerPortInt)

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
			case path == "/init-tunnel":
				fmt.Println("Init Tunnel Triggered")
				handlers.InitTunnel(w, r)
				fmt.Println("Init Tunnel Triggered v2")
			default:
				fmt.Println("Tunnel Triggered")
				handlers.Tunnel(w, r)
			}
			log.Printf("%d %s\t%s", http.StatusOK, r.Method, r.URL.Path)
		}),
	}
	log.Fatal(server.ListenAndServe())
}

// func ProxyServer(port int) {
// 	server := http.Server{
// 		Addr: fmt.Sprintf(":%d", port),
// 		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			w.Header().Set("Access-Control-Allow-Origin", "*")
// 			w.Header().Set("Access-Control-Allow-Headers", "*")

// 			if r.Method == http.MethodOptions {
// 				w.WriteHeader(http.StatusOK)
// 				return
// 			}
// 			switch path := r.URL.Path; {
// 			case path == "/":
// 				handlers.InitTunnel(w, r)
// 			default:
// 				handlers.Tunnel(w, r)
// 			}
// 		}),
// 	}
// 	log.Printf("Starting proxy server on port %d...", port)
// 	log.Fatal(server.ListenAndServe())
// }
