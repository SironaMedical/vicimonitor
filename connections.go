package main

import (
	"github.com/strongswan/govici/vici"
)

type Child struct {
	LocalTS      []string `vici:"local-ts"`
	RemoteTS     []string `vici:"remote-ts"`
	Mode         string   `vici:"mode"`
	Label        []string `vici:"label"`
	ReauthTime   int64    `vici:"reauth_time"`
	RekeyTime    int64    `vici:"rekey_time"`
	RekeyBytes   int64    `vici:"rekey_bytes"`
	RekeyPackets int64    `vici:"rekey_packets"`
}

type Auth struct {
	Class string `vici:"class"`
	ID    string `vici:"id"`
}

type Conn struct {
	Name         string
	LocalAddrs   []string         `vici:"local_addrs"`
	RemoteAddrs  []string         `vici:"remote_addrs"`
	ReauthTime   int64            `vici:"reauth_time"`
	RekeyTime    int64            `vici:"rekey_time"`
	Local        Auth             `vici:"local-1"`
	Remote       Auth             `vici:"remote-1"`
	Children     map[string]Child `vici:"children"`
	IKEVersion   string           `vici:"version"`
	IKEProposals []string         `vici:"proposals"`
}

func GetConns(session *vici.Session) (conns map[string]*Conn, err error) {
	conns = map[string]*Conn{}
	m := vici.NewMessage()
	ms, err := session.StreamedCommandRequest("list-conns", "list-conn", m)
	if err != nil {
		return
	}

	for _, msg := range ms {
		if msg.Err() != nil {
			return
		}
		for _, name := range msg.Keys() {
			if name != "" {
				var c Conn
				m := msg.Get(name).(*vici.Message)
				if err = vici.UnmarshalMessage(m, &c); err != nil {
					return
				}
				c.Name = name
				conns[name] = &c
			}
		}
	}
	return
}
