package main

import "strings"

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
