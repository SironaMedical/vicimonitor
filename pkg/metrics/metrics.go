package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
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

	ikeRekeyTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "ipsec",
			Subsystem: "ike_sa",
			Name:      "rekey_seconds",
			Help:      "Time until IKE rekey event.",
		},
		[]string{"name"},
	)
	ikeState = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "ipsec",
			Subsystem: "ike_sa",
			Name:      "state",
			Help:      "IKE SA state code",
		},
		[]string{"name"},
	)
	childState = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "ipsec",
			Subsystem: "child_sa",
			Name:      "state",
			Help:      "Child SA state code",
		},
		[]string{"name", "local_ts", "remote_ts", "parent_name"},
	)
	childBytesIn = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "ipsec",
			Subsystem: "child_sa",
			Name:      "in_bytes",
			Help:      "Child SA Bytes In",
		},
		[]string{"name", "local_ts", "remote_ts", "parent_name"},
	)
	childBytesOut = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "ipsec",
			Subsystem: "child_sa",
			Name:      "out_bytes",
			Help:      "Child SA Bytes Out",
		},
		[]string{"name", "local_ts", "remote_ts", "parent_name"},
	)
)

func init() {
	prometheus.MustRegister(
		childBytesIn,
		childBytesOut,
		childState,
		ikeRekeyTime,
		ikeState,
	)
}
