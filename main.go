package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/strongswan/govici/vici"

	"sironamedical/vicimonitor/pkg/metrics"
)

func main() {
	listenAddr := flag.String("listen", "0.0.0.0:9000", "The listen address")
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

	session, err := vici.NewSession(vici.WithSocketPath(*socketPath))
	if err != nil {
		log.Fatalln("error connecting to vici socket ", err)
	}
	defer session.Close()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	httpServer := &http.Server{Addr: *listenAddr, Handler: promhttp.Handler()}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("http server error ", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		collectMetrics(ctx, session, time.Duration(*tickerInterval))
	}()

	log.Println("vicimonitor started...")
	<-sigChan

	log.Println("http server shutting down")
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Println("http server shutdown error ", err)
	}

	cancel()
	wg.Wait()
}

func collectMetrics(ctx context.Context, session *vici.Session, interval time.Duration) {
	ticker := time.NewTicker(interval * time.Second)
	defer ticker.Stop()

	monitor := metrics.NewPrometheus(session)
	for {
		select {
		case <-ctx.Done():
			log.Println("stopping metric collection")
			return
		case <-ticker.C:
			monitor.Update()
		}
	}
}
