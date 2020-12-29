# Gluetun VPN client

*Lightweight swiss-knife-like VPN client to tunnel to Private Internet Access,
Mullvad, Windscribe, Surfshark Cyberghost, VyprVPN, NordVPN, PureVPN and Privado VPN servers, using Go, OpenVPN, iptables, DNS over TLS, ShadowSocks and an HTTP proxy*

**ANNOUNCEMENT**: *New Docker image name `qmcgaw/gluetun`*

<img height="250" src="https://raw.githubusercontent.com/qdm12/gluetun/master/title.svg?sanitize=true">

[![Build status](https://github.com/qdm12/gluetun/workflows/Buildx%20latest/badge.svg)](https://github.com/qdm12/gluetun/actions?query=workflow%3A%22Buildx+latest%22)

[![Docker Pulls](https://img.shields.io/docker/pulls/qmcgaw/private-internet-access.svg)](https://hub.docker.com/r/qmcgaw/private-internet-access)
[![Docker Stars](https://img.shields.io/docker/stars/qmcgaw/private-internet-access.svg)](https://hub.docker.com/r/qmcgaw/private-internet-access)

[![Docker Pulls](https://img.shields.io/docker/pulls/qmcgaw/gluetun.svg)](https://hub.docker.com/r/qmcgaw/gluetun)
[![Docker Stars](https://img.shields.io/docker/stars/qmcgaw/gluetun.svg)](https://hub.docker.com/r/qmcgaw/gluetun)

[![GitHub last commit](https://img.shields.io/github/last-commit/qdm12/gluetun.svg)](https://github.com/qdm12/gluetun/issues)
[![GitHub commit activity](https://img.shields.io/github/commit-activity/y/qdm12/gluetun.svg)](https://github.com/qdm12/gluetun/issues)
[![GitHub issues](https://img.shields.io/github/issues/qdm12/gluetun.svg)](https://github.com/qdm12/gluetun/issues)

[![Image size](https://images.microbadger.com/badges/image/qmcgaw/gluetun.svg)](https://microbadger.com/images/qmcgaw/gluetun)
[![Image version](https://images.microbadger.com/badges/version/qmcgaw/gluetun.svg)](https://microbadger.com/images/qmcgaw/gluetun)
[![Join Slack channel](https://img.shields.io/badge/slack-@qdm12-yellow.svg?logo=slack)](https://join.slack.com/t/qdm12/shared_invite/enQtOTE0NjcxNTM1ODc5LTYyZmVlOTM3MGI4ZWU0YmJkMjUxNmQ4ODQ2OTAwYzMxMTlhY2Q1MWQyOWUyNjc2ODliNjFjMDUxNWNmNzk5MDk)

## Videos

1. [**Introduction**](https://youtu.be/3jIbU6J2Hs0)
1. [**Connect a container**](https://youtu.be/mH7J_2JKNK0)
1. [**Connect LAN devices**](https://youtu.be/qvjrM15Y0uk)

## Features

- Based on Alpine 3.12 for a small Docker image of 52MB
- Supports **Private Internet Access**, **Mullvad**, **Windscribe**, **Surfshark**, **Cyberghost**, **Vyprvpn**, **NordVPN**, **PureVPN** and **Privado** servers
- Supports Openvpn only for now
- DNS over TLS baked in with service provider(s) of your choice
- DNS fine blocking of malicious/ads/surveillance hostnames and IP addresses, with live update every 24 hours
- Choose the vpn network protocol, `udp` or `tcp`
- Built in firewall kill switch to allow traffic only with needed the VPN servers and LAN devices
- Built in Shadowsocks proxy (protocol based on SOCKS5 with an encryption layer, tunnels TCP+UDP)
- Built in HTTP proxy (tunnels HTTP and HTTPS through TCP)
- [Connect other containers to it](https://github.com/qdm12/gluetun#connect-to-it)
- [Connect LAN devices to it](https://github.com/qdm12/gluetun#connect-to-it)
- Compatible with amd64, i686 (32 bit), **ARM** 64 bit, ARM 32 bit v6 and v7 🎆
- VPN server side port forwarding for Private Internet Access and Vyprvpn
- Possibility of split horizon DNS by selecting multiple DNS over TLS providers
- Subprograms all drop root privileges once launched
- Subprograms output streams are all merged together
- Can work as a Kubernetes sidecar container, thanks @rorph

## Setup

1. On some devices you may need to setup your tunnel kernel module on your host with `insmod /lib/modules/tun.ko` or `modprobe tun`
    - [Synology users Wiki page](https://github.com/qdm12/gluetun/wiki/Synology-setup)
1. Launch the container with:

    ```bash
    docker run -d --name gluetun --cap-add=NET_ADMIN \
    -e VPNSP="private internet access" -e REGION="CA Montreal" \
    -e OPENVPN_USER=js89ds7 -e OPENVPN_PASSWORD=8fd9s239G \
    -v /yourpath:/gluetun \
    qmcgaw/gluetun
    ```

    or use [docker-compose.yml](https://github.com/qdm12/gluetun/blob/master/docker-compose.yml) with:

    ```bash
    echo "your openvpn username" > openvpn_user
    echo "your openvpn password" > openvpn_password
    docker-compose up -d
    ```

    Note that you can:

    - Change the many [environment variables](#environment-variables) available
    - Use `-p 8888:8888/tcp` to access the HTTP web proxy
    - Use `-p 8388:8388/tcp -p 8388:8388/udp` to access the Shadowsocks proxy
    - Use `-p 8000:8000/tcp` to access the [HTTP control server](#HTTP-control-server) built-in
    - Use [Docker secrets](#Docker-secrets) to read your credentials instead of environment variables

    **If you encounter an issue with the tun device not being available, see [the FAQ](https://github.com/qdm12/gluetun/blob/master/doc/faq.md#how-to-fix-openvpn-failing-to-start)**

1. You can update the image with `docker pull qmcgaw/gluetun:latest`. See the [wiki](https://github.com/qdm12/gluetun/wiki/Common-issues#use-a-release-tag) for more information on other tags available.

## Testing

Check the VPN IP address matches your expectations

```sh
docker run --rm --network=container:gluetun alpine:3.12 wget -qO- https://ipinfo.io
```

▶ [Testing Wiki page](https://github.com/qdm12/gluetun/wiki/Testing-the-setup)

## Environment variables

**TLDR**; only set the 🏁 marked environment variables to get started.

💡 For all server filtering options such as `REGION`, you can have multiple values separated by a comma, i.e. `Germany,Singapore`

### VPN

| Variable | Default | Choices | Description |
| --- | --- | --- | --- |
| 🏁 `VPNSP` | `private internet access` | `private internet access`, `mullvad`, `windscribe`, `surfshark`, `vyprvpn`, `nordvpn`, `purevpn`, `privado` | VPN Service Provider |
| `PUBLICIP_FILE` | `/tmp/gluetun/ip` | Any filepath | Filepath to store the public IP address assigned |
| `PROTOCOL` | `udp` | `udp` or `tcp` | Network protocol to use |
| `OPENVPN_VERBOSITY` | `1` | `0` to `6` | Openvpn verbosity level |
| `OPENVPN_ROOT` | `no` | `yes` or `no` | Run OpenVPN as root |
| `OPENVPN_TARGET_IP` | | Valid IP address | Specify a target VPN IP address to use |
| `OPENVPN_CIPHER` | | i.e. `aes-256-gcm` | Specify a custom cipher to use. It will also set `ncp-disable` if using AES GCM for PIA |
| `OPENVPN_AUTH` | | i.e. `sha256` | Specify a custom auth algorithm to use |
| `OPENVPN_IPV6` | `off` | `on`, `off` | Enable tunneling of IPv6 (only for Mullvad) |

*For all providers below, server location parameters are all optional. By default a random server is picked using the filter settings provided.*

- Private Internet Access

    | Variable | Default | Choices | Description |
    | --- | --- | --- | --- |
    | 🏁 `OPENVPN_USER` | | | Your username |
    | 🏁 `OPENVPN_PASSWORD` | | | Your password |
    | `REGION` | | One of the [PIA regions](https://www.privateinternetaccess.com/pages/network/) | VPN server region |
    | `PIA_ENCRYPTION` | `strong` | `normal`, `strong` | Encryption preset |
    | `PORT_FORWARDING` | `off` | `on`, `off` | Enable port forwarding on the VPN server |
    | `PORT_FORWARDING_STATUS_FILE` | `/tmp/gluetun/forwarded_port` | Any filepath | Filepath to store the forwarded port number |

- Mullvad

    | Variable | Default | Choices | Description |
    | --- | --- | --- | --- |
    | 🏁 `OPENVPN_USER` | | | Your user ID |
    | `COUNTRY` | | One of the [Mullvad countries](https://mullvad.net/en/servers/#openvpn) | VPN server country |
    | `CITY` | | One of the [Mullvad cities](https://mullvad.net/en/servers/#openvpn) | VPN server city |
    | `ISP` | | One of the [Mullvad ISP](https://mullvad.net/en/servers/#openvpn) | VPN server ISP |
    | `PORT` | | `80`, `443` or `1401` for TCP; `53`, `1194`, `1195`, `1196`, `1197`, `1300`, `1301`, `1302`, `1303` or `1400` for UDP. Defaults to TCP `443` and UDP `1194` | Custom VPN port to use |
    | `OWNED` | `no` | `yes` or `no` | If the VPN server is owned by Mullvad |

    💡 [Mullvad IPv6 Wiki page](https://github.com/qdm12/gluetun/wiki/Mullvad-IPv6)

    For **port forwarding**, obtain a port from [here](https://mullvad.net/en/account/#/ports) and add it to `FIREWALL_VPN_INPUT_PORTS`

- Windscribe (see [this](https://github.com/qdm12/gluetun/blob/master/internal/constants/windscribe.go#L43) for the choices of regions, cities and hostnames)

    | Variable | Default | Choices | Description |
    | --- | --- | --- | --- |
    | 🏁 `OPENVPN_USER` | | | Your username |
    | 🏁 `OPENVPN_PASSWORD` | | | Your password |
    | `REGION` | | | Comma separated list of regions to choose the VPN server |
    | `CITY` | | | Comma separated list of cities to choose the VPN server |
    | `HOSTNAME` | | | Comma separated list of hostnames to choose the VPN server |
    | `PORT` | | One from the [this list of ports](https://windscribe.com/getconfig/openvpn) | Custom VPN port to use |

- Surfshark

    | Variable | Default | Choices | Description |
    | --- | --- | --- | --- |
    | 🏁 `OPENVPN_USER` | | | Your **service** username, found at the bottom of the [manual setup page](https://account.surfshark.com/setup/manual) |
    | 🏁 `OPENVPN_PASSWORD` | | | Your **service** password |
    | `REGION` | | One of the [Surfshark regions](https://github.com/qdm12/gluetun/wiki/Surfshark-Servers) | VPN server region |

- Cyberghost

    | Variable | Default | Choices | Description |
    | --- | --- | --- | --- |
    | 🏁 `OPENVPN_USER` | | | Your username |
    | 🏁 `OPENVPN_PASSWORD` | | | Your password |
    | 🏁 | | | **See additional setup steps below** |
    | `REGION` | | One of the Cyberghost regions, [Wiki page](https://github.com/qdm12/gluetun/wiki/Cyberghost-Servers) | VPN server country |
    | `CYBERGHOST_GROUP` | `Premium UDP Europe` | One of the server groups (see above Wiki page) | Server group |

    **Additional setup steps**: If you use docker Swarm or docker-compose, you should use the [Docker secrets](#Docker-secrets) `openvpn_clientkey` and `openvpn_clientcrt`.

    Otherwise, bind mount your `client.key` and `client.crt` files with, for example:

    ```sh
    -v /yourpath/client.key:/gluetun/client.key:ro -v /yourpath/client.crt:/gluetun/client.crt:ro
    ```

- Vyprvpn

    | Variable | Default | Choices | Description |
    | --- | --- | --- | --- |
    | 🏁 `OPENVPN_USER` | | | Your username |
    | 🏁 `OPENVPN_PASSWORD` | | | Your password |
    | `REGION` | | One of the [VyprVPN regions](https://www.vyprvpn.com/server-locations) | VPN server region |

    For **port forwarding**, add a port you want to be accessible to `FIREWALL_VPN_INPUT_PORTS`

- NordVPN

    | Variable | Default | Choices | Description |
    | --- | --- | --- | --- |
    | 🏁 `OPENVPN_USER` | | | Your username |
    | 🏁 `OPENVPN_PASSWORD` | | | Your password |
    | `REGION` | | One of the NordVPN server country, i.e. `Switzerland` | VPN server country |
    | `SERVER_NUMBER` | | Server integer number | Optional server number. For example `251` for `Italy #251` |

- PureVPN

    | Variable | Default | Choices | Description |
    | --- | --- | --- | --- |
    | 🏁 `OPENVPN_USER` | | | Your username |
    | 🏁 `OPENVPN_PASSWORD` | | | Your password |
    | `REGION` | | One of the [PureVPN regions](https://support.purevpn.com/vpn-servers) | VPN server region |
    | `COUNTRY` | | One of the [PureVPN countries](https://support.purevpn.com/vpn-servers) | VPN server country |
    | `CITY` | | One of the [PureVPN cities](https://support.purevpn.com/vpn-servers) | VPN server city |

- Privado

    | Variable | Default | Choices | Description |
    | --- | --- | --- | --- |
    | 🏁 `OPENVPN_USER` | | | Your username |
    | 🏁 `OPENVPN_PASSWORD` | | | Your password |
    | `HOSTNAME` | | [One of the Privado hostname](internal/constants/privado.go#L26), i.e. `ams-001.vpn.privado.io` | VPN server hostname |

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

### Firewall and routing

| Variable | Default | Choices | Description |
| --- | --- | --- | --- |
| `FIREWALL` | `on` | `on` or `off` | Turn on or off the container built-in firewall. You should use it for **debugging purposes** only. |
| `FIREWALL_VPN_INPUT_PORTS` | | i.e. `1000,8080` | Comma separated list of ports to allow from the VPN server side (useful for **vyprvpn** port forwarding) |
| `FIREWALL_INPUT_PORTS` | | i.e. `1000,8000` | Comma separated list of ports to allow through the default interface. This seems needed for Kubernetes sidecars. |
| `FIREWALL_DEBUG` | `off` | `on` or `off` | Prints every firewall related command. You should use it for **debugging purposes** only. |
| `FIREWALL_OUTBOUND_SUBNETS` | | i.e. `192.168.1.0/24,192.168.10.121,10.0.0.5/28` | Comma separated subnets that Gluetun and the containers sharing its network stack are allowed to access. This involves firewall and routing modifications. |

### Shadowsocks

| Variable | Default | Choices | Description |
| --- | --- | --- | --- |
| `SHADOWSOCKS` | `off` | `on`, `off` | Enable the internal Shadowsocks proxy |
| `SHADOWSOCKS_LOG` | `off` | `on`, `off` | Enable logging |
| `SHADOWSOCKS_PORT` | `8388` | `1024` to `65535` | Internal port number for Shadowsocks to listen on |
| `SHADOWSOCKS_PASSWORD` | |  | Password to use to connect to Shadowsocks |
| `SHADOWSOCKS_METHOD` | `chacha20-ietf-poly1305` | `chacha20-ietf-poly1305`, `aes-128-gcm`, `aes-256-gcm` | Method to use for Shadowsocks |

### HTTP proxy

| Variable | Default | Choices | Description |
| --- | --- | --- | --- |
| `HTTPPROXY` | `off` | `on`, `off` | Enable the internal HTTP proxy |
| `HTTPPROXY_LOG` | `off` | `on` or `off` | Logs every tunnel requests |
| `HTTPPROXY_PORT` | `8888` | `1024` to `65535` | Internal port number for the HTTP proxy to listen on |
| `HTTPPROXY_USER` | | | Username to use to connect to the HTTP proxy |
| `HTTPPROXY_PASSWORD` | | | Password to use to connect to the HTTP proxy |
| `HTTPPROXY_STEALTH` | `off` | `on` or `off` | Stealth mode means HTTP proxy headers are not added to your requests |

### System

| Variable | Default | Choices | Description |
| --- | --- | --- | --- |
| `TZ` | | i.e. `Europe/London` | Specify a timezone to use to have correct log times |
| `PUID` | `1000` | | User ID to run as non root and for ownership of files written |
| `PGID` | `1000` | | Group ID to run as non root and for ownership of files written |

### HTTP Control server

| Variable | Default | Choices | Description |
| --- | --- | --- | --- |
| `HTTP_CONTROL_SERVER_PORT` | `8000` | `1` to `65535` | Listening port for the HTTP control server |
| `HTTP_CONTROL_SERVER_LOG` | `on` | `on` or `off` | Enable logging of HTTP requests |

### Other

| Variable | Default | Choices | Description |
| --- | --- | --- | --- |
| `PUBLICIP_PERIOD` | `12h` | Valid duration | Period to check for public IP address. Set to `0` to disable. |
| `VERSION_INFORMATION` | `on` | `on`, `off` | Logs a message indicating if a newer version is available once the VPN is connected |
| `UPDATER_PERIOD` | `0` | Valid duration string such as `24h` | Period to update all VPN servers information in memory and to /gluetun/servers.json. Set to `0` to disable. This does a burst of DNS over TLS requests, which may be blocked if you set `BLOCK_MALICIOUS=on` for example. |

## Docker secrets

If you use Docker Compose or Docker Swarm, you can optionally use [Docker secret files](https://docs.docker.com/engine/swarm/secrets/) for all sensitive values such as your Openvpn credentials, instead of using environment variables.

The following secrets can be used:

- `openvpn_user`
- `openvpn_password`
- `openvpn_clientkey`
- `openvpn_clientcrt`
- `httpproxy_username`
- `httpproxy_password`
- `shadowsocks_password`

By default, `openvpn_user` and `openvpn_password` are set in [docker-compose.yml](docker-compose.yml)

Note that you can change the secret file path in the container by changing the environment variable in the form `<capitalizedSecretName>_SECRETFILE`.
For example, `OPENVPN_PASSWORD_SECRETFILE` defaults to `/run/secrets/openvpn_password` which you can change.

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
- <details><summary>Connect LAN devices through the built-in HTTP proxy (i.e. with Chrome, Kodi, etc.)</summary><p>

    ⚠️ You might want to use Shadowsocks instead which tunnels UDP as well as TCP and does not leak your credentials.
    The HTTP proxy will not encrypt your username and password every time you send a request to the HTTP proxy server.

    1. Setup an HTTP proxy client, such as [SwitchyOmega for Chrome](https://chrome.google.com/webstore/detail/proxy-switchyomega/padekgcemlokbadohgkifijomclgjgif?hl=en)
    1. Ensure the Gluetun container is launched with:
        - port `8888` published `-p 8888:8888/tcp`
    1. With your HTTP proxy client, connect to the Docker host (i.e. `192.168.1.10`) on port `8888`. You need to enter your credentials if you set them with `HTTPPROXY_USER` and `HTTPPROXY_PASSWORD`. Note that Chrome does not support authentication.
    1. If you set `HTTPPROXY_LOG` to `on`, more information will be logged in the Docker logs

    </p></details>
- <details><summary>Connect LAN devices through the built-in *Shadowsocks* proxy (per app, system wide, etc.)</summary><p>

    1. Setup a Shadowsocks proxy client, there is a list of [ShadowSocks clients for **all platforms**](https://shadowsocks.org/en/download/clients.html)
        - **note** some clients do not tunnel UDP so your DNS queries will be done locally and not through Gluetun and its built in DNS over TLS
        - Clients that support such UDP tunneling are, as far as I know:
            - iOS: Potatso Lite
            - OSX: ShadowsocksX
            - Android: Shadowsocks by Max Lv
    1. Ensure the Gluetun container is launched with:
        - port `8388` published `-p 8388:8388/tcp -p 8388:8388/udp`
    1. With your Shadowsocks proxy client
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
        image: qmcgaw/gluetun
        container_name: gluetun
        cap_add:
          - NET_ADMIN
        environment:
          - OPENVPN_USER=js89ds7
          - OPENVPN_PASSWORD=8fd9s239G
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

When `PORT_FORWARDING=on`, a port will be forwarded on the VPN server side and written to the file specified by `PORT_FORWARDING_STATUS_FILE=/tmp/gluetun/forwarded_port`.
It can be useful to mount this file as a volume to read it from other containers, for example to configure a torrenting client.

For `VPNSP=private internet access` (default), you will keep the same forwarded port for 60 days as long as you bind mount the `/gluetun` directory.

You can also use the HTTP control server (see below) to get the port forwarded.

## HTTP control server

[Wiki page](https://github.com/qdm12/gluetun/wiki/HTTP-Control-server)

## Development and contributing

- Contribute with code: start with [this Wiki page](https://github.com/qdm12/gluetun/wiki/Developement-setup)
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
