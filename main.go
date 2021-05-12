package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	programName    = "Service Discovery PoC"
	programVersion = "0.1.0"
)

func main() {
	fmt.Printf("%s v%s\n", programName, programVersion)

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

	done := make(chan error)

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

		done <- fmt.Errorf("%s", <-sigs)
	}()

	<-done
	fmt.Println("Bye Bye")
}
