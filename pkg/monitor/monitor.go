package monitor

import (
	"log"
	"strings"
	"time"

	"sironamedical/vicimonitor/pkg/metrics"
	"sironamedical/vicimonitor/pkg/vici/messages"

	"github.com/strongswan/govici/vici"
)

type SAMapKey struct {
	name string
	id   int64
}

func NewMonitor(session *vici.Session, interval time.Duration, initiate bool) *Monitor {
	return &Monitor{
		session:      session,
		ticker:       time.NewTicker(interval),
		initiate:     initiate,
		shutdownChan: make(chan struct{}),
	}
}

type Monitor struct {
	ticker       *time.Ticker
	session      *vici.Session
	initiate     bool
	shutdownChan chan struct{}
}

func (m *Monitor) Run() {
	for {
		select {
		case <-m.ticker.C:
			if err := m.monitor(); err != nil {
				// exit on failure
				// sessions are not handled
				// simpler to let systemd restart for now
				log.Fatalln(err)
			}
		case <-m.shutdownChan:
			return
		}
	}
}

func (m *Monitor) Shutdown() error {
	m.ticker.Stop()
	m.shutdownChan <- struct{}{}
	if err := m.session.Close(); err != nil {
		return err
	}
	return nil
}

// Each connection has an expected number of Child SAs based on the number of local and remote traffic selectors.
// We compare the number of Child SAs to the expected number and if they do not match we initiate the IKE SAs.
func (m *Monitor) monitor() error {
	// reset these metrics
	// when strongswan deletes an SA we have no way to know
	metrics.ChildState.Reset()
	metrics.IkeState.Reset()

	connections, err := m.getConnections()
	if err != nil {
		return err
	}

	connTunnelCountMap := make(map[string]map[string]int)
	for connName, conn := range connections {
		for tunnelName, tunnel := range conn.Children {
			tunnelCount := len(tunnel.LocalTS) * len(tunnel.RemoteTS)
			connTunnelCountMap[connName] = make(map[string]int)
			connTunnelCountMap[connName][tunnelName] += tunnelCount
		}
	}

	securityAssociations, err := m.getSecurityAssociations()
	if err != nil {
		return err
	}

	// We only care about the following states
	// based on ChildSAStateMap values
	accetableStates := []int{0, 1, 2, 3, 4, 5, 6, 7}
	accetableCountMap := make(map[string]int)
	for mk, sa := range securityAssociations {
		metrics.IkeState.WithLabelValues(mk.name).Set(metrics.IkeSAStateMap[sa.State])
		metrics.IkeRekeyTime.WithLabelValues(mk.name).Set(float64(sa.ReKeyTime))
		for _, child := range sa.Children {
			promLabels := []string{child.Name, strings.Join(child.LocalTS, ","), strings.Join(child.RemoteTS, ","), mk.name}
			childState := metrics.ChildSAStateMap[child.State]

			metrics.ChildBytesIn.WithLabelValues(promLabels...).Set(float64(child.BytesIn))
			metrics.ChildBytesOut.WithLabelValues(promLabels...).Set(float64(child.BytesOut))
			metrics.ChildState.WithLabelValues(promLabels...).Set(childState)

			for _, v := range accetableStates {
				if v == int(childState) {
					accetableCountMap[child.Name]++
				}
			}
		}
	}

	if m.initiate {
		for conn, tunnels := range connTunnelCountMap {
			for tunnel, count := range tunnels {
				if count != accetableCountMap[tunnel] {
					log.Printf("Connection %s, Tunnel %s, Found %d, Expected %d\n", conn, tunnel, accetableCountMap[tunnel], count)
					if err := m.initiateIkeSAs(conn); err != nil {
						return err
					}
				}
			}
		}
	}

	allKeys := []SAMapKey{}
	for k := range securityAssociations {
		allKeys = append(allKeys, k)
	}

	metrics.OverLappingSAs.Reset()
	for i := 0; i < len(allKeys); i++ {
		for j := i + 1; j < len(allKeys); j++ {
			if allKeys[i].name == allKeys[j].name {
				metrics.OverLappingSAs.WithLabelValues(allKeys[i].name).Set(1)
			}
		}
	}

	return nil
}

func (m *Monitor) getConnections() (map[string]messages.ListConn, error) {
	conns, err := m.session.StreamedCommandRequest("list-conns", "list-conn", nil)
	if err != nil {
		return nil, err
	}

	connections := make(map[string]messages.ListConn)
	for _, conn := range conns {
		if err := conn.Err(); err != nil {
			return nil, err
		}
		for _, key := range conn.Keys() {
			inner := conn.Get(key).(*vici.Message)

			var listConn messages.ListConn
			if err := vici.UnmarshalMessage(inner, &listConn); err != nil {
				return nil, err
			}
			listConn.LocalAuth = make(map[string]messages.ListConnAuthSection)
			listConn.RemoteAuth = make(map[string]messages.ListConnAuthSection)
			for _, k := range inner.Keys() {
				if strings.HasPrefix(k, "local-") {
					var auth messages.ListConnAuthSection
					if err := vici.UnmarshalMessage(inner.Get(k).(*vici.Message), &auth); err != nil {
						return nil, err
					}
					newKey, _ := strings.CutPrefix(k, "local-")
					listConn.LocalAuth[newKey] = auth
				}
				if strings.HasPrefix(k, "remote-") {
					var auth messages.ListConnAuthSection
					if err := vici.UnmarshalMessage(inner.Get(k).(*vici.Message), &auth); err != nil {
						return nil, err
					}
					newKey, _ := strings.CutPrefix(k, "remote-")
					listConn.RemoteAuth[newKey] = auth
				}
			}
			connections[key] = listConn
		}
	}
	return connections, nil
}

func (m *Monitor) getSecurityAssociations() (map[SAMapKey]messages.ListSAS, error) {
	sas, err := m.session.StreamedCommandRequest("list-sas", "list-sa", nil)
	if err != nil {
		return nil, err
	}

	securityAssociations := make(map[SAMapKey]messages.ListSAS)
	for _, sa := range sas {
		if err := sa.Err(); err != nil {
			return nil, err
		}
		for _, key := range sa.Keys() {
			inner := sa.Get(key).(*vici.Message)

			var listSAS messages.ListSAS
			if err := vici.UnmarshalMessage(inner, &listSAS); err != nil {
				return nil, err
			}
			sk := SAMapKey{name: key, id: listSAS.UniqueID}
			securityAssociations[sk] = listSAS
		}
	}

	return securityAssociations, nil
}

func (m *Monitor) initiateIkeSAs(ike string) error {
	initiate := &messages.Initiate{
		Ike:        ike,
		Timeout:    3600,
		InitLimits: true,
		LogLevel:   2,
	}
	initiateArgs, err := vici.MarshalMessage(initiate)
	if err != nil {
		return err
	}
	log.Printf("Initiating Ike SA: %v\n", ike)
	resp, err := m.session.CommandRequest("initiate", initiateArgs)
	if err != nil {
		return err
	}
	if err := resp.Err(); err != nil {
		return err
	}
	return nil
}
