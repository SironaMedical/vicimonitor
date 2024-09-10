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
}

func NewManager(session *vici.Session, interval time.Duration) *Manager {
	return &Manager{
		session: session,
		ticker:  time.NewTicker(interval),
		metrics: *metrics.NewCollector(session),
		monitor: *monitor.NewMonitor(session),
	}
}

func (m *Manager) Run() {
	for {
		select {
		case <-m.ticker.C:
			if err := m.metrics.Update(); err != nil {
				log.Println(err)
			}
		case message := <-m.metrics.C:
			if err := m.monitor.InitiateSA(message); err != nil {
				log.Println(err)
			}
		}
	}
}

func (m *Manager) Shutdown() error {
	m.ticker.Stop()
	if err := m.session.Close(); err != nil {
		return err
	}
	return nil
}
