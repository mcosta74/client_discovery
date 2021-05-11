package main

import (
	"context"
	"log"
	"time"

	"github.com/grandcat/zeroconf"
)

const (
	ServiceName = "_mytimeserver._tcp"
)

func RegisterService(port int) *zeroconf.Server {
	meta := []string{
		"version=0.1.0",
		"hello=world",
	}

	service, err := zeroconf.Register(
		"My Time Server",
		ServiceName,
		"local.",
		port,
		meta,
		nil,
	)

	if err != nil {
		log.Fatal(err)
	}

	return service
}

func LookupService() []zeroconf.ServiceEntry {
	log.Println("Looking for service instances....")
	defer log.Println("Done")

	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatal(err)
	}

	services := make([]zeroconf.ServiceEntry, 0)

	// Channel to receive discovered service entries
	entries := make(chan *zeroconf.ServiceEntry)

	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			log.Println("Found service:", entry.ServiceInstanceName(), entry.Text)
			services = append(services, *entry)
		}
	}(entries)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err = resolver.Browse(ctx, ServiceName, "local.", entries)
	if err != nil {
		log.Fatalln("Failed to browse:", err.Error())
	}

	<-ctx.Done()
	return services
}
