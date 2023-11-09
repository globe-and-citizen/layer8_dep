package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"globe-and-citizen/layer8/resource_server/backend/middleware"
	router "globe-and-citizen/layer8/resource_server/backend/router"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	serverPort := os.Getenv("SERVER_PORT")

	// Register the routes using the RegisterRoutes() function with logger middleware
	http.HandleFunc("/api/v1/", middleware.LogRequest(middleware.Cors(router.RegisterRoutes())))

	fmt.Printf("Server listening on localhost:%s\n", serverPort)

	// Start the server on localhost and log any errors
	err = http.ListenAndServe(fmt.Sprintf(":%s", serverPort), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Server stopped")
}
