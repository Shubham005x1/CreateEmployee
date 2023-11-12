package main

import (
	"log"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
)

// main function sets up the HTTP server and routes.
func main() {
	// Create a new router using gorilla/mux.
	//router := mux.NewRouter()

	// // Initialize Firestore connection.
	// initializeFirestore()

	// // Define a route that listens for POST requests on the "/employees" endpoint and
	// // calls the CreateEmployee function to handle the request.

	// // Start the server and listen on port :8080.
	// log.Println("Server started on :8080")
	// log.Fatal(http.ListenAndServe(":8080", router))
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	// By default, listen on all interfaces. If testing locally, run with
	// LOCAL_ONLY=true to avoid triggering firewall warnings and
	// exposing the server outside of your own machine.
	hostname := ""
	if localOnly := os.Getenv("LOCAL_ONLY"); localOnly == "true" {
		hostname = "127.0.0.1"
	}
	if err := funcframework.StartHostPort(hostname, port); err != nil {
		log.Fatalf("funcframework.StartHostPort: %v\n", err)
	}
}
