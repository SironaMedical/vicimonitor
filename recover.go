package main

import (
	"github.com/strongswan/govici/vici"
)

func InitiateSA(name string, session *vici.Session) (err error) {
	m := vici.NewMessage()
	if err := m.Set("child", name); err != nil {
		return err
	}
	if err := m.Set("timeout", initiateTimeout); err != nil {
		return err
	}
	if err := m.Set("loglevel", 1); err != nil {
		return err
	}
	messages, err := session.StreamedCommandRequest("initiate", "control-log", m)
	if err != nil {
		return err
	}
	for _, msg := range messages {
		if err = msg.Err(); err != nil {
			return err
		}
		var c ControlLog
		if err = vici.UnmarshalMessage(msg, &c); err != nil {
			return err
		}
		if c.Message == "" {
			break
		}
		EmitControlLog(c)
	}
	return nil
}
