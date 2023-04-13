package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"golang.org/x/exp/slog"

	"github.com/mcosta74/service_discovery"
)

const (
	programName    = "Service Discovery PoC"
	programVersion = "0.1.0"
)

func main() {
	config := service_discovery.LoadConfig(os.Args[1:])

	logger := setupLogger(config)
	logger.Info("Started", "name", programName, "version", programVersion)
	defer logger.Info("Bye Bye")

	entries := service_discovery.LookupService(logger.With("component", "Lookup"))
	logger.Info(fmt.Sprintf("Found %d entries", len(entries)))

	if len(entries) == 0 {
		logger.Info("Starting new Service instance...")
		go service_discovery.StartService(config.HTTPPort, logger.With("component", "Service"))

		service := service_discovery.RegisterService(config.HTTPPort)
		defer service.Shutdown()

	} else {
		go service_discovery.StartClients(entries, logger.With("component", "Client"))
	}

	done := make(chan error)

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

		done <- fmt.Errorf("%s", <-sigs)
	}()

	<-done
}

func setupLogger(config service_discovery.Config) *slog.Logger {
	opts := slog.HandlerOptions{
		AddSource: true,
		Level:     config.LogLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.TimeKey:
				a.Value = slog.TimeValue(a.Value.Time().UTC())

			case slog.SourceKey:
				a.Value = slog.StringValue(filepath.Base(a.Value.String()))
			}
			return a
		},
	}

	var handler slog.Handler = opts.NewJSONHandler(os.Stdout)
	if config.LogFormat == "text" {
		handler = opts.NewTextHandler(os.Stdout)
	}

	return slog.New(handler)
}
