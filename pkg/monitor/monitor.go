package monitor

import (
	"fmt"
	"log"
	"sironamedical/vicimonitor/pkg/vici/messages"

	"github.com/strongswan/govici/vici"
)

func NewMonitor(session *vici.Session) *Monitor {
	return &Monitor{
		session: session,
		C:       make(chan *vici.Message, 1),
	}
}

type Monitor struct {
	session *vici.Session
	C       chan *vici.Message
}

func (m *Monitor) InitiateSA(message *vici.Message) error {
	ike := message.Get("ike")
	child := message.Get("child")

	if child == nil {
		log.Printf("initiating ike SA %v\n", ike)
	} else {
		log.Printf("initiating child sa %v from %v", child, ike)
	}

	mesgs, err := m.session.StreamedCommandRequest("initiate", "control-log", message)
	if err != nil {
		return err
	}
	for _, msg := range mesgs {
		if err := msg.Err(); err != nil {
			return err
		}

		var cLog messages.ControlLog
		if err = vici.UnmarshalMessage(msg, &cLog); err != nil {
			return err
		}

		log.Println(fmt.Printf("%v %v %v", cLog.Level, cLog.Group, cLog.Message))
		if cLog.Message == "" {
			break
		}
	}
	return nil
}
