package main

import (
	"strings"

	"github.com/mitchellh/go-ps"
)

func ChildSAExists(s string, sa *IkeSA) (exists bool) {
	exists = false
	if sa != nil {
		for _, childsa := range sa.ChildSAs {
			if strings.EqualFold(childsa.Name, s) {
				exists = true
			}
		}
	}
	return
}

func IsIkeSAEstablished(sa IkeSA) (connected bool) {
	return strings.EqualFold(sa.State, "ESTABLISHED")
}

func IsIpsecProcressRunning(procs []ps.Process) (running bool) {
	running = false
	for _, p := range procs {
		if strings.Contains(p.Executable(), "charon") {
			running = true
			return
		}
	}
	return
}
