package main

import (
	"context"
	"embed"
	"fmt"
	"globe-and-citizen/layer8/proxy/config"
	"globe-and-citizen/layer8/proxy/handlers"
	"globe-and-citizen/layer8/proxy/internals/repository"
	"globe-and-citizen/layer8/proxy/internals/usecases"
	"globe-and-citizen/layer8/proxy/resource_server/middleware"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	router "globe-and-citizen/layer8/proxy/resource_server/router"

	"github.com/joho/godotenv"
)

//go:embed dist/*
var StaticFiles embed.FS
var workingDirectory string

func indexHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	var relativePathFavicon = "dist/index.html"
	faviconPath := filepath.Join(workingDirectory, relativePathFavicon)
	fmt.Println("faviconPath: ", faviconPath)
	if r.URL.Path == "/favicon.ico" {
		http.ServeFile(w, r, faviconPath)
		return
	}
	var relativePathIndex = "/dist/index.html"
	indexPath := filepath.Join(workingDirectory, relativePathIndex)
	fmt.Println("indexPath: ", indexPath)
	http.ServeFile(w, r, indexPath)

}

func routerFunc2() http.Handler {

	mux := http.NewServeMux()

	// index
	mux.HandleFunc("/", indexHandler)

	// static files
	fmt.Println("anything?")
	staticFS, _ := fs.Sub(StaticFiles, "dist")
	httpFS := http.FileServer(http.FS(staticFS))
	// httpFS := http.FileServer(http.Dir("resource_server\\frontend\\dist"))
	mux.Handle("/assets/", httpFS)

	// api
	mux.HandleFunc("/api/v1/", middleware.LogRequest(middleware.Cors(router.RegisterRoutes())))
	return mux
}

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
