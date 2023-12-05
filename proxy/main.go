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

	proxyServerPort := os.Getenv("PROXY_SERVER_PORT")

	proxyServerPortInt, err := strconv.Atoi(proxyServerPort)
	if err != nil {
		log.Fatal(err)
	}

	Server(proxyServerPortInt)

}

func Server(port int) {

	repo, err := repository.CreateRepository("postgres")
	if err != nil {
		log.Fatal(err)
	}
	usecase := &usecases.UseCase{Repo: repo}

	_, err = usecase.AddTestClient()
	if err != nil {
		log.Fatal(err)
	}

	server := http.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "*")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), "usecase", usecase))

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
			case path == "/":
				handlers.InitTunnel(w, r)
			default:
				handlers.Tunnel(w, r)
			}
		}),
	}
	log.Printf("Starting server on port %d...", port)
	log.Fatal(server.ListenAndServe())
}
