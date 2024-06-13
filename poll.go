package main

import (
	"log"
	"log/slog"
	"strings"
	"time"

	"github.com/mitchellh/go-ps"
	"github.com/prometheus/client_golang/prometheus"
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
			procs, err := ps.Processes()
			if err != nil {
				Errorf("%s", err)
				break
			}
			if !IsIpsecProcressRunning(procs) {
				Errorf("IPSec process not running, attempting to start.")
				err := StartIpsec()
				if err != nil {
					Errorf("%s", err)
					break
				}
			}
			conns, err := GetConns(session)
			if err != nil {
				log.Fatal(err)
			}
			sas, err := GetSAs(session)
			if err != nil {
				log.Fatal(err)
			}
			for _, sa := range sas {
				UpdateSAMetrics(sa)
			}
			for _, c := range conns {
				for childName, child := range c.Children {
					if !ChildSAExists(childName, sas[c.Name]) {
						// Mark Missing CHILD_SA as deleted
						labels := prometheus.Labels{
							"name":        childName,
							"local_ts":    strings.Join(child.LocalTS, ","),
							"remote_ts":   strings.Join(child.RemoteTS, ","),
							"parent_name": c.Name,
						}
						childState.With(labels).Set((ChildSAStateMap["DELETED"]))
						Errorf("CHILD_SA %s does not exist, attempting initialization.", childName)
						err := InitiateSA(childName, session)
						if err != nil {
							Errorf("%s", err)
						}
					}
				}
			}
		}
	}
}
