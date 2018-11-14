# Private Internet Access Client (OpenVPN+Iptables+DNS over TLS on Alpine Linux)

*Lightweight VPN client to tunnel to private internet access servers*

**WARNING: auth.conf is now replaced by the environment variables `USER` and `PASSWORD`, please update your configuration**

[![PIA Docker OpenVPN](https://github.com/qdm12/private-internet-access-docker/raw/master/readme/title.png)](https://hub.docker.com/r/qmcgaw/private-internet-access/)

[![Build Status](https://travis-ci.org/qdm12/private-internet-access-docker.svg?branch=master)](https://travis-ci.org/qdm12/private-internet-access-docker)
[![Docker Build Status](https://img.shields.io/docker/build/qmcgaw/private-internet-access.svg)](https://hub.docker.com/r/qmcgaw/private-internet-access)

[![GitHub last commit](https://img.shields.io/github/last-commit/qdm12/private-internet-access-docker.svg)](https://github.com/qdm12/private-internet-access-docker/issues)
[![GitHub commit activity](https://img.shields.io/github/commit-activity/y/qdm12/private-internet-access-docker.svg)](https://github.com/qdm12/private-internet-access-docker/issues)
[![GitHub issues](https://img.shields.io/github/issues/qdm12/private-internet-access-docker.svg)](https://github.com/qdm12/private-internet-access-docker/issues)

[![Docker Pulls](https://img.shields.io/docker/pulls/qmcgaw/private-internet-access.svg)](https://hub.docker.com/r/qmcgaw/private-internet-access)
[![Docker Stars](https://img.shields.io/docker/stars/qmcgaw/private-internet-access.svg)](https://hub.docker.com/r/qmcgaw/private-internet-access)
[![Docker Automated](https://img.shields.io/docker/automated/qmcgaw/private-internet-access.svg)](https://hub.docker.com/r/qmcgaw/private-internet-access)

[![Image size](https://images.microbadger.com/badges/image/qmcgaw/private-internet-access.svg)](https://microbadger.com/images/qmcgaw/private-internet-access)
[![Image version](https://images.microbadger.com/badges/version/qmcgaw/private-internet-access.svg)](https://microbadger.com/images/qmcgaw/private-internet-access)

| Image size | RAM usage | CPU usage |
| --- | --- | --- |
| 20MB | 14MB to 80MB | Low to Medium |

It is based on:

- [Alpine 3.8](https://alpinelinux.org) for a tiny image
- [OpenVPN 2.4.6-r3](https://pkgs.alpinelinux.org/package/v3.8/main/x86_64/openvpn) to tunnel to PIA servers
- [IPtables 1.6.2-r0](https://pkgs.alpinelinux.org/package/v3.8/main/x86_64/iptables) enforces the container to communicate only through the VPN or with other containers in its virtual network (acts as a killswitch)
- [Unbound 1.7.3-r0](https://pkgs.alpinelinux.org/package/v3.8/main/x86_64/unbound) configured with Cloudflare's [1.1.1.1](https://1.1.1.1) DNS over TLS
- [Malicious hostnames list](https://github.com/qdm12/malicious-hostnames-docker) used with Unbound (see `BLOCK_MALICIOUS` environment variable)
- [Malicious IPs list](https://github.com/qdm12/malicious-ips-docker) used with Unbound (see `BLOCK_MALICIOUS`)

## Extra features

- Only use environment variables:
    - the [destination region]((https://www.privateinternetaccess.com/pages/network/))
    - the protocol `tcp` or `udp`
    - the level of encryption `normal` or `strong`
- Connect other containers to it
- The *iptables* firewall allows traffic only with needed PIA servers (IP addresses, port, protocol) combination
- OpenVPN restarts on failure using another PIA IP address for the same region
- Docker healthchecks using [duckduckgo.com](https://duckduckgo.com) to obtain your public IP address and compare it with your initial non-VPN IP address
- Openvpn and Unbound do not run as root

## Requirements

- A Private Internet Access **username** and **password** - [Sign up](https://www.privateinternetaccess.com/pages/buy-vpn/)
- [Docker](https://docs.docker.com/install/) installed on the host
- If you use a firewall on the host:
  - Allow outgoing TCP port 853 for Cloudflare DNS over TLS initial resolution of PIA server domain name, **you should then BLOCK it**
  - Allow outgoing TCP port 443 for querying duckduckgo.com to obtain the initial IP address *only at the start of the container*, **you should then BLOCK it**
  - Allow outgoing TCP port 501 for TCP strong encryption
  - Allow outgoing TCP port 502 for TCP normal encryption
  - Allow outgoing UDP port 1197 for UDP strong encryption
  - Allow outgoing UDP port 1198 for UDP normal encryption

## Setup

1. Make sure you have your `/dev/net/tun` device setup on your host with one of the following commands, depending on your OS:

    ```bash
    insmod /lib/modules/tun.ko
    ```

    Or

    ```bash
    modprobe tun
    ```

1. Launch the container with:

    ```bash
    docker run -d --name=pia -v ./auth.conf:/auth.conf:ro \
    --cap-add=NET_ADMIN --device=/dev/net/tun --network=pianet \
    -e REGION="CA Montreal" -e PROTOCOL=udp -e ENCRYPTION=strong \
    -e USER=js89ds7 -e PASSWORD=8fd9s239G \
    qmcgaw/private-internet-access
    ```

    or use [docker-compose.yml](https://github.com/qdm12/private-internet-access-docker/blob/master/docker-compose.yml) with:

    ```bash
    docker-compose up -d
    ```

    Note that you can change all the [environment variables](#environment-variables)
1. Wait about 5 seconds for it to connect to the PIA server. You can check with:

    ```bash
    docker logs -f pia
    ```

1. Follow the [**Testing section**](#testing)

## Testing

You can simply use the Docker healthcheck. The container will mark itself as **unhealthy** if the public IP address is the same as your initial public IP address. Otherwise you can follow these instructions:

1. Check your host IP address with:

    ```bash
    wget -qO- https://ipinfo.io/ip
    ```

1. Run the same command in a Docker container using your *pia* container as network with:

    ```bash
    docker run --rm --network=container:pia alpine:3.8 wget -qO- https://ipinfo.io/ip
    ```

    If the displayed IP address appears and is different that your host IP address, the PIA client works !

## Environment variables

| Environment variable | Default | Description |
| --- | --- | --- |
| `REGION` | `CA Montreal` | One of the [PIA regions](https://www.privateinternetaccess.com/pages/network/) |
| `PROTOCOL` | `udp` | `tcp` or `udp` |
| `ENCRYPTION` | `strong` | `normal` or `strong` |
| `BLOCK_MALICIOUS` | `off` | `on` or `off` |
| `USER` | `` | Your PIA username |
| `PASSWORD` | `` | Your PIA password |
| `EXTRA_SUBNETS` | `` | Comma separated subnets allowed in the container firewall |

`EXTRA_SUBNETS` can be in example: `192.168.1.0/24,192.168.10.121,10.0.0.5/28`

## Connect other containers to it

Connect other Docker containers to the PIA VPN connection by adding `--network=container:pia` when launching them.

## For the paranoids

- You can review the code which essential consits in the [Dockerfile](https://github.com/qdm12/private-internet-access-docker/blob/master/Dockerfile) and [entrypoint.sh](https://github.com/qdm12/private-internet-access-docker/blob/master/entrypoint.sh)
- Build the images yourself:

    ```bash
    docker build -t qmcgaw/private-internet-access https://github.com/qdm12/private-internet-access-docker.git
    ```

- The download and unziping of PIA openvpn files is done at build for the ones not able to download the zip files
- Checksums for PIA openvpn zip files are not used as these files change often (but HTTPS is used)
- Use `-e ENCRYPTION=strong -e BLOCK_MALICIOUS=on`

## TODOs

- [ ] Malicious IPs and hostnames with wget at launch+checksums
- [ ] Su Exec (fork and addition)
- [ ] SOCKS proxy/Hiproxy/VPN server for other devices to use the container

## License

This repository is under an [MIT license](https://github.com/qdm12/private-internet-access-docker/master/license)
