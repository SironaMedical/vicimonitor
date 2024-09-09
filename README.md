# vicimonitor

## Introduction
`vicimonitor` is a tool designed to be run as a daemon that interacts with the [Versatile IKE Control Interface (VICI) protocol](https://github.com/strongswan/strongswan/blob/master/src/libcharon/plugins/vici/README.md) of [Strongswan](https://www.strongswan.org) to monitor IPSec VPN connections.
`vicimonitor` exposes Security Association metrics as an OpenMetrics endpoint.

## Usage
```$ vicimonitor --help
Usage of vicimonitor:
  -interval int
    	The interval to update metrics in seconds (default 30)
  -listen string
    	The listen address (default "0.0.0.0:9000")
  -socket string
    	The vici socket path (default "/var/run/charon.vici")
  -version
    	Display the version and exit
```

## OpenMetrics endpoint
By default the OpenMetrics endpoint is exposed on `0.0.0.0:9000/metrics`

## Development
`vicimonitor` requires Go version 1.21.
After cloning, run `go mod tidy` to install dependencies.
Tests can be run with `go test`.

### Releasing
Github releases will automatically build and upload artifacts via this [github action](.github/workflows/release.yaml)
