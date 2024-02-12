package main

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	reg           = prometheus.NewRegistry()
	factory       = promauto.With(reg)
	IkeSAStateMap = map[string]float64{
		"CREATED":     0,
		"CONNECTING":  1,
		"ESTABLISHED": 2,
		"PASSIVE":     3,
		"REKEYING":    4,
		"REKEYED":     5,
		"DELETING":    6,
		"DESTROYING":  7,
	}
	ChildSAStateMap = map[string]float64{
		"CREATED":    0,
		"ROUTED":     1,
		"INSTALLING": 2,
		"INSTALLED":  3,
		"UPDATING":   4,
		"REKEYING":   5,
		"REKEYED":    6,
		"RETRYING":   7,
		"DELETING":   8,
		"DELETED":    9,
		"DESTROYING": 10,
	}

	ikeRekeyTime = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "ipsec",
			Subsystem: "ike_sa",
			Name:      "rekey_seconds",
			Help:      "Time until IKE rekey event.",
		},
		[]string{"name"},
	)
	ikeState = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "ipsec",
			Subsystem: "ike_sa",
			Name:      "state",
			Help:      "IKE SA state code",
		},
		[]string{"name"},
	)
	childState = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "ipsec",
			Subsystem: "child_sa",
			Name:      "state",
			Help:      "Child SA state code",
		},
		[]string{"name", "local_ts", "remote_ts", "parent_name"},
	)
	childBytesIn = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "ipsec",
			Subsystem: "child_sa",
			Name:      "in_bytes",
			Help:      "Child SA Bytes In",
		},
		[]string{"name", "local_ts", "remote_ts", "parent_name"},
	)
	childBytesOut = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "ipsec",
			Subsystem: "child_sa",
			Name:      "out_bytes",
			Help:      "Child SA Bytes Out",
		},
		[]string{"name", "local_ts", "remote_ts", "parent_name"},
	)
)

func UpdateSAMetrics(sa *IkeSA) {
	ikeRekeyTime.With(prometheus.Labels{"name": sa.Name}).Set(float64(sa.RekeyTime))
	ikeState.With(prometheus.Labels{"name": sa.Name}).Set(IkeSAStateMap[sa.State])
	for _, csa := range sa.ChildSAs {
		labels := prometheus.Labels{
			"name":        csa.Name,
			"local_ts":    strings.Join(csa.LocalTS, ","),
			"remote_ts":   strings.Join(csa.RemoteTS, ","),
			"parent_name": sa.Name,
		}
		childState.With(labels).Set(ChildSAStateMap[csa.State])
		childBytesIn.With(labels).Set(float64(csa.BytesIn))
		childBytesOut.With(labels).Set(float64(csa.BytesOut))
	}
}
