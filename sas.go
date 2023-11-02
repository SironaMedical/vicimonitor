package main

import (
	"github.com/strongswan/govici/vici"
)

type ChildSA struct {
	Name         string   `vici:"name"`
	UniqueID     string   `vici:"uniqueid"`
	ReqID        int64    `vici:"reqid"`
	State        string   `vici:"state"`
	Mode         string   `vici:"mode"`
	Protocol     string   `vici:"protocol"`
	Encap        string   `vici:"encap"`
	SPIIn        string   `vici:"spi-in"`
	SPIOut       string   `vici:"spi-out"`
	CPIIn        string   `vici:"cpi-in"`
	CPIOut       string   `vici:"cpi-out"`
	MarkIn       string   `vici:"mark-in"`
	MarkMaskIn   string   `vici:"mark-mask-in"`
	MarkOut      string   `vici:"mark-out"`
	MarkMaskOut  string   `vici:"mark-mask-out"`
	IfIDIn       string   `vici:"if-id-in"`
	IfIDOut      string   `vici:"if-id-out"`
	Label        string   `vici:"label"`
	EncrAlg      string   `vici:"encr-alg"`
	EncrKeySize  int64    `vici:"encr-keysize"`
	IntegAlg     string   `vici:"integ-alg"`
	IntegKeySize int64    `vici:"integ-keysize"`
	PrfAlg       string   `vici:"prf-alg"`
	DHGroup      string   `vici:"dh-group"`
	ESN          int64    `vici:"esn"`
	BytesIn      int64    `vici:"bytes-in"`
	PacketsIn    int64    `vici:"packets-in"`
	UseIn        int64    `vici:"use-in"`
	BytesOut     int64    `vici:"bytes-out"`
	PacketsOut   int64    `vici:"packets-out"`
	UseOut       int64    `vici:"use-out"`
	RekeyTime    int64    `vici:"rekey-time"`
	Lifetime     int64    `vici:"life-time"`
	InstallTime  int64    `vici:"install-time"`
	LocalTS      []string `vici:"local-ts"`
	RemoteTS     []string `vici:"remote-ts"`
}

type IkeSA struct {
	Name          string
	UniqueID      string             `vici:"uniqueid"`
	Version       int64              `vici:"version"`
	State         string             `vici:"state"`
	LocalHost     string             `vici:"local-host"`
	LocalPort     int64              `vici:"local-port"`
	LocalID       string             `vici:"local-id"`
	RemoteHost    string             `vici:"remote-host"`
	RemotePort    int64              `vici:"remote-port"`
	RemoteID      string             `vici:"remote-id"`
	RemoteXAuthID string             `vici:"remote-xauth-id"`
	RemoteEapID   string             `vici:"remote-eap-id"`
	Initiator     string             `vici:"initiator"`
	InitiatorSPI  string             `vici:"initiator-spi"`
	ResponderSPI  string             `vici:"responder-spi"`
	NatLocal      string             `vici:"nat-local"`
	NatRemote     string             `vici:"nat-remote"`
	NatFake       string             `vici:"nat-fake"`
	NatAny        string             `vici:"nat-any"`
	IfIDIn        string             `vici:"if-id-in"`
	IfIDOut       string             `vici:"if-id-out"`
	EncrAlg       string             `vici:"encr-alg"`
	EncrKeySize   int64              `vici:"encr-keysize"`
	IntegAlg      string             `vici:"integ-alg"`
	IntegKeySize  int64              `vici:"integ-keysize"`
	PrfAlg        string             `vici:"prf-alg"`
	DHGroup       string             `vici:"dh-group"`
	Established   int64              `vici:"established"`
	RekeyTime     int64              `vici:"rekey-time"`
	ReauthTime    int64              `vici:"reauth-time"`
	LocalVIPs     []string           `vici:"local-vips"`
	RemoteVIPs    []string           `vici:"remote-vips"`
	TasksQueued   []string           `vici:"tasks-queued"`
	TasksActive   []string           `vici:"tasks-active"`
	TasksPassive  []string           `vici:"tasks-passive"`
	ChildSAs      map[string]ChildSA `vici:"child-sas"`
}

func GetSAs(session *vici.Session) (sas map[string]*IkeSA, err error) {
	sas = map[string]*IkeSA{}
	m := vici.NewMessage()
	ms, err := session.StreamedCommandRequest("list-sas", "list-sa", m)
	if err != nil {
		return
	}

	for _, msg := range ms {
		if m.Err() != nil {
			return
		}
		for _, name := range msg.Keys() {
			if name != "" {
				var sa IkeSA
				m := msg.Get(name).(*vici.Message)
				if err = vici.UnmarshalMessage(m, &sa); err != nil {
					return
				}
				sa.Name = name
				sas[name] = &sa
			}
		}
	}
	return
}
