# plexheadend

[![GoDoc](https://godoc.org/github.com/wrboyce/plexheadend?status.svg)](https://godoc.org/github.com/wrboyce/plexheadend)
[![Go Report Card](https://goreportcard.com/badge/github.com/wrboyce/plexheadend)](https://goreportcard.com/report/github.com/wrboyce/plexheadend)
[![CircleCI](https://circleci.com/gh/wrboyce/plexheadend.png?style=shield)](https://circleci.com/gh/wrboyce/plexheadend)

Proxy requests between PlexDVR and TVHeadend. Now supports non-numeric channel numbers (i.e. "21.1").

## Installation

### Binary Release

Download the latest release from the [downloads page](https://github.com/wrboyce/plexheadend/releases).

### Build from Source

Download and build the project and its dependencies with the standard Go tooling, `go get github.com/wrboyce/plexheadend`.

### Docker Container

There is also a Docker container made available for use at `wrboyce/plexheadend`.

## Usage

All configuration options can be specified as either a commandline parameter or an environment variable.

| Name               | Commandline                 | Environment                   |
|--------------------|-----------------------------|-------------------------------|
| Device ID          | `--device-id` `-i`          | `PLEXHEADEND_DEVICE_ID`       |
| Name               | `--name` `-n`               | `PLEXHEADEND_NAME`            |
| Proxy Bind         | `--proxy-bind` `-b`         | `PLEXHEADEND_PROXY_BIND`      |
| Proxy Hostname     | `--proxy-hostname` `-H`     | `PLEXHEADEND_PROXY_HOSTNAME`  |
| Proxy Listen       | `--proxy-listen` `-l`       | `PLEXHEADEND_PROXY_LISTEN`    |
| Filter Tag         | `--tag` `-f`                | `PLEXHEADEND_TAG`             |
| Tuners             | `--tuners` `-t`             | `PLEXHEADEND_TUNERS`          |
| TVHeadend Host     | `--tvh-host` `-h`           | `PLEXHEADEND_TVH_HOST`        |
| TVHeadend Pass     | `--tvh-pass` `-P`           | `PLEXHEADEND_TVH_PASS`        |
| TVHeadend Port     | `--tvh-port` `-p`           | `PLEXHEADEND_TVH_PORT`        |
| TVHeadend User     | `--tvh-user` `-u`           | `PLEXHEADEND_TVH_USER`        |

```
Usage of plexheadend:
  -i, --device-id string        Device ID reported to Plex (default "1")
  -n, --name string             Friendly name reported to Plex (default "plexHeadend")
  -b, --proxy-bind string       Bind address (default all)
  -H, --proxy-hostname string   Hostname reported to Plex (default "localhost")
  -l, --proxy-listen string     Listen port (default "80")
  -f, --tag string              TVHeadend tag to filter reported channels (default none)
  -t, --tuners int              Number of Tuners reported to Plex (default 1)
  -h, --tvh-host string         TVHeadend Host (default "localhost")
  -P, --tvh-pass string         TVHeadend Password (default "plex")
  -p, --tvh-port string         TVHeadend Port (default "9981")
  -u, --tvh-user string         TVHeadend Username (default "plex")
```

### Example `docker-compose` Usage

```yaml
services:
    tvheadend:
        image: linuxserver/tvheadend
        container_name: tvheadend
        environment:
            - VERSION=latest
            - TZ=UTC
        ports:
            - 9981:9981
            - 9982:9982
        restart: unless-stopped
        
    plexheadend:
        image: wrboyce/plexheadend
        container_name: plexheadend
        environment:
            - PLEXHEADEND_TVH_HOST=tvheadend
            - PLEXHEADEND_PROXY_HOSTNAME=plexheadend
        restart: unless-stopped
        
    plex:
        image: linuxserver/plex
        container_name: plex
        environment:
            - VERSION=latest
            - TZ=UTC
        ports:
            - 32400:32400
            - 32400:32400/udp
        restart: unless-stopped
```
