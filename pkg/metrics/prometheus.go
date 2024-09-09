package metrics

import (
	"fmt"
	"log"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/strongswan/govici/vici"

	"sironamedical/vicimonitor/pkg/vici/messages"
)

func NewPrometheus(session *vici.Session) *Prometheus {
	return &Prometheus{session: session}
}

type Prometheus struct {
	session *vici.Session
	C       chan string
}

func (p *Prometheus) Update() {
	p.updateSecurityAssociation()
}

func (p *Prometheus) forceInitiateIke(sa string) {
	ikeForceRestart.WithLabelValues(sa).Add(1)
	log.Println(fmt.Sprintf("force restart for ike sa %v",sa))
	p.C <- sa
}

func (p *Prometheus) updateSecurityAssociation() {
	sas, err := p.session.StreamedCommandRequest("list-sas", "list-sa", nil)
	if err != nil {
		log.Println("Unable to list SAS:", err)
		return
	}

	for _, mesg := range sas {
		if err := mesg.Err(); err != nil {
			log.Println("Error in message:", err)
			continue
		}

		for _, key := range mesg.Keys() {
			inner := mesg.Get(key).(*vici.Message)
			var listSAS messages.ListSAS
			if err := vici.UnmarshalMessage(inner, &listSAS); err != nil {
				log.Println("Error unmarshalling message:", err)
				continue
			}

			currentState := IkeSAStateMap[listSAS.State]
			if currentState > 2 {
				p.forceInitiateIke(key)
				return
			}

			ikeState.WithLabelValues(key).Set(currentState)
			ikeRekeyTime.WithLabelValues(key).Set(float64(listSAS.ReKeyTime))

			for _, child := range listSAS.Children {
				lables := prometheus.Labels{
					"name":        child.Name,
					"parent_name": key,
					"local_ts":    strings.Join(child.LocalTS, ","),
					"remote_ts":   strings.Join(child.RemoteTS, ","),
				}
				childBytesIn.With(lables).Set(float64(child.BytesIn))
				childBytesOut.With(lables).Set(float64(child.BytesOut))

				currentChildState := ChildSAStateMap[child.State]
				if currentChildState > 3 {
					p.forceInitiateIke(key)
					return
				}
				childState.With(lables).Set(currentChildState)
			}
		}
	}
}
