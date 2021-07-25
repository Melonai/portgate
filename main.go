package main

import (
	"github.com/valyala/fasthttp"
	"log"
)

func main() {
	log.Print("Starting Portgate...")

	// Get global Portgate config.
	config, err := GetConfig()
	if err != nil {
		log.Fatal("Failed to get Portgate config.")
	}

	// Create handler for requests
	handler := RequestHandler{
		config: &config,
		client: fasthttp.Client{},
	}

	// Start to listen to the outside world.
	log.Print("Listening for requests on port 8080.")
	err = fasthttp.ListenAndServe(config.PortgateAddress(), handler.handleRequest)
	if err != nil {
		log.Fatalf("Portgate server could not be started: %s", err)
	}
}
