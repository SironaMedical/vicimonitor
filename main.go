package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sironamedical/vicimonitor/pkg/metrics"
	"sironamedical/vicimonitor/pkg/monitor"
	"sync"
	"syscall"
	"time"

	"github.com/strongswan/govici/vici"
)

func main() {
	listenAddr := flag.String("listen", "0.0.0.0:9000", "The listen address")
	reinitiate := flag.Bool("reinitiate", false, "Attempt to initiate SAs")
	socketPath := flag.String("socket", "/var/run/charon.vici", "The vici socket path")
	tickerInterval := flag.Int("interval", 30, "The interval to update metrics in seconds")
	version := flag.Bool("version", false, "Display the version and exit")
	flag.Parse()

	if *version {
		fmt.Println(Version)
		return
	}

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	session, err := vici.NewSession(vici.WithSocketPath(*socketPath))
	if err != nil {
		log.Fatalln("error connecting to vici socket ", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	mux := http.NewServeMux()
	mux.Handle("/metrics", metrics.Handler)
	httpServer := &http.Server{Addr: *listenAddr, Handler: mux}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("http server error ", err)
		}
	}()

	monitor := monitor.NewMonitor(session, time.Duration(*tickerInterval)*time.Second, *reinitiate)

	wg.Add(1)
	go func() {
		defer wg.Done()
		monitor.Run()
	}()

	log.Println("vicimonitor started...")
	<-sigChan

	log.Println("vicimonitor shuting down...")
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Println("http server shutdown error ", err)
	}

	if err := monitor.Shutdown(); err != nil {
		log.Println("reactor shutdown error ", err)
	}

	wg.Wait()
}
