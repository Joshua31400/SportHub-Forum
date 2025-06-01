package main

import (
	"SportHub-Forum/internal/server"
	"log"
)

func main() {
	// Initialise and start the server
	srv := server.New()
	log.Fatal(srv.Start())
}
