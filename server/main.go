package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"globe-and-citizen/layer8/server/config"
	"globe-and-citizen/layer8/server/handlers"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	Ctl "globe-and-citizen/layer8/server/resource_server/controller"
	"globe-and-citizen/layer8/server/resource_server/interfaces"

	repo "globe-and-citizen/layer8/server/resource_server/repository"

	svc "globe-and-citizen/layer8/server/resource_server/service"

	"globe-and-citizen/layer8/server/internals/repository"
	OauthSvc "globe-and-citizen/layer8/server/internals/service"

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

	// Use flags to set the port
	port := flag.Int("port", 8080, "Port to run the server on")
	jwtKey := flag.String("jwtKey", "secret", "Key to sign JWT tokens")
	MpKey := flag.String("MpKey", "secret", "Key to sign mpJWT tokens")
	UpKey := flag.String("UpKey", "secret", "Key to sign upJWT tokens")
	flag.Parse()

	if *port != 8080 {
		os.Setenv("SERVER_PORT", strconv.Itoa(*port))
		os.Setenv("JWT_SECRET_KEY", *jwtKey)
		os.Setenv("MP_123_SECRET_KEY", *MpKey)
		os.Setenv("UP_999_SECRET_KEY", *UpKey)
		repository := repo.NewMemoryRepository()
		service := svc.NewService(repository)
		Server(*port, service, repository)
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	if os.Getenv("DB_USER") != "" || os.Getenv("DB_PASSWORD") != "" {
		config.InitDB()
	}

	proxyServerPort := os.Getenv("SERVER_PORT")

	proxyServerPortInt, err := strconv.Atoi(proxyServerPort)
	if err != nil {
		log.Fatal(err)
	}

	// Register repository
	repository := repo.NewRepository(config.DB)

	// Register service(usecase)
	service := svc.NewService(repository)

	Server(proxyServerPortInt, service, repository)

}

func Server(port int, service interfaces.IService, MemoryRepository interfaces.IRepository) {

	// Uncomment below line and use `repository` instead of `MemoryRepository` in `OauthService` if you want to use local postgres db
	postgresRepository := repository.InitDB()

	OauthService := &OauthSvc.Service{Repo: postgresRepository}

	_, err := OauthService.AddTestClient()
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

			r = r.WithContext(context.WithValue(r.Context(), "Oauthservice", OauthService))
			r = r.WithContext(context.WithValue(r.Context(), "service", service))

			staticFS, _ := fs.Sub(StaticFiles, "dist")
			httpFS := http.FileServer(http.FS(staticFS))

			switch path := r.URL.Path; {

			// Authorization Server endpoints
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
			case strings.HasPrefix(path, "/assets-v1"):
				http.StripPrefix("/assets-v1", http.FileServer(http.Dir("./assets-v1"))).ServeHTTP(w, r)

			// Resource Server endpoints
			case path == "/":
				Ctl.IndexHandler(w, r)
			case path == "/user":
				Ctl.UserHandler(w, r)
			case path == "/register":
				Ctl.ClientHandler(w, r)
			case path == "/api/v1/register-user":
				Ctl.RegisterUserHandler(w, r)
			case path == "/api/v1/register-client":
				Ctl.RegisterClientHandler(w, r)
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
			case path == "/init-tunnel":
				handlers.InitTunnel(w, r)
			case path == "/error":
				handlers.TestError(w, r)
			// TODO: For later, to be discussed more
			// case path == "/tunnel":
			// 	handlers.Tunnel(w, r)
			default:
				handlers.Tunnel(w, r)
			}
		}),
	}
	log.Printf("Starting server on port %d...", port)
	log.Fatal(server.ListenAndServe())
}
