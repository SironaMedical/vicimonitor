package metrics

import (
	"fmt"
	"log"
	"strings"

	"sironamedical/vicimonitor/pkg/vici/messages"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/strongswan/govici/vici"
)

func NewCollector(session *vici.Session) *Collector {
	return &Collector{
		session: session,
		C:       make(chan *vici.Message, 1),
	}
}

type Collector struct {
	session *vici.Session
	C       chan *vici.Message
}

func (c *Collector) Update() error {
	sas, err := c.session.StreamedCommandRequest("list-sas", "list-sa", nil)
	if err != nil {
		return fmt.Errorf("Unable to run list-sas %v", err)
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
			ikeState.WithLabelValues(key).Set(currentState)
			ikeRekeyTime.WithLabelValues(key).Set(float64(listSAS.ReKeyTime))

			initateMessage := vici.NewMessage()
			if err := initateMessage.Set("ike", key); err != nil {
				return err
			}

			if currentState > 2 {
				c.C <- initateMessage
				ikeForceRestart.WithLabelValues(key).Add(1)
				return nil
			}

			for _, child := range listSAS.Children {
				lables := prometheus.Labels{
					"name":        child.Name,
					"parent_name": key,
					"local_ts":    strings.Join(child.LocalTS, ","),
					"remote_ts":   strings.Join(child.RemoteTS, ","),
				}
				currentChildState := ChildSAStateMap[child.State]
				childState.With(lables).Set(currentChildState)
				childBytesIn.With(lables).Set(float64(child.BytesIn))
				childBytesOut.With(lables).Set(float64(child.BytesOut))
				if currentChildState > 3 {
					if err := initateMessage.Set("child", child.Name); err != nil {
						return err
					}
					c.C <- initateMessage
					ikeForceRestart.WithLabelValues(key).Add(1)
					return nil
				}
			}
		}
	}
	return nil
}
