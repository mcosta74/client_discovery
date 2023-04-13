package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"golang.org/x/exp/slog"
)

const (
	programName    = "Service Discovery PoC"
	programVersion = "0.1.0"
)

func main() {
	opts := slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
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

	handler := opts.NewTextHandler(os.Stdout)

	logger := slog.New(handler)

	logger.Info("Started", "name", programName, "version", programVersion)
	defer logger.Info("Bye Bye")
	servicePort := flag.Int("port", 8080, "listening port for the service")
	flag.Parse()

	entries := LookupService(logger.With("component", "Lookup"))
	logger.Info(fmt.Sprintf("Found %d entries", len(entries)))

	if len(entries) == 0 {
		logger.Info("Starting new Service instance...")
		go StartService(*servicePort, logger.With("component", "Service"))

		service := RegisterService(*servicePort)
		defer service.Shutdown()

	} else {
		go StartClients(entries, logger.With("component", "Client"))
	}

	done := make(chan error)

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

		done <- fmt.Errorf("%s", <-sigs)
	}()

	<-done
}
