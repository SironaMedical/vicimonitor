package main

import (
	"log/slog"
	"time"

	"github.com/strongswan/govici/vici"
)

func PollVici(session *vici.Session, channel chan bool) (conns map[string]*Conn, sas map[string]*IkeSA, err error) {
	ticker := time.NewTicker(time.Duration(pollInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-channel:
			slog.Info("stopping vici poller")
			return
		case <-ticker.C:
			conns, err := GetConns(session)
			if err != nil {
				Errorf("%s", err)
			}
			sas, err := GetSAs(session)
			if err != nil {
				Errorf("%s", err)
			}
			for _, sa := range sas {
				UpdateSAMetrics(sa)
			}
			for _, c := range conns {
				for childName := range c.Children {
					if !ChildSAExists(childName, sas[c.Name]) {
						Errorf("CHILD_SA %s does not exist, attempting initialization.", childName)
						err := InitiateSA(childName, session)
						if err != nil {
							Errorf("%s", err)
						}
					} else {
						slog.Info("all CHILD_SAs exist")
					}
				}
			}
		}
	}
}
