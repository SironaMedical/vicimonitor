package manager

import (
	"log"
	"time"

	"sironamedical/vicimonitor/pkg/metrics"
	"sironamedical/vicimonitor/pkg/monitor"

	"github.com/strongswan/govici/vici"
)

type Manager struct {
	session *vici.Session
	ticker  *time.Ticker
	metrics metrics.Collector
	monitor monitor.Monitor

	initiateChan chan *vici.Message
	shutdownChan chan struct{}
}

func NewManager(session *vici.Session, interval time.Duration) *Manager {
	return &Manager{
		session: session,
		ticker:  time.NewTicker(interval),
		metrics: *metrics.NewCollector(session),
		monitor: *monitor.NewMonitor(session),

		initiateChan: make(chan *vici.Message, 1),
		shutdownChan: make(chan struct{}, 1),
	}
}

func (m *Manager) Run() {
	for {
		select {
		case <-m.shutdownChan:
			m.ticker.Stop()
			if err := m.session.Close(); err != nil {
				log.Println("unable to close vici session ", err)
			}
			return
		case <-m.ticker.C:
			if err := m.metrics.Update(); err != nil {
				log.Println(err)
			}
		case message := <-m.metrics.C:
			m.initiateChan <- message
		case sa := <-m.initiateChan:
			if err := m.monitor.InitiateSA(sa); err != nil {
				log.Println(err)
			}
		}
	}
}

func (m *Manager) Shutdown() error {
	close(m.shutdownChan)
	m.ticker.Stop()
	if err := m.session.Close(); err != nil {
		return err
	}
	return nil
}
