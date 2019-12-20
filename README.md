# Private Internet Access Client

*Lightweight swiss-knife-like VPN client to tunnel to private internet access servers, using OpenVPN, iptables, DNS over TLS, ShadowSocks, Tinyproxy and more*

**ANNOUCEMENT**: I just published [*Kape acquisition of Private Internet Access: not worry you must*](https://link.medium.com/e70B1j0wz2)

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
[![Join Slack channel](https://img.shields.io/badge/slack-@qdm12-yellow.svg?logo=slack)](https://join.slack.com/t/qdm12/shared_invite/enQtODMwMDQyMTAxMjY1LTU1YjE1MTVhNTBmNTViNzJiZmQwZWRmMDhhZjEyNjVhZGM4YmIxOTMxOTYzN2U0N2U2YjQ2MDk3YmYxN2NiNTc)

<details><summary>Click to show base components</summary><p>

- [Alpine 3.10](https://alpinelinux.org) for a tiny image
- [OpenVPN 2.4.7](https://pkgs.alpinelinux.org/package/v3.10/main/x86_64/openvpn) to tunnel to PIA servers
- [IPtables 1.8.3](https://pkgs.alpinelinux.org/package/v3.10/main/x86_64/iptables) enforces the container to communicate only through the VPN or with other containers in its virtual network (acts as a killswitch)
- [Unbound 1.9.1](https://pkgs.alpinelinux.org/package/v3.10/main/x86_64/unbound) configured with Cloudflare's [1.1.1.1](https://1.1.1.1) DNS over TLS
- [Files and blocking lists built periodically](https://github.com/qdm12/updated/tree/master/files) used with Unbound (see `BLOCK_MALICIOUS` and `BLOCK_NSA` environment variables)
- [TinyProxy 1.10.0](https://pkgs.alpinelinux.org/package/v3.10/main/x86_64/tinyproxy)

</p></details>

## Features

- <details><summary>Configure everything with environment variables</summary><p>

    - [Destination region](https://www.privateinternetaccess.com/pages/network)
    - Internet protocol
    - Level of encryption
    - PIA Username and password
    - DNS over TLS
    - Malicious DNS blocking
    - Internal firewall
    - Web HTTP proxy
    - Run openvpn without root

    </p></details>
- Connect other containers to it, [see this](https://github.com/qdm12/private-internet-access-docker#connect-to-it)
- Compatible with amd64, i686 (32 bit), ARM 64 bit, ARM 32 bit v6 and v7, ppc64le and even that s390x 🎆
- Port forwarding
- The *iptables* firewall allows traffic only with needed PIA servers (IP addresses, port, protocol) combinations
- OpenVPN reconnects automatically on failure
- Docker healthcheck pings the DNS 1.1.1.1 to verify the connection is up
- Unbound DNS runs *without root*
- OpenVPN runs *without root* by default. You can run it with root with the environment variable `NONROOT=no`
- Connect your LAN devices
  - HTTP Web proxy *tinyproxy*
  - SOCKS5 proxy *shadowsocks* (better as it does UDP too)

## Setup

1. <details><summary>Requirements</summary><p>

    - A Private Internet Access **username** and **password** - [Sign up](https://www.privateinternetaccess.com/pages/buy-vpn/)
    - External firewall requirements, if you have one
        - Allow outbound TCP 853 to 1.1.1.1 to allow Unbound to resolve the PIA domain name at start. You can then block it once the container is started.
        - For UDP strong encryption, allow outbound UDP 1197
        - For UDP normal encryption, allow outbound UDP 1198
        - For TCP strong encryption, allow outbound TCP 501
        - For TCP normal encryption, allow outbound TCP 502
        - For the built-in web HTTP proxy, allow inbound TCP 8888
        - For the built-in SOCKS5 proxy, allow inbound TCP 8388 and UDP 8388
    - Docker API 1.25 to support `init`
    - If you use Docker Compose, docker-compose >= 1.22.0, to support `init: true`

    </p></details>

1. Ensure `/dev/net/tun` is setup on your host with either:

    ```sh
    insmod /lib/modules/tun.ko
    # or...
    modprobe tun
    ```

1. Launch the container with:

    ```bash
    docker run -d --init --name=pia --cap-add=NET_ADMIN --device=/dev/net/tun \
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
| `NONROOT` | `yes` | Run OpenVPN without root, `yes` or `no` |
| `DOT` | `on` | `on` or `off`, to activate DNS over TLS to 1.1.1.1 |
| `BLOCK_MALICIOUS` | `off` | `on` or `off`, blocks malicious hostnames and IPs |
| `BLOCK_NSA` | `off` | `on` or `off`, blocks NSA hostnames |
| `UNBLOCK` | | comma separated string (i.e. `web.com,web2.ca`) to unblock hostnames |
| `EXTRA_SUBNETS` | | comma separated subnets allowed in the container firewall (i.e. `192.168.1.0/24,192.168.10.121,10.0.0.5/28`) |
| `PORT_FORWARDING` | `off` | Set to `on` to forward a port on PIA server |
| `PORT_FORWARDING_STATUS_FILE` | `/forwarded_port` | File path to store the forwarded port number |
| `TINYPROXY` | `on` | `on` or `off`, to enable the internal HTTP proxy tinyproxy |
| `TINYPROXY_LOG` | `Critical` | `Info`, `Warning`, `Error` or `Critical` |
| `TINYPROXY_PORT` | `8888` | `1024` to `65535` internal port for HTTP proxy |
| `TINYPROXY_USER` | | Username to use to connect to the HTTP proxy |
| `TINYPROXY_PASSWORD` | | Passsword to use to connect to the HTTP proxy |
| `SHADOWSOCKS` | `on` | `on` or `off`, to enable the internal SOCKS5 proxy Shadowsocks |
| `SHADOWSOCKS_LOG` | `on` | `on` or `off` to enable logging for Shadowsocks  |
| `SHADOWSOCKS_PORT` | `8388` | `1024` to `65535` internal port for SOCKS5 proxy |
| `SHADOWSOCKS_PASSWORD` | | Passsword to use to connect to the SOCKS5 proxy |
| `TZ` | | Specify a timezone to use e.g. `Europe/London` |

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

    1. Setup a HTTP proxy client, such as [SwitchyOmega for Chrome](https://chrome.google.com/webstore/detail/proxy-switchyomega/padekgcemlokbadohgkifijomclgjgif?hl=en)
    1. Ensure the PIA container is launched with:
        - port `8888` published `-p 8888:8888/tcp`
        - your LAN subnet, i.e. `192.168.1.0/24`, set as `-e EXTRA_SUBNETS=192.168.1.0/24`
    1. With your HTTP proxy client, connect to the Docker host (i.e. `192.168.1.10`) on port `8888`. You need to enter your credentials if you set them with `TINYPROXY_USER` and `TINYPROXY_PASSWORD`.
    1. If you set `TINYPROXY_LOG` to `Info`, more information will be logged in the Docker logs, merged with the OpenVPN logs.
       `TINYPROXY_LOG` defaults to `Critical` to avoid logging everything, for privacy purposes.

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
    1. If you set `SHADOWSOCKS_LOG` to `on`, more information will be logged in the Docker logs, merged with the OpenVPN logs.

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
        devices:
          - /dev/net/tun
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

- You can review the code which essential consists in the [Dockerfile](https://github.com/qdm12/private-internet-access-docker/blob/master/Dockerfile) and [entrypoint.sh](https://github.com/qdm12/private-internet-access-docker/blob/master/entrypoint.sh)
- Build the images yourself:

    ```bash
    docker build -t qmcgaw/private-internet-access https://github.com/qdm12/private-internet-access-docker.git
    ```

- The download and unziping of PIA openvpn files is done at build for the ones not able to download the zip files
- Checksums for PIA openvpn zip files are not used as these files change often (but HTTPS is used)
- Use `-e ENCRYPTION=strong -e BLOCK_MALICIOUS=on`
- You can test DNSSEC using [internet.nl/connection](https://www.internet.nl/connection/)
- Check DNS leak tests with [https://www.dnsleaktest.com](https://www.dnsleaktest.com)
- DNS Leaks tests might not work because of [this](https://github.com/qdm12/cloudflare-dns-server#verify-dns-connection) (*TLDR*: DNS server is a local caching intermediary)

## Troubleshooting

- Password problems `AUTH: Received control message: AUTH_FAILED`
    - Your password may contain a special character such as `$`.
     You need to escape it with `\` in your run command or docker-compose.yml.
     For example you would set `-e PASSWORD=mypa\$\$word`.
- Fallback to a previous version
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

## TODOs

- Golang binary to setup the container at start, and:
  - Mix logs of unbound, tinyproxy, shadowsocks and openvpn together somehow
  - support other VPN providers
- Maybe use `--inactive 3600 --ping 10 --ping-exit 60` as default behavior
- Try without tun

## License

This repository is under an [MIT license](https://github.com/qdm12/private-internet-access-docker/master/license)
