package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/libp2p/zeroconf/v2"
)

// The service for now "just" reply with the current time
func StartService(port int) {
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(rw, "Hello now is: %v", time.Now())
	})

	serverAddress := fmt.Sprintf(":%d", port)
	log.Printf("starting server at %s...\n", serverAddress)
	if err := http.ListenAndServe(serverAddress, nil); err != nil {
		log.Fatal(err)
	}
}

var wg sync.WaitGroup

func StartClients(entries []zeroconf.ServiceEntry) {

	for _, entry := range entries {
		wg.Add(1)
		go startClient(entry)
	}
	wg.Wait()
}

func startClient(entry zeroconf.ServiceEntry) {
	ticker := time.NewTicker(time.Second * 2)

	url := fmt.Sprintf("http://%v:%v", entry.AddrIPv4[0].String(), entry.Port)
	for range ticker.C {
		func() {
			resp, err := http.Get(url)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()

			data, _ := ioutil.ReadAll(resp.Body)
			log.Printf("Got response: %s", data)
		}()
	}
	wg.Done()
}
