package main

import (
	"flag"
	"log"
)

const (
	programName    = "Service Discovery PoC"
	programVersion = "0.1.0"
)

func main() {
	log.Printf("Welcome to %s v%s", programName, programVersion)

	servicePort := flag.Int("port", 8080, "listening port for the service")
	flag.Parse()

	entries := LookupService()
	log.Printf("Found %d entries", len(entries))

	if len(entries) == 0 {
		log.Printf("Starting new Service instance...\n")
		go StartService(*servicePort)

		service := RegisterService(*servicePort)
		defer service.Shutdown()

	} else {
		go StartClients(entries)
	}

	select {}
}
