# vicimonitor

## Introduction
`vicimonitor` is a tool designed to be run as a daemon that interacts with the [Versatile IKE Control Interface (VICI) protocol](https://www.strongswan.org/apidoc/md_src_libcharon_plugins_vici_README.html) of [Strongswan](https://www.strongswan.org) to monitor IPSec VPN connections.
`vicimonitor` exposes Security Association metrics as an OpenMetrics endpoint as well attempts to initiate Security Associations missing from the list of connections.

## Usage
```$ vicimonitor --help
Usage of vicimonitor:
  -i int
    	poll interval duration in seconds for the vici socket (default 5)
  -l string
    	listener address and port for the OpenMetrics endpoint (default ":9903")
  -s string
    	path to vici socket (default "/var/run/charon.vici")
  -t int
    	SA initiate timeout (default 3000)
  -version
      print the version of vicimonitor and exit
```

## OpenMetrics endpoint
By default the OpenMetrics endpoint is exposed on `0.0.0.0:9903/metrics`
