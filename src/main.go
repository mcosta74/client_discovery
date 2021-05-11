package main

import (
	"log"
)

const (
	ProgramName    = "Service Discovery PoC"
	ProgramVersion = "0.1.0"
	ServicePort    = 8080
)

func main() {
	log.Printf("Welcome to %s v%s", ProgramName, ProgramVersion)

	entries := LookupService()
	log.Printf("Found %d entries", len(entries))

	if len(entries) == 0 {
		log.Printf("Starting new Service instance...\n")
		go StartService(ServicePort)

		service := RegisterService(ServicePort)
		defer service.Shutdown()

	} else {
		go StartClients(entries)
	}

	select {}
}
