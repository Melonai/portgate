package main

import (
	"github.com/valyala/fasthttp"
	"log"

	"portgate"
	"portgate/handlers"
)

func main() {
	log.Print("Starting Portgate...")

	// Get global Portgate config.
	config, err := portgate.GetConfig()
	if err != nil {
		log.Fatal("Failed to get Portgate config.")
	}

	// Create handler for requests
	handler := handlers.NewRequestHandler(&config)

	// Start to listen to the outside world.
	log.Print("Listening for requests on port 8080.")
	err = fasthttp.ListenAndServe(config.PortgateAddress(), handler.HandleRequest)
	if err != nil {
		log.Fatalf("Portgate server could not be started: %s", err)
	}
}
