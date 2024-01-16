package main

import (
	"context"
	"embed"
	"fmt"
	"globe-and-citizen/layer8/server/config"
	"globe-and-citizen/layer8/server/handlers"
	"globe-and-citizen/layer8/server/internals/repository"
	"globe-and-citizen/layer8/server/internals/usecases"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	Ctl "globe-and-citizen/layer8/server/resource_server/controller"

	"github.com/joho/godotenv"
)

// go:embed dist
var StaticFiles embed.FS

var workingDirectory string

func getPwd() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	workingDirectory = dir
}

func main() {

	config.InitDB()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	proxyServerPort := os.Getenv("SERVER_PORT")

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

	getPwd()

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

			staticFS, _ := fs.Sub(StaticFiles, "dist")
			httpFS := http.FileServer(http.FS(staticFS))

			switch path := r.URL.Path; {

			// Authorization Server endpoints
			case path == "/login":
				handlers.Login(w, r)
			// case path == "/register":
			// 	handlers.Register(w, r)
			case path == "/authorize":
				handlers.Authorize(w, r)
			case path == "/error":
				handlers.Error(w, r)
			case path == "/api/oauth":
				handlers.OAuthToken(w, r)
			case path == "/api/user":
				handlers.UserInfo(w, r)
			case strings.HasPrefix(path, "/assets-v1"):
				http.StripPrefix("/assets-v1", http.FileServer(http.Dir("./assets-v1"))).ServeHTTP(w, r)

			// Resource Server endpoints
			case path == "/api/v1/register-user":
				Ctl.RegisterUserHandler(w, r)
			case path == "/api/v1/register-client":
				Ctl.RegisterClientHandler(w, r)
			case path == "/register":
				Ctl.ClientHandler(w, r)
			case path == "/api/v1/getClient":
				Ctl.GetClientData(w, r)
			case path == "/api/v1/login-precheck":
				Ctl.LoginPrecheckHandler(w, r)
			case path == "/api/v1/login-user":
				Ctl.LoginUserHandler(w, r)
			case path == "/api/v1/profile":
				Ctl.ProfileHandler(w, r)
			case path == "/api/v1/verify-email":
				Ctl.VerifyEmailHandler(w, r)
			case path == "/api/v1/change-display-name":
				Ctl.UpdateDisplayNameHandler(w, r)
			case path == "/favicon.ico":
				faviconPath := workingDirectory + "/dist/favicon.ico"
				http.ServeFile(w, r, faviconPath)
			case strings.HasPrefix(path, "/assets/"):
				httpFS.ServeHTTP(w, r)

			// Proxy Server endpoints
			case path == "/":
				Ctl.IndexHandler(w, r)
			case path == "/user":
				Ctl.UserHandler(w, r)
			case path == "/init-tunnel":
				handlers.InitTunnel(w, r)
			case path == "/error":
				handlers.TestError(w, r)
			default:
				handlers.Tunnel(w, r)
			}
		}),
	}
	log.Printf("Starting server on port %d...", port)
	log.Fatal(server.ListenAndServe())
}
