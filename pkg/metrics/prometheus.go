package metrics

import (
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
}

func (p *Prometheus) Update() {
	p.updateSecurityAssociation()
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

			ikeState.WithLabelValues(key).Set(IkeSAStateMap[listSAS.State])
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
				childState.With(lables).Set(ChildSAStateMap[child.State])
			}
		}
	}
}
