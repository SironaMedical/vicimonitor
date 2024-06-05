package main

import "os/exec"

func StartIpsec() (err error) {
	err = exec.Command("/usr/sbin/ipsec", "start").Run()
	return
}
