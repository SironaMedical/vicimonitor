package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/strongswan/govici/vici"
)

// Set flags
var pollInterval int
var initiateTimeout int
var listenerString string
var socketPath string

func init() {
	flag.IntVar(&pollInterval, "i", 5, "poll interval duration in seconds for the vici socket")
	flag.IntVar(&initiateTimeout, "t", 3000, "SA initiate timeout")
	flag.StringVar(&listenerString, "l", ":9903", "listener address and port")
	flag.StringVar(&socketPath, "s", "/var/run/charon.vici", "path to vici socket")
}

func main() {
	flag.Parse()

	// Set up vici client
	sessionOption := vici.WithSocketPath(socketPath)
	session, err := vici.NewSession(sessionOption)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	done := make(chan bool)

	// Set up OpenMetrics endpoint
	http.Handle("/metrics", promhttp.HandlerFor(
		reg,
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		},
	))

	go PollVici(session, done)

	err = http.ListenAndServe(listenerString, nil)
	if err != nil {
		log.Fatal(err)
	}
}
