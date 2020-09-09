# Gluetun VPN client

*Lightweight swiss-knife-like VPN client to tunnel to Private Internet Access,
Mullvad, Windscribe, Surfshark Cyberghost, VyprVPN, NordVPN and PureVPN VPN servers, using Go, OpenVPN,
iptables, DNS over TLS, ShadowSocks and Tinyproxy*

**ANNOUNCEMENT**: *Youtube videos added*

<img height="250" src="https://raw.githubusercontent.com/qdm12/gluetun/master/title.svg?sanitize=true">

[![Build status](https://github.com/qdm12/gluetun/workflows/Buildx%20latest/badge.svg)](https://github.com/qdm12/gluetun/actions?query=workflow%3A%22Buildx+latest%22)
[![Docker Pulls](https://img.shields.io/docker/pulls/qmcgaw/private-internet-access.svg)](https://hub.docker.com/r/qmcgaw/private-internet-access)
[![Docker Stars](https://img.shields.io/docker/stars/qmcgaw/private-internet-access.svg)](https://hub.docker.com/r/qmcgaw/private-internet-access)

[![GitHub last commit](https://img.shields.io/github/last-commit/qdm12/gluetun.svg)](https://github.com/qdm12/gluetun/issues)
[![GitHub commit activity](https://img.shields.io/github/commit-activity/y/qdm12/gluetun.svg)](https://github.com/qdm12/gluetun/issues)
[![GitHub issues](https://img.shields.io/github/issues/qdm12/gluetun.svg)](https://github.com/qdm12/gluetun/issues)

[![Image size](https://images.microbadger.com/badges/image/qmcgaw/private-internet-access.svg)](https://microbadger.com/images/qmcgaw/private-internet-access)
[![Image version](https://images.microbadger.com/badges/version/qmcgaw/private-internet-access.svg)](https://microbadger.com/images/qmcgaw/private-internet-access)
[![Join Slack channel](https://img.shields.io/badge/slack-@qdm12-yellow.svg?logo=slack)](https://join.slack.com/t/qdm12/shared_invite/enQtOTE0NjcxNTM1ODc5LTYyZmVlOTM3MGI4ZWU0YmJkMjUxNmQ4ODQ2OTAwYzMxMTlhY2Q1MWQyOWUyNjc2ODliNjFjMDUxNWNmNzk5MDk)

## Videos

1. [**Introduction**](https://youtu.be/3jIbU6J2Hs0)
1. [**Connect a container**](https://youtu.be/mH7J_2JKNK0)
1. [**Connect LAN devices**](https://youtu.be/qvjrM15Y0uk)

## Features

- Based on Alpine 3.12 for a small Docker image of 52MB
- Supports **Private Internet Access** (new and old), **Mullvad**, **Windscribe**, **Surfshark**, **Cyberghost**, **Vyprvpn**, **NordVPN** and **PureVPN** servers
- Supports Openvpn only for now
- DNS over TLS baked in with service provider(s) of your choice
- DNS fine blocking of malicious/ads/surveillance hostnames and IP addresses, with live update every 24 hours
- Choose the vpn network protocol, `udp` or `tcp`
- Built in firewall kill switch to allow traffic only with needed the VPN servers and LAN devices
- Built in SOCKS5 proxy (Shadowsocks, tunnels TCP+UDP)
- Built in HTTP proxy (Tinyproxy, tunnels TCP)
- [Connect other containers to it](https://github.com/qdm12/gluetun#connect-to-it)
- [Connect LAN devices to it](https://github.com/qdm12/gluetun#connect-to-it)
- Compatible with amd64, i686 (32 bit), **ARM** 64 bit, ARM 32 bit v6 and v7 🎆
- VPN server side port forwarding for Private Internet Access and Vyprvpn
- Possibility of split horizon DNS by selecting multiple DNS over TLS providers
- Subprograms all drop root privileges once launched
- Subprograms output streams are all merged together
- Can work as a Kubernetes sidecar container, thanks @rorph

## Setup

1. Requirements
    - A VPN account with one of the service providers supported
    - If you have a host or router firewall, please refer [to the firewall documentation](https://github.com/qdm12/gluetun/wiki/External-firewall-requirements)
1. On some devices you may need to setup your tunnel kernel module on your host with `insmod /lib/modules/tun.ko` or `modprobe tun`
    - *Synology users*: please read [this part of the Wiki](https://github.com/qdm12/gluetun/wiki/Common-issues#synology)
1. Launch the container with:

    ```bash
    docker run -d --name gluetun --cap-add=NET_ADMIN \
    -e REGION="CA Montreal" -e USER=js89ds7 -e PASSWORD=8fd9s239G \
    -v /yourpath:/gluetun \
    qmcgaw/private-internet-access
    ```

    or use [docker-compose.yml](https://github.com/qdm12/gluetun/blob/master/docker-compose.yml) with:

    ```bash
    docker-compose up -d
    ```

    Note that you can:

    - Change the many [environment variables](#environment-variables) available
    - Use `-p 8888:8888/tcp` to access the HTTP web proxy (and put your LAN in `EXTRA_SUBNETS` environment variable, in example `192.168.1.0/24`)
    - Use `-p 8388:8388/tcp -p 8388:8388/udp` to access the SOCKS5 proxy (and put your LAN in `EXTRA_SUBNETS` environment variable, in example `192.168.1.0/24`)
    - Use `-p 8000:8000/tcp` to access the [HTTP control server](#HTTP-control-server) built-in

    **If you encounter an issue with the tun device not being available, see [the FAQ](https://github.com/qdm12/gluetun/blob/master/doc/faq.md#how-to-fix-openvpn-failing-to-start)**

1. You can update the image with `docker pull qmcgaw/private-internet-access:latest`. See the [wiki](https://github.com/qdm12/gluetun/wiki/Common-issues#use-a-release-tag) for more information on other tags available.

## Testing

Check the VPN IP address matches your expectations

```sh
docker run --rm --network=container:gluetun alpine:3.12 wget -qO- https://ipinfo.io
```

Want more testing? ▶ [see the Wiki](https://github.com/qdm12/gluetun/wiki/Testing)

## Environment variables

**TLDR**; only set the 🏁 marked environment variables to get started.

### VPN

| Variable | Default | Choices | Description |
| --- | --- | --- | --- |
| 🏁 `VPNSP` | `private internet access` | `private internet access`, `private internet access old`, `mullvad`, `windscribe`, `surfshark`, `vyprvpn`, `nordvpn`, `purevpn` | VPN Service Provider |
| `IP_STATUS_FILE` | `/tmp/gluetun/ip` | Any filepath | Filepath to store the public IP address assigned |
| `PROTOCOL` | `udp` | `udp` or `tcp` | Network protocol to use |
| `OPENVPN_VERBOSITY` | `1` | `0` to `6` | Openvpn verbosity level |
| `OPENVPN_ROOT` | `no` | `yes` or `no` | Run OpenVPN as root |
| `OPENVPN_TARGET_IP` | | Valid IP address | Specify a target VPN server (or gateway) IP address to use |
| `OPENVPN_CIPHER` | | i.e. `aes-256-gcm` | Specify a custom cipher to use. It will also set `ncp-disable` if using AES GCM for PIA |
| `OPENVPN_AUTH` | | i.e. `sha256` | Specify a custom auth algorithm to use |

*For all providers below, server location parameters are all optional. By default a random server is picked using the filter settings provided.*

- Private Internet Access

    | Variable | Default | Choices | Description |
    | --- | --- | --- | --- |
    | 🏁 `USER` | | | Your username |
    | 🏁 `PASSWORD` | | | Your password |
    | `REGION` | | One of the [PIA regions](https://www.privateinternetaccess.com/pages/network/) | VPN server region |
    | `PIA_ENCRYPTION` | `strong` | `normal`, `strong` | Encryption preset |
    | `PORT_FORWARDING` | `off` | `on`, `off` | Enable port forwarding on the VPN server **for old only** |
    | `PORT_FORWARDING_STATUS_FILE` | `/tmp/gluetun/forwarded_port` | Any filepath | Filepath to store the forwarded port number **for old only** |

- Mullvad

    | Variable | Default | Choices | Description |
    | --- | --- | --- | --- |
    | 🏁 `USER` | | | Your user ID |
    | `COUNTRY` | | One of the [Mullvad countries](https://mullvad.net/en/servers/#openvpn) | VPN server country |
    | `CITY` | | One of the [Mullvad cities](https://mullvad.net/en/servers/#openvpn) | VPN server city |
    | `ISP` | | One of the [Mullvad ISP](https://mullvad.net/en/servers/#openvpn) | VPN server ISP |
    | `PORT` | | `80`, `443` or `1401` for TCP; `53`, `1194`, `1195`, `1196`, `1197`, `1300`, `1301`, `1302`, `1303` or `1400` for UDP. Defaults to TCP `443` and UDP `1194` | Custom VPN port to use |

- Windscribe

    | Variable | Default | Choices | Description |
    | --- | --- | --- | --- |
    | 🏁 `USER` | | | Your username |
    | 🏁 `PASSWORD` | | | Your password |
    | `REGION` | | One of the [Windscribe regions](https://windscribe.com/status) | VPN server region |
    | `PORT` | | One from the [this list of ports](https://windscribe.com/getconfig/openvpn) | Custom VPN port to use |

- Surfshark

    | Variable | Default | Choices | Description |
    | --- | --- | --- | --- |
    | 🏁 `USER` | | | Your **service** username, found at the bottom of the [manual setup page](https://account.surfshark.com/setup/manual) |
    | 🏁 `PASSWORD` | | | Your **service** password |
    | `REGION` | | One of the [Surfshark regions](https://github.com/qdm12/gluetun/wiki/surfshark) | VPN server region |

- Cyberghost

    | Variable | Default | Choices | Description |
    | --- | --- | --- | --- |
    | 🏁 `USER` | | | Your username |
    | 🏁 `PASSWORD` | | | Your password |
    | 🏁 `CLIENT_KEY` | | | Your device client key content, **see below** |
    | `REGION` | | One of the [Cyberghost countries](https://github.com/qdm12/gluetun/wiki/Cyberghost#regions) | VPN server country |
    | `CYBERGHOST_GROUP` | `Premium UDP Europe` | One of the [server groups](https://github.com/qdm12/gluetun/wiki/Cyberghost#server-groups) | Server group |

    To specify your client key, you can either:

    - Bind mount it at `/files/client.key`, for example with `-v /yourpath/client.key:/files/client.key:ro`
    - Convert it to a single line value using:

        ```sh
        docker run -it --rm -v /yourpath/client.key:/files/client.key:ro qmcgaw/private-internet-access clientkey
        ```

        And use the line produced as the value for the environment variable `CLIENT_KEY`.

- Vyprvpn

    | Variable | Default | Choices | Description |
    | --- | --- | --- | --- |
    | 🏁 `USER` | | | Your username |
    | 🏁 `PASSWORD` | | | Your password |
    | `REGION` | | One of the [VyprVPN regions](https://www.vyprvpn.com/server-locations) | VPN server region |

- NordVPN

    | Variable | Default | Choices | Description |
    | --- | --- | --- | --- |
    | 🏁 `USER` | | | Your username |
    | 🏁 `PASSWORD` | | | Your password |
    | `REGION` | | One of the NordVPN server country, i.e. `Switzerland` | VPN server country |
    | `SERVER_NUMBER` | | Server integer number | Optional server number. For example `251` for `Italy #251` |

- PureVPN

    | Variable | Default | Choices | Description |
    | --- | --- | --- | --- |
    | 🏁 `USER` | | | Your user ID |
    | 🏁 `REGION` | | One of the [PureVPN regions](https://support.purevpn.com/vpn-servers) | VPN server region |
    | `COUNTRY` | | One of the [PureVPN countries](https://support.purevpn.com/vpn-servers) | VPN server country |
    | `CITY` | | One of the [PureVPN cities](https://support.purevpn.com/vpn-servers) | VPN server city |

### DNS over TLS

None of the following values are required.

| Variable | Default | Choices | Description |
| --- | --- | --- | --- |
| `DOT` | `on` | `on`, `off` | Activate DNS over TLS with Unbound |
| `DOT_PROVIDERS` | `cloudflare` | `cloudflare`, `google`, `quad9`, `quadrant`, `cleanbrowsing`, `securedns`, `libredns` | Comma delimited list of DNS over TLS providers |
| `DOT_CACHING` | `on` | `on`, `off` | Unbound caching |
| `DOT_IPV6` | `off` | `on`, `off` | DNS IPv6 resolution |
| `DOT_PRIVATE_ADDRESS` | All private CIDRs ranges | | Comma separated list of CIDRs or single IP addresses Unbound won't resolve to. Note that the default setting prevents DNS rebinding |
| `DOT_VERBOSITY` | `1` | `0` to `5` | Unbound verbosity level |
| `DOT_VERBOSITY_DETAILS` | `0` | `0` to `4` | Unbound details verbosity level |
| `DOT_VALIDATION_LOGLEVEL` | `0` | `0` to `2` | Unbound validation log level |
| `DNS_UPDATE_PERIOD` | `24h` | i.e. `0`, `30s`, `5m`, `24h` | Period to update block lists and cryptographic files and restart Unbound. Set to `0` to deactivate updates |
| `BLOCK_MALICIOUS` | `on` | `on`, `off` | Block malicious hostnames and IPs with Unbound |
| `BLOCK_SURVEILLANCE` | `off` | `on`, `off` | Block surveillance hostnames and IPs with Unbound |
| `BLOCK_ADS` | `off` | `on`, `off` | Block ads hostnames and IPs with Unbound |
| `UNBLOCK` | |i.e. `domain1.com,x.domain2.co.uk` | Comma separated list of domain names to leave unblocked with Unbound |
| `DNS_PLAINTEXT_ADDRESS` | `1.1.1.1` | Any IP address | IP address to use as DNS resolver if `DOT` is `off` |
| `DNS_KEEP_NAMESERVER` | `off` | `on` or `off` | Keep the nameservers in /etc/resolv.conf untouched, but disabled DNS blocking features |

### Firewall

That one is important if you want to connect to the container from your LAN for example, using Shadowsocks or Tinyproxy.

| Variable | Default | Choices | Description |
| --- | --- | --- | --- |
| `FIREWALL` | `on` | `on` or `off` | Turn on or off the container built-in firewall. You should use it for **debugging purposes** only. |
| `EXTRA_SUBNETS` | | i.e. `192.168.1.0/24,192.168.10.121,10.0.0.5/28` | Comma separated subnets allowed in the container firewall |
| `FIREWALL_VPN_INPUT_PORTS` | | i.e. `1000,8080` | Comma separated list of ports to allow from the VPN server side (useful for **vyprvpn** port forwarding) |
| `FIREWALL_DEBUG` | `off` | `on` or `off` | Prints every firewall related command. You should use it for **debugging purposes** only. |

### Shadowsocks

| Variable | Default | Choices | Description |
| --- | --- | --- | --- |
| `SHADOWSOCKS` | `off` | `on`, `off` | Enable the internal SOCKS5 proxy Shadowsocks |
| `SHADOWSOCKS_LOG` | `off` | `on`, `off` | Enable logging |
| `SHADOWSOCKS_PORT` | `8388` | `1024` to `65535` | Internal port number for Shadowsocks to listen on |
| `SHADOWSOCKS_PASSWORD` | |  | Password to use to connect to Shadowsocks |
| `SHADOWSOCKS_METHOD` | `chacha20-ietf-poly1305` | `chacha20-ietf-poly1305`, `aes-128-gcm`, `aes-256-gcm` | Method to use for Shadowsocks |

### Tinyproxy

| Variable | Default | Choices | Description |
| --- | --- | --- | --- |
| `TINYPROXY` | `off` | `on`, `off` | Enable the internal HTTP proxy tinyproxy |
| `TINYPROXY_LOG` | `Info` | `Info`, `Connect`, `Notice`, `Warning`, `Error`, `Critical` | Tinyproxy log level |
| `TINYPROXY_PORT` | `8888` | `1024` to `65535` | Internal port number for Tinyproxy to listen on |
| `TINYPROXY_USER` | | | Username to use to connect to Tinyproxy |
| `TINYPROXY_PASSWORD` | | | Password to use to connect to Tinyproxy |

### System

| Variable | Default | Choices | Description |
| --- | --- | --- | --- |
| `TZ` | | i.e. `Europe/London` | Specify a timezone to use to have correct log times |
| `UID` | `1000` | | User ID to run as non root and for ownership of files written |
| `GID` | `1000` | | Group ID to run as non root and for ownership of files written |

### Other

| Variable | Default | Choices | Description |
| --- | --- | --- | --- |
| `PUBLICIP_PERIOD` | `12h` | Valid duration | Period to check for public IP address. Set to `0` to disable. |
| `VERSION_INFORMATION` | `on` | `on`, `off` | Logs a message indicating if a newer version is available once the VPN is connected |

## Connect to it

There are various ways to achieve this, depending on your use case.

- <details><summary>Connect containers in the same docker-compose.yml as Gluetun</summary><p>

    Add `network_mode: "service:gluetun"` to your *docker-compose.yml* (no need for `depends_on`)

    </p></details>
- <details><summary>Connect other containers to Gluetun</summary><p>

    Add `--network=container:gluetun` when launching the container, provided Gluetun is already running

    </p></details>
- <details><summary>Connect containers from another docker-compose.yml</summary><p>

    Add `network_mode: "container:gluetun"` to your *docker-compose.yml*, provided Gluetun is already running

    </p></details>
- <details><summary>Connect LAN devices through the built-in HTTP proxy *Tinyproxy* (i.e. with Chrome, Kodi, etc.)</summary><p>

    You might want to use Shadowsocks instead which tunnels UDP as well as TCP, whereas Tinyproxy only tunnels TCP.

    1. Setup a HTTP proxy client, such as [SwitchyOmega for Chrome](https://chrome.google.com/webstore/detail/proxy-switchyomega/padekgcemlokbadohgkifijomclgjgif?hl=en)
    1. Ensure the Gluetun container is launched with:
        - port `8888` published `-p 8888:8888/tcp`
        - your LAN subnet, i.e. `192.168.1.0/24`, set as `-e EXTRA_SUBNETS=192.168.1.0/24`
    1. With your HTTP proxy client, connect to the Docker host (i.e. `192.168.1.10`) on port `8888`. You need to enter your credentials if you set them with `TINYPROXY_USER` and `TINYPROXY_PASSWORD`.
    1. If you set `TINYPROXY_LOG` to `Info`, more information will be logged in the Docker logs

    </p></details>
- <details><summary>Connect LAN devices through the built-in SOCKS5 proxy *Shadowsocks* (per app, system wide, etc.)</summary><p>

    1. Setup a SOCKS5 proxy client, there is a list of [ShadowSocks clients for **all platforms**](https://shadowsocks.org/en/download/clients.html)
        - **note** some clients do not tunnel UDP so your DNS queries will be done locally and not through Gluetun and its built in DNS over TLS
        - Clients that support such UDP tunneling are, as far as I know:
            - iOS: Potatso Lite
            - OSX: ShadowsocksX
            - Android: Shadowsocks by Max Lv
    1. Ensure the Gluetun container is launched with:
        - port `8388` published `-p 8388:8388/tcp -p 8388:8388/udp`
        - your LAN subnet, i.e. `192.168.1.0/24`, set as `-e EXTRA_SUBNETS=192.168.1.0/24`
    1. With your SOCKS5 proxy client
        - Enter the Docker host (i.e. `192.168.1.10`) as the server IP
        - Enter port TCP (and UDP, if available) `8388` as the server port
        - Use the password you have set with `SHADOWSOCKS_PASSWORD`
        - Choose the encryption method/algorithm to the method you specified in `SHADOWSOCKS_METHOD`
    1. If you set `SHADOWSOCKS_LOG` to `on`, (a lot) more information will be logged in the Docker logs

    </p></details>
- <details><summary>Access ports of containers connected to Gluetun</summary><p>

    In example, to access port `8000` of container `xyz`  and `9000` of container `abc` connected to Gluetun,
    publish ports `8000` and `9000` for the Gluetun container and access them as you would with any other container

    </p></details>
- <details><summary>Access ports of containers connected to Gluetun, all in the same docker-compose.yml</summary><p>

    In example, to access port `8000` of container `xyz`  and `9000` of container `abc` connected to Gluetun, publish port `8000` and `9000` for the Gluetun container.
    The docker-compose.yml file would look like:

    ```yml
    version: '3.7'
    services:
      gluetun:
        image: qmcgaw/private-internet-access
        container_name: gluetun
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
        network_mode: "service:gluetun"
      xyz:
        image: xyz
        container_name: xyz
        network_mode: "service:gluetun"
    ```

    </p></details>

## Private Internet Access port forwarding

Note that [not all regions support port forwarding](https://www.privateinternetaccess.com/helpdesk/kb/articles/how-do-i-enable-port-forwarding-on-my-vpn).

When `PORT_FORWARDING=on`, a port will be forwarded on the VPN server side and written to the file specified by `PORT_FORWARDING_STATUS_FILE=/forwarded_port`.
It can be useful to mount this file as a volume to read it from other containers, for example to configure a torrenting client.

You can also use the HTTP control server (see below) to get the port forwarded.

## HTTP control server

See [its Wiki page](https://github.com/qdm12/gluetun/wiki/HTTP-control-server)

## Development and contributing

- Contribute with code: see [the Wiki](https://github.com/qdm12/gluetun/wiki/Contributing).
- [The list of existing contributors 👍](https://github.com/qdm12/gluetun/blob/master/.github/CONTRIBUTING.md#Contributors)
- [Github workflows](https://github.com/qdm12/gluetun/actions) to know what's building
- [List of issues and feature requests](https://github.com/qdm12/gluetun/issues)

## License

This repository is under an [MIT license](https://github.com/qdm12/gluetun/master/license)

## Support

Sponsor me on [Github](https://github.com/sponsors/qdm12), donate to [paypal.me/qmcgaw](https://www.paypal.me/qmcgaw) or subscribe to a VPN provider through one of my affiliate links:

[![https://github.com/sponsors/qdm12](https://raw.githubusercontent.com/qdm12/gluetun/master/doc/sponsors.jpg)](https://github.com/sponsors/qdm12)
[![https://www.paypal.me/qmcgaw](https://raw.githubusercontent.com/qdm12/gluetun/master/doc/paypal.jpg)](https://www.paypal.me/qmcgaw)

[![https://windscribe.com/?affid=mh7nyafu](https://raw.githubusercontent.com/qdm12/gluetun/master/doc/windscribe.jpg)](https://windscribe.com/?affid=mh7nyafu)

Feel also free to have a look at [the Kanban board](https://github.com/qdm12/gluetun/projects/1) and [contribute](#Development-and-contributing) to the code or the issues discussion.

Many thanks to @Frepke, @Ralph521, G. Mendez, M. Otmar Weber, J. Perez and A. Cooper for supporting me financially 🥇👍
