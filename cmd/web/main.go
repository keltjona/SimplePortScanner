package main

import (
	"flag"
	"fmt"
	"os"
	
	"github.com/user/portscanner/web/server"
)

func main() {
	// Parse command line flags
	port := flag.Int("port", 8080, "Port to run the web server on")
	
	flag.Parse()
	
	// Start the web server
	err := server.StartServer(*port)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}
}
