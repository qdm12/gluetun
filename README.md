# Gluetun VPN client

*Lightweight swiss-knife-like VPN client to tunnel to Private Internet Access, Mullvad and Windscribe VPN servers, using Go, OpenVPN, iptables, DNS over TLS, ShadowSocks and Tinyproxy*

**ANNOUNCEMENT**: *New wiki available [here](https://github.com/qdm12/private-internet-access-docker/wiki)*

<img height="250" src="https://raw.githubusercontent.com/qdm12/private-internet-access-docker/master/title.svg?sanitize=true">

[![Build status](https://github.com/qdm12/private-internet-access-docker/workflows/Buildx%20latest/badge.svg)](https://github.com/qdm12/private-internet-access-docker/actions?query=workflow%3A%22Buildx+latest%22)
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

- Based on Alpine 3.11 for a small Docker image below 50MB
- Supports **Private Internet Access**, **Mullvad** and **Windscribe** servers
- DNS over TLS baked in with service provider(s) of your choice
- DNS fine blocking of malicious/ads/surveillance hostnames and IP addresses, with live update every 24 hours
- Choose the vpn network protocol, `udp` or `tcp`
- Built in firewall kill switch to allow traffic only with needed PIA servers and LAN devices
- Built in SOCKS5 proxy (Shadowsocks, tunnels TCP+UDP)
- Built in HTTP proxy (Tinyproxy, tunnels TCP)
- [Connect other containers to it](https://github.com/qdm12/private-internet-access-docker#connect-to-it)
- [Connect LAN devices to it](https://github.com/qdm12/private-internet-access-docker#connect-to-it)
- Compatible with amd64, i686 (32 bit), **ARM** 64 bit, ARM 32 bit v6 and v7 🎆

### Private Internet Access

- Pick the [region](https://www.privateinternetaccess.com/pages/network/)
- Pick the level of encryption
- Enable port forwarding

### Mullvad

- Pick the [country, city and ISP](https://mullvad.net/en/servers/#openvpn)
- Pick the port to use (i.e. `53` (udp) or `80` (tcp))

### Windscribe

- Pick the [region](https://windscribe.com/status)

### Extra niche features

- Possibility of split horizon DNS by selecting multiple DNS over TLS providers
- Subprograms all drop root privileges once launched
- Subprograms output streams are all merged together
- Can work as a Kubernetes sidecar container, thanks @rorph

## Setup

1. Requirements
    - *Ideally*, Docker 1.13, in order to have Docker API 1.25 which supports `init` (and, if you use docker-compose, docker-compose version 1.22.0)
    - A VPN account with one of the service providers:
        - Private Internet Access: **username** and **password**
        - Mullvad: user ID ([sign up](https://mullvad.net/en/account/))
        - Windscribe: **username** and **password** | Signup up using my affiliate link below

            [![https://windscribe.com/?affid=mh7nyafu](https://raw.githubusercontent.com/qdm12/private-internet-access-docker/master/doc/windscribe.jpg)](https://windscribe.com/?affid=mh7nyafu)

    - If you have a host or router firewall, please refer [to the firewall documentation](https://github.com/qdm12/private-internet-access-docker/blob/master/doc/firewall.md)

1. On some devices such as Synology or Qnap machines, it's required to setup your tunnel device `/dev/net/tun` on your host:

    ```sh
    insmod /lib/modules/tun.ko
    # or
    modprobe tun
    ```

    You can verify it's here with `ls /dev/net`

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
    - Use `-p 8888:8888/tcp` to access the HTTP web proxy (and put your LAN in `EXTRA_SUBNETS` environment variable, in example `192.168.1.0/24`)
    - Use `-p 8388:8388/tcp -p 8388:8388/udp` to access the SOCKS5 proxy (and put your LAN in `EXTRA_SUBNETS` environment variable, in example `192.168.1.0/24`)
    - Use `-p 8000:8000/tcp` to access the [HTTP control server](#HTTP-control-server) built-in
    - Pass additional arguments to *openvpn* using Docker's command function (commands after the image name)

    **If you encounter an issue with the tun device not being available, see [the FAQ](https://github.com/qdm12/private-internet-access-docker/blob/master/doc/faq.md#how-to-fix-openvpn-failing-to-start)**

1. You can update the image with `docker pull qmcgaw/private-internet-access:latest`. See the [wiki](https://github.com/qdm12/private-internet-access-docker/wiki/Common-issues#use-a-release-tag) for more information on other tags available.

## Testing

Check the PIA IP address matches your expectations

```sh
docker run --rm --network=container:pia alpine:3.11 wget -qO- https://ipinfo.io
```

Want more testing? ▶ [see the Wiki](https://github.com/qdm12/private-internet-access-docker/wiki/Testing)

## Environment variables

**TLDR**; only set the 🏁 marked environment variables to get started.

### VPN

| Variable | Default | Choices | Description |
| --- | --- | --- | --- |
| 🏁 `VPNSP` | `private internet access` | `private internet access`, `mullvad`, `windscribe`, `surfshark` | VPN Service Provider |
| `IP_STATUS_FILE` | `/ip` | Any filepath | Filepath to store the public IP address assigned |
| `PROTOCOL` | `udp` | `udp` or `tcp` | Network protocol to use |
| `OPENVPN_VERBOSITY` | `1` | `0` to `6` | Openvpn verbosity level |
| `OPENVPN_ROOT` | `no` | `yes` or `no` | Run OpenVPN as root |
| `OPENVPN_TARGET_IP` | | Valid IP address | Specify a target VPN server IP address to use |
| `OPENVPN_CIPHER` | | i.e. `aes-256-gcm` | Specify a custom cipher to use. It will also set `ncp-disable` if using AES GCM for PIA |
| `OPENVPN_AUTH` | | i.e. `sha256` | Specify a custom auth algorithm to use |

- Private Internet Access

    | Variable | Default | Choices | Description |
    | --- | --- | --- | --- |
    | 🏁 `USER` | | | Your username |
    | 🏁 `PASSWORD` | | | Your password |
    | `REGION` | `Austria` | One of the [PIA regions](https://www.privateinternetaccess.com/pages/network/) | VPN server region |
    | `PIA_ENCRYPTION` | `strong` | `normal`, `strong` | Encryption preset |
    | `PORT_FORWARDING` | `off` | `on`, `off` | Enable port forwarding on the VPN server |
    | `PORT_FORWARDING_STATUS_FILE` | `/forwarded_port` | Any filepath | Filepath to store the forwarded port number |

- Mullvad

    | Variable | Default | Choices | Description |
    | --- | --- | --- | --- |
    | 🏁 `USER` | | | Your user ID |
    | `COUNTRY` | `Sweden` | One of the [Mullvad countries](https://mullvad.net/en/servers/#openvpn) | VPN server country |
    | `CITY` | | One of the [Mullvad cities](https://mullvad.net/en/servers/#openvpn) | VPN server city |
    | `ISP` | | One of the [Mullvad ISP](https://mullvad.net/en/servers/#openvpn) | VPN server ISP |
    | `PORT` | | `80` or `443` for TCP; or `53` for UDP. Leave blank for default Mullvad server port | Custom VPN port to use |

- Windscribe

    | Variable | Default | Choices | Description |
    | --- | --- | --- | --- |
    | 🏁 `USER` | | | Your username |
    | 🏁 `PASSWORD` | | | Your password |
    | `REGION` | `Austria` | One of the [Windscribe regions](https://windscribe.com/status) | VPN server region |
    | `PORT` | | One from the [this list of ports](https://windscribe.com/getconfig/openvpn) | Custom VPN port to use |

- Surfshark

    | Variable | Default | Choices | Description |
    | --- | --- | --- | --- |
    | 🏁 `USER` | | | Your username |
    | 🏁 `PASSWORD` | | | Your password |
    | `REGION` | `Austria` | One of the [Surfshark regions (subdomains)](https://github.com/qdm12/private-internet-access-docker/wiki/surfshark) | VPN server region |

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

### Firewall

That one is important if you want to connect to the container from your LAN for example, using Shadowsocks or Tinyproxy.

| Variable | Default | Choices | Description |
| --- | --- | --- | --- |
| `EXTRA_SUBNETS` | | i.e. `192.168.1.0/24,192.168.10.121,10.0.0.5/28` | Comma separated subnets allowed in the container firewall |

### Shadowsocks

| Variable | Default | Choices | Description |
| --- | --- | --- | --- |
| `SHADOWSOCKS` | `on` | `on`, `off` | Enable the internal SOCKS5 proxy Shadowsocks |
| `SHADOWSOCKS_LOG` | `off` | `on`, `off` | Enable logging |
| `SHADOWSOCKS_PORT` | `8388` | `1024` to `65535` | Internal port number for Shadowsocks to listen on |
| `SHADOWSOCKS_PASSWORD` | `8388` | `1024` to `65535` | Password to use to connect to Shadowsocks |
| `SHADOWSOCKS_METHOD` | `chacha20-ietf-poly1305` | One of [these ciphers](https://shadowsocks.org/en/config/quick-guide.html) | Method to use for Shadowsocks |

### Tinyproxy

| Variable | Default | Choices | Description |
| --- | --- | --- | --- |
| `TINYPROXY` | `on` | `on`, `off` | Enable the internal HTTP proxy tinyproxy |
| `TINYPROXY_LOG` | `Info` | `Info`, `Connect`, `Notice`, `Warning`, `Error`, `Critical` | Tinyproxy log level |
| `TINYPROXY_PORT` | `8888` | `1024` to `65535` | Internal port number for Tinyproxy to listen on |
| `TINYPROXY_USER` | `8388` | `1024` to `65535` | Username to use to connect to Tinyproxy |
| `TINYPROXY_PASSWORD` | `8388` | `1024` to `65535` | Password to use to connect to Tinyproxy |

### System

| Variable | Default | Choices | Description |
| --- | --- | --- | --- |
| `TZ` | | i.e. `Europe/London` | Specify a timezone to use to have correct log times |
| `UID` | `1000` | | User ID to run as non root and for ownership of files written |
| `GID` | `1000` | | Group ID to run as non root and for ownership of files written |

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
        - Choose the encryption method/algorithm to the method you specified in `SHADOWSOCKS_METHOD`
    1. If you set `SHADOWSOCKS_LOG` to `on`, (a lot) more information will be logged in the Docker logs

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

## Private Internet Access port forwarding

Note that [not all regions support port forwarding](https://www.privateinternetaccess.com/helpdesk/kb/articles/how-do-i-enable-port-forwarding-on-my-vpn).

When `PORT_FORWARDING=on`, a port will be forwarded on the PIA server side and written to the file specified by `PORT_FORWARDING_STATUS_FILE=/forwarded_port`.

It can be useful to mount this file as a volume to read it from other containers, for example to configure a torrenting client.

## HTTP control server

A built-in HTTP server listens on port `8000` to modify the state of the container. You have the following routes available:

- `http://<your-docker-host-ip>:8000/openvpn/actions/restart` restarts the openvpn process
- `http://<your-docker-host-ip>:8000/unbound/actions/restart` re-downloads the DNS files (crypto and block lists) and restarts the unbound process

## Development and contributing

- Contribute with code: see [the Wiki](https://github.com/qdm12/private-internet-access-docker/wiki/Contributing).
- [The list of existing contributors 👍](https://github.com/qdm12/private-internet-access-docker/blob/master/.github/CONTRIBUTING.md#Contributors)
- [Github workflows](https://github.com/qdm12/private-internet-access-docker/actions) to know what's building
- [List of issues and feature requests](https://github.com/qdm12/private-internet-access-docker/issues)

## License

This repository is under an [MIT license](https://github.com/qdm12/private-internet-access-docker/master/license)

## Support

Sponsor me on [Github](https://github.com/sponsors/qdm12), donate to [paypal.me/qmcgaw](https://www.paypal.me/qmcgaw) or subscribe to a VPN provider through one of my affiliate links:

[![https://github.com/sponsors/qdm12](https://raw.githubusercontent.com/qdm12/private-internet-access-docker/master/doc/sponsors.jpg)](https://github.com/sponsors/qdm12)
[![https://www.paypal.me/qmcgaw](https://raw.githubusercontent.com/qdm12/private-internet-access-docker/master/doc/paypal.jpg)](https://www.paypal.me/qmcgaw)

[![https://windscribe.com/?affid=mh7nyafu](https://raw.githubusercontent.com/qdm12/private-internet-access-docker/master/doc/windscribe.jpg)](https://windscribe.com/?affid=mh7nyafu)

Feel also free to have a look at [the Kanban board](https://github.com/qdm12/private-internet-access-docker/projects/1) and [contribute](#Development-and-contributing) to the code or the issues discussion.

Many thanks to @Frepke and @Ralph521 for supporting me financially 🥇👍
