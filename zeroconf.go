package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/libp2p/zeroconf/v2"
	"golang.org/x/exp/slog"
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

func LookupService(logger *slog.Logger) []zeroconf.ServiceEntry {
	logger.Info("Looking for service instances....")
	defer logger.Info("Done")

	services := make([]zeroconf.ServiceEntry, 0)

	// Channel to receive discovered service entries
	entries := make(chan *zeroconf.ServiceEntry)

	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			log.Println("Found service", "instance_name", entry.ServiceInstanceName(), "meta", entry.Text)
			services = append(services, *entry)
		}
	}(entries)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := zeroconf.Browse(ctx, ServiceName, "local.", entries)
	if err != nil {
		logger.Error("Failed to browse", "err", err)
		os.Exit(1)
	}

	<-ctx.Done()
	return services
}
