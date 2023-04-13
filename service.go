package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/libp2p/zeroconf/v2"
	"golang.org/x/exp/slog"
)

// The service for now "just" reply with the current time
func StartService(port int, logger *slog.Logger) {
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(rw, "Hello now is: %v", time.Now())
	})

	serverAddress := fmt.Sprintf(":%d", port)
	logger.Info("starting server", "address", serverAddress)
	if err := http.ListenAndServe(serverAddress, nil); err != nil {
		log.Fatal(err)
	}
}

var wg sync.WaitGroup

func StartClients(entries []zeroconf.ServiceEntry, logger *slog.Logger) {

	for _, entry := range entries {
		wg.Add(1)
		go startClient(entry, logger)
	}
	wg.Wait()
}

func startClient(entry zeroconf.ServiceEntry, logger *slog.Logger) {
	ticker := time.NewTicker(time.Second * 2)

	url := fmt.Sprintf("http://%v:%v", entry.AddrIPv4[0].String(), entry.Port)
	for range ticker.C {
		func() {
			resp, err := http.Get(url)
			if err != nil {
				logger.Error("failed to get response", "err", err)
				os.Exit(1)
			}
			defer resp.Body.Close()

			data, _ := ioutil.ReadAll(resp.Body)
			logger.Info("Got response", "data", data)
		}()
	}
	wg.Done()
}
