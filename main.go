package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/strongswan/govici/vici"
)

var Version = "development"

// Set flags
var pollInterval int
var initiateTimeout int
var listenerString string
var socketPath string
var versionFlag bool

func init() {
	flag.IntVar(&pollInterval, "i", 5, "poll interval duration in seconds for the vici socket")
	flag.IntVar(&initiateTimeout, "t", 3000, "SA initiate timeout")
	flag.StringVar(&listenerString, "l", ":9903", "listener address and port")
	flag.StringVar(&socketPath, "s", "/var/run/charon.vici", "path to vici socket")
	flag.BoolVar(&versionFlag, "version", false, "print the version of vicimonitor and exit")
}

func main() {

	flag.Parse()
	if versionFlag {
		fmt.Println(Version)
		os.Exit(0)
	}

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

	slog.Info("starting vici poller")
	go PollVici(session, done)

	Infof("starting OpenMetrics http listener on: %s", listenerString)
	err = http.ListenAndServe(listenerString, nil)
	if err != nil {
		log.Fatal(err)
	}
}
