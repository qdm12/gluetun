# Private Internet Access Client

*Lightweight swiss-knife-like VPN client to tunnel to private internet access servers, using OpenVPN, iptables, DNS over TLS, ShadowSocks, Tinyproxy and more*

**ANNOUCEMENT**: *Total rewrite in Go: see the new features [below](#Features)* (in case something break use the image with tag `:old`)

<a href="https://hub.docker.com/r/qmcgaw/private-internet-access">
    <img width="100%" height="320" src="https://raw.githubusercontent.com/qdm12/private-internet-access-docker/master/title.svg?sanitize=true">
</a>

[![Build Status](https://travis-ci.org/qdm12/private-internet-access-docker.svg?branch=master)](https://travis-ci.org/qdm12/private-internet-access-docker)
[![Docker Pulls](https://img.shields.io/docker/pulls/qmcgaw/private-internet-access.svg)](https://hub.docker.com/r/qmcgaw/private-internet-access)
[![Docker Stars](https://img.shields.io/docker/stars/qmcgaw/private-internet-access.svg)](https://hub.docker.com/r/qmcgaw/private-internet-access)

[![GitHub last commit](https://img.shields.io/github/last-commit/qdm12/private-internet-access-docker.svg)](https://github.com/qdm12/private-internet-access-docker/issues)
[![GitHub commit activity](https://img.shields.io/github/commit-activity/y/qdm12/private-internet-access-docker.svg)](https://github.com/qdm12/private-internet-access-docker/issues)
[![GitHub issues](https://img.shields.io/github/issues/qdm12/private-internet-access-docker.svg)](https://github.com/qdm12/private-internet-access-docker/issues)

[![Image size](https://images.microbadger.com/badges/image/qmcgaw/private-internet-access.svg)](https://microbadger.com/images/qmcgaw/private-internet-access)
[![Image version](https://images.microbadger.com/badges/version/qmcgaw/private-internet-access.svg)](https://microbadger.com/images/qmcgaw/private-internet-access)
[![Join Slack channel](https://img.shields.io/badge/slack-@qdm12-yellow.svg?logo=slack)](https://join.slack.com/t/qdm12/shared_invite/enQtOTE0NjcxNTM1ODc5LTYyZmVlOTM3MGI4ZWU0YmJkMjUxNmQ4ODQ2OTAwYzMxMTlhY2Q1MWQyOWUyNjc2ODliNjFjMDUxNWNmNzk5MDk)

<details><summary>Click to show base components</summary><p>

- [Alpine 3.11](https://alpinelinux.org) for a tiny image (37MB of packages, 6.7MB of Go binary and 5.6MB for Alpine)
- [OpenVPN 2.4.8](https://pkgs.alpinelinux.org/package/v3.11/main/x86_64/openvpn) to tunnel to PIA servers
- [IPtables 1.8.3](https://pkgs.alpinelinux.org/package/v3.11/main/x86_64/iptables) enforces the container to communicate only through the VPN or with other containers in its virtual network (acts as a killswitch)
- [Unbound 1.9.6](https://pkgs.alpinelinux.org/package/v3.11/main/x86_64/unbound) configured with Cloudflare's [1.1.1.1](https://1.1.1.1) DNS over TLS (configurable with 5 different providers)
- [Files and blocking lists built periodically](https://github.com/qdm12/updated/tree/master/files) used with Unbound (see `BLOCK_MALICIOUS`, `BLOCK_SURVEILLANCE` and `BLOCK_ADS` environment variables)
- [TinyProxy 1.10.0](https://pkgs.alpinelinux.org/package/v3.11/main/x86_64/tinyproxy)
- [Shadowsocks 3.3.4](https://pkgs.alpinelinux.org/package/edge/testing/x86/shadowsocks-libev)

</p></details>

## Features

- **New features**
    - Choice to block ads, malicious and surveillance at the DNS level
    - All program output streams are merged (openvpn, unbound, shadowsocks, tinyproxy, etc.)
    - Choice of DNS over TLS provider(s)
    - Possibility of split horizon DNS by selecting multiple DNS over TLS providers
    - Download block lists and cryptographic files at start instead of at build time
    - Can work as a Kubernetes sidecar container, thanks @rorph
    - Pick a random region if no region is given, thanks @rorph
- <details><summary>Configure everything with environment variables</summary><p>

    - [Destination region](https://www.privateinternetaccess.com/pages/network)
    - Internet protocol
    - Level of encryption
    - PIA Username and password
    - DNS over TLS
    - DNS blocking: ads, malicious, surveillance
    - Internal firewall
    - Socks5 proxy
    - Web HTTP proxy

    </p></details>
- Connect
    - [Other containers to it](https://github.com/qdm12/private-internet-access-docker#connect-to-it)
    - [LAN devices to it](https://github.com/qdm12/private-internet-access-docker#connect-to-it)
- Killswitch using *iptables* to allow traffic only with needed PIA servers and LAN devices
- Port forwarding
- Compatible with amd64, i686 (32 bit), **ARM** 64 bit, ARM 32 bit v6 and v7, ppc64le and even that s390x 🎆
- Sub programs drop root privileges once launched: Openvpn, Unbound, Shadowsocks, Tinyproxy

## Setup

1. <details><summary>Requirements</summary><p>

    - A Private Internet Access **username** and **password** - [Sign up](https://www.privateinternetaccess.com/pages/buy-vpn/)
    - Docker API 1.25 to support `init`
    - If you use Docker Compose, docker-compose >= 1.22.0, to support `init: true`
    - <details><summary>External firewall requirements, if you have one</summary><p>

        - At start only
            - Allow outbound TCP 443 to github.com and privateinternetaccess.com
            - If `DOT=on`, allow outbound TCP 853 to 1.1.1.1 to allow Unbound to resolve the PIA domain name.
            - If `DOT=off`, allow outbound UDP 53 to your DNS provider to resolve the PIA domain name.
        - For UDP strong encryption, allow outbound UDP 1197 to the corresponding VPN server IPs
        - For UDP normal encryption, allow outbound UDP 1198 to the corresponding VPN server IPs
        - For TCP strong encryption, allow outbound TCP 501 to the corresponding VPN server IPs
        - For TCP normal encryption, allow outbound TCP 502 to the corresponding VPN server IPs
        - If `SHADOWSOCKS=on`, allow inbound TCP 8388 and UDP 8388 from your LAN
        - If `TINYPROXY=on`, allow inbound TCP 8888 from your LAN

    </p></details>

    </p></details>

1. Launch the container with:

    ```bash
    docker run -d --init --name=pia --cap-add=NET_ADMIN \
    -e REGION="CA Montreal" -e USER=js89ds7 -e PASSWORD=8fd9s239G \
    qmcgaw/private-internet-access
    ```

    or use [docker-compose.yml](https://github.com/qdm12/private-internet-access-docker/blob/master/docker-compose.yml) with:

    ```bash
    docker-compose up -d
    ```

    Note that you can:
    - Change the many [environment variables](#environment-variables) available
    - Use `-p 8888:8888/tcp` to access the HTTP web proxy (and put your LAN in `EXTRA_SUBNETS` environment variable)
    - Use `-p 8388:8388/tcp -p 8388:8388/udp` to access the SOCKS5 proxy (and put your LAN in `EXTRA_SUBNETS` environment variable)
    - Pass additional arguments to *openvpn* using Docker's command function (commands after the image name)
1. You can update the image with `docker pull qmcgaw/private-internet-access:latest`. There are also docker tags available:
    - `qmcgaw/private-internet-access:v1` linked to the [v1 release](https://github.com/qdm12/private-internet-access-docker/releases/tag/v1.0)

## Testing

Check the PIA IP address matches your expectations

```sh
docker run --rm --network=container:pia alpine:3.10 wget -qO- https://ipinfo.io
```

## Environment variables

| Environment variable | Default | Description |
| --- | --- | --- |
| `REGION` | `CA Montreal` | One of the [PIA regions](https://www.privateinternetaccess.com/pages/network/) |
| `PROTOCOL` | `udp` | `tcp` or `udp` |
| `ENCRYPTION` | `strong` | `normal` or `strong` |
| `USER` | | Your PIA username |
| `PASSWORD` | | Your PIA password |
| `DOT` | `on` | `on` or `off`, to activate DNS over TLS to 1.1.1.1 |
| `DOT_PROVIDERS` | `cloudflare` | Comma delimited list of DNS over TLS providers from `cloudflare`, `google`, `quad9`, `quadrant`, `cleanbrowsing`, `securedns`, `libredns` |
| `DOT_CACHING` | `on` | Unbound caching feature, `on` or `off` |
| `DOT_PRIVATE_ADDRESS` | All IPv4 and IPv6 CIDRs private ranges | Comma separated list of CIDRs or single IP addresses. Note that the default setting prevents DNS rebinding |
| `DOT_VERBOSITY` | `1` | Unbound verbosity level from `0` to `5` (full debug) |
| `DOT_VERBOSITY_DETAILS` | `0` | Unbound details verbosity level from `0` to `4` |
| `DOT_VALIDATION_LOGLEVEL` | `0` | Unbound validation log level from `0` to `2` |
| `BLOCK_MALICIOUS` | `on` | `on` or `off`, blocks malicious hostnames and IPs |
| `BLOCK_SURVEILLANCE` | `off` | `on` or `off`, blocks surveillance hostnames and IPs |
| `BLOCK_ADS` | `off` | `on` or `off`, blocks ads hostnames and IPs |
| `UNBLOCK` | | comma separated string (i.e. `web.com,web2.ca`) to unblock hostnames |
| `EXTRA_SUBNETS` | | comma separated subnets allowed in the container firewall (i.e. `192.168.1.0/24,192.168.10.121,10.0.0.5/28`) |
| `PORT_FORWARDING` | `off` | Set to `on` to forward a port on PIA server |
| `PORT_FORWARDING_STATUS_FILE` | `/forwarded_port` | File path to store the forwarded port number |
| `TINYPROXY` | `off` | `on` or `off`, to enable the internal HTTP proxy tinyproxy |
| `TINYPROXY_LOG` | `Info` | `Info`, `Connect`, `Notice`, `Warning`, `Error` or `Critical` |
| `TINYPROXY_PORT` | `8888` | `1024` to `65535` internal port for HTTP proxy |
| `TINYPROXY_USER` | | Username to use to connect to the HTTP proxy |
| `TINYPROXY_PASSWORD` | | Passsword to use to connect to the HTTP proxy |
| `SHADOWSOCKS` | `off` | `on` or `off`, to enable the internal SOCKS5 proxy Shadowsocks |
| `SHADOWSOCKS_LOG` | `on` | `on` or `off` to enable logging for Shadowsocks  |
| `SHADOWSOCKS_PORT` | `8388` | `1024` to `65535` internal port for SOCKS5 proxy |
| `SHADOWSOCKS_PASSWORD` | | Passsword to use to connect to the SOCKS5 proxy |
| `TZ` | | Specify a timezone to use i.e. `Europe/London` |

## Connect to it

There are various ways to achieve this, depending on your use case.

- <details><summary>Connect containers in the same docker-compose.yml as PIA</summary><p>

    Add `network_mode: "service:pia"` to your *docker-compose.yml* (no need for `depends_on`)

    </p></details>
- <details><summary>Connect other containers to PIA</summary><p>

    Add `--network=container:pia` when launching the container, provided PIA is already running

    </p></details>
- <details><summary>Connect containers from another docker-compose.yml</summary><p>

    Add `network_mode: "container:pia"` to your *docker-compose.yml*, provided PIA is already running

    </p></details>
- <details><summary>Connect LAN devices through the built-in HTTP proxy *Tinyproxy* (i.e. with Chrome, Kodi, etc.)</summary><p>

    You might want to use Shadowsocks instead which tunnels UDP as well as TCP, whereas Tinyproxy only tunnels TCP.

    1. Setup a HTTP proxy client, such as [SwitchyOmega for Chrome](https://chrome.google.com/webstore/detail/proxy-switchyomega/padekgcemlokbadohgkifijomclgjgif?hl=en)
    1. Ensure the PIA container is launched with:
        - port `8888` published `-p 8888:8888/tcp`
        - your LAN subnet, i.e. `192.168.1.0/24`, set as `-e EXTRA_SUBNETS=192.168.1.0/24`
    1. With your HTTP proxy client, connect to the Docker host (i.e. `192.168.1.10`) on port `8888`. You need to enter your credentials if you set them with `TINYPROXY_USER` and `TINYPROXY_PASSWORD`.
    1. If you set `TINYPROXY_LOG` to `Info`, more information will be logged in the Docker logs

    </p></details>
- <details><summary>Connect LAN devices through the built-in SOCKS5 proxy *Shadowsocks* (per app, system wide, etc.)</summary><p>

    1. Setup a SOCKS5 proxy client, there is a list of [ShadowSocks clients for **all platforms**](https://shadowsocks.org/en/download/clients.html)
        - **note** some clients do not tunnel UDP so your DNS queries will be done locally and not through PIA and its built in DNS over TLS
        - Clients that support such UDP tunneling are, as far as I know:
            - iOS: Potatso Lite
            - OSX: ShadowsocksX
            - Android: Shadowsocks by Max Lv
    1. Ensure the PIA container is launched with:
        - port `8388` published `-p 8388:8388/tcp -p 8388:8388/udp`
        - your LAN subnet, i.e. `192.168.1.0/24`, set as `-e EXTRA_SUBNETS=192.168.1.0/24`
    1. With your SOCKS5 proxy client
        - Enter the Docker host (i.e. `192.168.1.10`) as the server IP
        - Enter port TCP (and UDP, if available) `8388` as the server port
        - Use the password you have set with `SHADOWSOCKS_PASSWORD`
        - Choose the encryption method/algorithm `chacha20-ietf-poly1305`
    1. If you set `SHADOWSOCKS_LOG` to `on`, more information will be logged in the Docker logs

    </p></details>
- <details><summary>Access ports of containers connected to PIA</summary><p>

    In example, to access port `8000` of container `xyz`  and `9000` of container `abc` connected to PIA,
    publish ports `8000` and `9000` for the PIA container and access them as you would with any other container

    </p></details>
- <details><summary>Access ports of containers connected to PIA, all in the same docker-compose.yml</summary><p>

    In example, to access port `8000` of container `xyz`  and `9000` of container `abc` connected to PIA, publish port `8000` and `9000` for the PIA container.
    The docker-compose.yml file would look like:

    ```yml
    version: '3.7'
    services:
      pia:
        image: qmcgaw/private-internet-access
        container_name: pia
        init: true
        cap_add:
          - NET_ADMIN
        environment:
          - USER=js89ds7
          - PASSWORD=8fd9s239G
        ports:
          - 8000:8000/tcp
          - 9000:9000/tcp
      abc:
        image: abc
        container_name: abc
        network_mode: "service:pia"
      xyz:
        image: xyz
        container_name: xyz
        network_mode: "service:pia"
    ```

    </p></details>

## Port forwarding

By setting `PORT_FORWARDING` environment variable to `on`, the forwarded port will be read and written to the file specified in `PORT_FORWARDING_STATUS_FILE` (by default, this is set to `/forwarded_port`). If the location for this file does not exist, it will be created automatically.

You can mount this file as a volume to read it from other containers.

Note that not all regions support port forwarding.

## For the paranoids

- You can review the code which consists in:
    - [Dockerfile](https://github.com/qdm12/private-internet-access-docker/blob/master/Dockerfile)
    - [main.go](https://github.com/qdm12/private-internet-access-docker/blob/master/cmd/main.go), the main code entrypoint
    - [internal package](https://github.com/qdm12/private-internet-access-docker/blob/master/internal)
    - [github.com/qdm12/golibs](https://github.com/qdm12/golibs) dependency
    - [github.com/qdm12/files](https://github.com/qdm12/files) for files downloaded at start (DNS root hints, block lists, etc.)
- Build the image yourself:

    ```bash
    docker build -t qmcgaw/private-internet-access https://github.com/qdm12/private-internet-access-docker.git
    ```

- The download and parsing of all needed files is done at start (openvpn config files, Unbound files, block lists, etc.)
- Use `-e ENCRYPTION=strong -e BLOCK_MALICIOUS=on`
- You can test DNSSEC using [internet.nl/connection](https://www.internet.nl/connection/)
- Check DNS leak tests with [https://www.dnsleaktest.com](https://www.dnsleaktest.com)
- DNS Leaks tests might not work because of [this](https://github.com/qdm12/cloudflare-dns-server#verify-dns-connection) (*TLDR*: DNS server is a local caching intermediary)

## Troubleshooting

- If openvpn fails to start, you may need to:
    - Install the tun kernel module on your host with `insmod /lib/modules/tun.ko` or `modprobe tun`
    - Add `--device=/dev/net/tun` to your docker run command (equivalent for docker-compose, kubernetes, etc.)

- Fallback to a previous Docker image tags:
    - `v1` tag, stable shell scripting based (no support)
    - `old` tag, latest shell scripting version (no support)
    - `v2`... waiting for `latest` to become more stable

- Fallback to a precise previous version
    1. Clone the repository on your machine

        ```sh
        git clone https://github.com/qdm12/private-internet-access-docker.git pia
        cd pia
        ```

    1. Look up which commit you want to go back to [here](https://github.com/qdm12/private-internet-access-docker/commits/master), i.e. `942cc7d4d10545b6f5f89c907b7dd1dbc39368e0`
    1. Revert to this commit locally

        ```sh
        git reset --hard 942cc7d4d10545b6f5f89c907b7dd1dbc39368e0
        ```

    1. Build the Docker image

        ```sh
        docker build -t qmcgaw/private-internet-access .
        ```

## Development

### Using VSCode and Docker

1. Install [Docker](https://docs.docker.com/install)
    - On Windows, share a drive with Docker Desktop and have the project on that partition
1. With [Visual Studio Code](https://code.visualstudio.com/download), install the [remote containers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)
1. In Visual Studio Code, press on `F1` and select `Remote-Containers: Open Folder in Container...`
1. Your dev environment is ready to go!... and it's running in a container :+1:

## TODOs

- Support other VPN providers
    - Mullvad
    - Windscribe
- Gotify support for notificactions
- Periodic update of malicious block lists with Unbound restart
- Improve healthcheck
    - Check IP address belongs to selected region
    - Check for DNS provider somehow if this is even possible
- Support for other VPN protocols
    - Wireguard (wireguard-go)
- Show new versions/commits at start
- Colors & emojis
    - Setup
    - Logging streams
- More unit tests
- Write in Go
    - DNS over TLS to replace Unbound
    - HTTP proxy to replace tinyproxy
    - use [go-Shadowsocks2](https://github.com/shadowsocks/go-shadowsocks2)
    - DNS over HTTPS, maybe use [github.com/likexian/doh-go](https://github.com/likexian/doh-go)
    - use [iptables-go](https://github.com/coreos/go-iptables) to replace iptables
    - wireguard-go
    - Openvpn to replace openvpn

## License

This repository is under an [MIT license](https://github.com/qdm12/private-internet-access-docker/master/license)
