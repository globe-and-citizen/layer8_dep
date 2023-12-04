package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"globe-and-citizen/layer8/resource_server/backend/config"
	"globe-and-citizen/layer8/resource_server/backend/middleware"
	router "globe-and-citizen/layer8/resource_server/backend/router"

	"github.com/joho/godotenv"
)

//go:embed dist/*
var StaticFiles embed.FS
var workingDirectory string

func init() {
	wD, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	workingDirectory = wD
	fmt.Println("workingDirectory: ", workingDirectory)
}

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
	// http.ServeFile(w, r, "C:\\Ottawa_DT_Dev\\Learning_Computers\\layer8\\resource_server\\frontend\\dist\\index.html")
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
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	serverPort := os.Getenv("SERVER_PORT")

	srv := &http.Server{
		Addr:        ":" + serverPort,
		Handler:     routerFunc2(),
		IdleTimeout: time.Minute,
	}

	db := config.SetupDatabaseConnection()
	fmt.Println("db: ", db)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Server stopped")

	// Register the routes using the RegisterRoutes() function with logger middleware
	// http.HandleFunc("/api/v1/", middleware.LogRequest(middleware.Cors(router.RegisterRoutes())))

	// fmt.Printf("Server listening on localhost:%s\n", serverPort)

	// Start the server on localhost and log any errors
	// err = http.ListenAndServe(fmt.Sprintf(":%s", serverPort), nil)
	// if err != nil {
	// log.Fatal(err)
	// }
	// fmt.Println("Server stopped")
}
