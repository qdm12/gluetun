# Gluetun VPN client

*Lightweight swiss-knife-like VPN client to tunnel to Private Internet Access, Mullvad and Windscribe VPN servers, using Go, OpenVPN, iptables, DNS over TLS, ShadowSocks and Tinyproxy*

**ANNOUNCEMENT**: *Support for [Windscribe](https://windscribe.com/)*

<img height="200" src="title.svg?sanitize=true">

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
- DNS fine blocking of malicious/ads/surveillance hostnames and IP addresses
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
    - Docker 1.13, in order to have Docker API 1.25 which supports `init` (and, if you use docker-compose, docker-compose version 1.22.0)
    - A Private Internet Access **username** and **password** ([sign up](https://www.privateinternetaccess.com/pages/buy-vpn/)) or Mullvad user ID ([sign up](https://mullvad.net/en/account/))
    - <details><summary>External firewall requirements, if you have one</summary><p>

        - At start only
            - Allow outbound TCP 443 to github.com
            - If `DOT=on`, allow outbound TCP 853 to allow Unbound to resolve github.com and the PIA subdomain name if you use PIA.
            - If `DOT=off` and `VPNSP=pia`, allow outbound UDP 53 to your DNS provider to resolve the PIA subdomain name.
        - If `VPNSP=pia`, `ENCRYPTION=strong` and `PROTOCOL=udp`: allow outbound UDP 1197 to the corresponding VPN server IPs
        - If `VPNSP=pia`, `ENCRYPTION=normal` and `PROTOCOL=udp`: allow outbound UDP 1198 to the corresponding VPN server IPs
        - If `VPNSP=pia`, `ENCRYPTION=strong` and `PROTOCOL=tcp`: allow outbound TCP 501 to the corresponding VPN server IPs
        - If `VPNSP=pia`, `ENCRYPTION=normal` and `PROTOCOL=tcp`: allow outbound TCP 502 to the corresponding VPN server IPs
        - If `VPNSP=mullvad` and `PORT=`, please refer to the mapping of Mullvad servers in [these source code lines](https://github.com/qdm12/private-internet-access-docker/blob/master/internal/constants/mullvad.go#L64-L667) to find the corresponding UDP port number and IP address(es) of your choice
        - If `VPNSP=mullvad` and `PORT=53`, allow outbound UDP 53 to the corresponding VPN server IPs, which you can fine in [the mapping of Mullvad servers](https://github.com/qdm12/private-internet-access-docker/blob/master/internal/constants/mullvad.go#L64-L667)
        - If `VPNSP=mullvad` and `PORT=80`, allow outbound TCP 80 to the corresponding VPN server IPs, which you can fine in [the mapping of Mullvad servers](https://github.com/qdm12/private-internet-access-docker/blob/master/internal/constants/mullvad.go#L64-L667)
        - If `VPNSP=mullvad` and `PORT=443`, allow outbound TCP 443 to the corresponding VPN server IPs, which you can fine in [the mapping of Mullvad servers](https://github.com/qdm12/private-internet-access-docker/blob/master/internal/constants/mullvad.go#L64-L667)
        - If `VPNSP=windscribe` and `PROTOCOL=udp`: allow outbound UDP 443 to the corresponding VPN server IPs
        - If `VPNSP=windscribe` and `PROTOCOL=tcp`: allow outbound TCP 1194 to the corresponding VPN server IPs
        - If `SHADOWSOCKS=on`, allow inbound TCP 8388 and UDP 8388 from your LAN
        - If `TINYPROXY=on`, allow inbound TCP 8888 from your LAN

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

1. You can update the image with `docker pull qmcgaw/private-internet-access:latest`. There are also docker tags for older versions available:
    - `qmcgaw/private-internet-access:v2` linked to the [v2 release](https://github.com/qdm12/private-internet-access-docker/releases/tag/v2.0) (Golang based, only PIA)
    - `qmcgaw/private-internet-access:v1` linked to the [v1 release](https://github.com/qdm12/private-internet-access-docker/releases/tag/v1.0) (shell scripting based, no support, only PIA)
    - `qmcgaw/private-internet-access:old` tag, which is the latest shell scripting version (shell scripting based, no support, only PIA)

## Testing

Check the PIA IP address matches your expectations

```sh
docker run --rm --network=container:pia alpine:3.11 wget -qO- https://ipinfo.io
```

## Environment variables

| Environment variable | Default | Description |
| --- | --- | --- |
| `VPNSP` | `pia` | VPN Service Provider, one of `pia`, `mullvad` or `windscribe` |
| `REGION` | `Netherlands` | (PIA & Windscribe only) one of the [PIA regions](https://www.privateinternetaccess.com/pages/network/) or one of the [Windscribe regions](https://windscribe.com/status) |
| `COUNTRY` | `Sweden` | (Mullvad only) one of the [Mullvad countries](https://mullvad.net/en/servers/#openvpn) |
| `CITY` |  | (Mullvad only, *optional*) one of the [Mullvad cities](https://mullvad.net/en/servers/#openvpn) |
| `ISP` |  | (Mullvad only, *optional*) one of the [Mullvad ISP](https://mullvad.net/en/servers/#openvpn) |
| `PORT` |  | (Mullvad only, *optional*) For TCP, `80` or `443`, or `53` for UDP. Leave blank for default Mullvad server port |
| `PROTOCOL` | `udp` | `tcp` or `udp` |
| `ENCRYPTION` | `strong` | (PIA only) `normal` or `strong` |
| `USER` | | PIA username **or** Mullvad user ID **or** Windscribe username |
| `PASSWORD` | | Your PIA password **or** Windscribe password |
| `DOT` | `on` | `on` or `off`, to activate DNS over TLS to 1.1.1.1 |
| `DOT_PROVIDERS` | `cloudflare` | Comma delimited list of DNS over TLS providers from `cloudflare`, `google`, `quad9`, `quadrant`, `cleanbrowsing`, `securedns`, `libredns` |
| `DOT_CACHING` | `on` | Unbound caching feature, `on` or `off` |
| `DOT_IPV6` | `on` | Unbound will resolve domain names using IPv6 as well as IPv4 |
| `DOT_PRIVATE_ADDRESS` | All IPv4 and IPv6 CIDRs private ranges | Comma separated list of CIDRs or single IP addresses. Note that the default setting prevents DNS rebinding |
| `DOT_VERBOSITY` | `1` | Unbound verbosity level from `0` to `5` (full debug) |
| `DOT_VERBOSITY_DETAILS` | `0` | Unbound details verbosity level from `0` to `4` |
| `DOT_VALIDATION_LOGLEVEL` | `0` | Unbound validation log level from `0` to `2` |
| `BLOCK_MALICIOUS` | `on` | `on` or `off`, blocks malicious hostnames and IPs |
| `BLOCK_SURVEILLANCE` | `off` | `on` or `off`, blocks surveillance hostnames and IPs |
| `BLOCK_ADS` | `off` | `on` or `off`, blocks ads hostnames and IPs |
| `UNBLOCK` | | comma separated string (i.e. `web.com,web2.ca`) to unblock hostnames |
| `EXTRA_SUBNETS` | | comma separated subnets allowed in the container firewall (i.e. `192.168.1.0/24,192.168.10.121,10.0.0.5/28`) |
| `PORT_FORWARDING` | `off` | (PIA only) Set to `on` to forward a port on PIA server |
| `PORT_FORWARDING_STATUS_FILE` | `/forwarded_port` | (PIA only) File path to store the forwarded port number |
| `TINYPROXY` | `off` | `on` or `off`, to enable the internal HTTP proxy tinyproxy |
| `TINYPROXY_LOG` | `Info` | `Info`, `Connect`, `Notice`, `Warning`, `Error` or `Critical` |
| `TINYPROXY_PORT` | `8888` | `1024` to `65535` internal port for HTTP proxy |
| `TINYPROXY_USER` | | Username to use to connect to the HTTP proxy |
| `TINYPROXY_PASSWORD` | | Passsword to use to connect to the HTTP proxy |
| `SHADOWSOCKS` | `off` | `on` or `off`, to enable the internal SOCKS5 proxy Shadowsocks |
| `SHADOWSOCKS_LOG` | `off` | `on` or `off` to enable logging for Shadowsocks  |
| `SHADOWSOCKS_PORT` | `8388` | `1024` to `65535` internal port for SOCKS5 proxy |
| `SHADOWSOCKS_PASSWORD` | | Passsword to use to connect to the SOCKS5 proxy |
| `TZ` | | Specify a timezone to use i.e. `Europe/London` |
| `OPENVPN_VERBOSITY` | `1` | Openvpn verbosity level from 0 to 6 |
| `OPENVPN_ROOT` | `no` | Run OpenVPN as root, `yes` or `no` |
| `OPENVPN_TARGET_IP` | | (Optional) Specify a target VPN server IP address to use |

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

## FAQ

<details><summary>Openvpn disconnects because of a ping timeout</summary><p>

It happens especially on some PIA servers where they change their configuration or the server goes offline.

You will obtain an error similar to:

```s
openvpn: Wed Mar 18 22:13:00 2020 [3a51ae90324bcb0719cb399b650c64d4] Inactivity timeout (--ping-restart), restarting,
openvpn: Wed Mar 18 22:13:00 2020 SIGUSR1[soft,ping-restart] received, process restarting,
...
openvpn: Wed Mar 18 22:13:17 2020 Preserving previous TUN/TAP instance: tun0,
openvpn: Wed Mar 18 22:13:17 2020 NOTE: Pulled options changed on restart, will need to close and reopen TUN/TAP device.,
openvpn: Wed Mar 18 22:13:17 2020 ERROR: Linux route delete command failed: external program exited with error status: 2,
openvpn: Wed Mar 18 22:13:17 2020 ERROR: Linux route delete command failed: external program exited with error status: 2,
openvpn: Wed Mar 18 22:13:17 2020 ERROR: Linux route delete command failed: external program exited with error status: 2,
openvpn: Wed Mar 18 22:13:17 2020 ERROR: Linux route delete command failed: external program exited with error status: 2,
openvpn: Wed Mar 18 22:13:17 2020 /sbin/ip addr del dev tun0 local 10.6.11.6 peer 10.6.11.5,
openvpn: Wed Mar 18 22:13:17 2020 Linux ip addr del failed: external program exited with error status: 2,
openvpn: Wed Mar 18 22:13:18 2020 ERROR: Cannot ioctl TUNSETIFF tun: Operation not permitted (errno=1),
openvpn: Wed Mar 18 22:13:18 2020 Exiting due to fatal error,
exit status 1
```

To fix it, you would have to run openvpn with root, by setting the environment variable `OPENVPN_ROOT=yes`.

</p></details>

<details><summary>Private Internet Access: Why do I see openvpn warnings at start?</summary><p>

You might see some warnings similar to:

```s
openvpn: Sat Feb 22 15:55:02 2020 WARNING: this configuration may cache passwords in memory -- use the auth-nocache option to prevent this
openvpn: Sat Feb 22 15:55:02 2020 WARNING: 'link-mtu' is used inconsistently, local='link-mtu 1569', remote='link-mtu 1542'
openvpn: Sat Feb 22 15:55:02 2020 WARNING: 'cipher' is used inconsistently, local='cipher AES-256-CBC', remote='cipher BF-CBC'
openvpn: Sat Feb 22 15:55:02 2020 WARNING: 'auth' is used inconsistently, local='auth SHA256', remote='auth SHA1'
openvpn: Sat Feb 22 15:55:02 2020 WARNING: 'keysize' is used inconsistently, local='keysize 256', remote='keysize 128'
openvpn: Sat Feb 22 15:55:02 2020 WARNING: 'comp-lzo' is present in remote config but missing in local config, remote='comp-lzo'
openvpn: Sat Feb 22 15:55:02 2020 [a121ce520d670b71bfd3aa475485539b] Peer Connection Initiated with [AF_INET]xx.xx.xx.xx:1197
```

It is mainly because the option [disable-occ](https://openvpn.net/community-resources/reference-manual-for-openvpn-2-4/) was removed for transparency with you.

Private Internet Access explains [here why](https://www.privateinternetaccess.com/helpdesk/kb/articles/why-do-i-get-cipher-auth-warnings-when-i-connect) the warnings show up.

</p></details>

<details><summary>What files does it download at start before tunneling?</summary><p>

At start, the Go entrypoint only downloads, depending on your settings:

- If `DOT=on`: [DNS over TLS named root](https://github.com/qdm12/files/blob/master/named.root.updated) for Unbound
- If `DOT=on`: [DNS over TLS root key](https://github.com/qdm12/files/blob/master/root.key.updated) for Unbound
- If `BLOCK_MALICIOUS=on`: [Malicious hostnames and IP addresses block lists](https://github.com/qdm12/files) for Unbound
- If `BLOCK_SURVEILLANCE=on`: [Surveillance hostnames and IP addresses block lists](https://github.com/qdm12/files) for Unbound
- If `BLOCK_ADS=on`: [Ads hostnames and IP addresses block lists](https://github.com/qdm12/files) for Unbound

</p></details>

<details><summary>How to build Docker images of older or alternate versions</summary><p>

First, install [Git](https://git-scm.com/).

The following will build the Docker image locally and replace the previous one you built or pulled.

- Build the latest image

    ```sh
    docker build -t qmcgaw/private-internet-access https://github.com/qdm12/private-internet-access-docker.git
    ```

- Find a [commit](https://github.com/qdm12/private-internet-access-docker/commits/master) you want to build for, in example `095623925a9cc0e5cf89d5b9b510714792267d9b`, then:

    ```sh
    docker build -t qmcgaw/private-internet-access https://github.com/qdm12/private-internet-access-docker.git#095623925a9cc0e5cf89d5b9b510714792267d9b
    ```

- Find a [branch](https://github.com/qdm12/private-internet-access-docker/branches) you want to build for, in example `mullvad`, then:

    ```sh
    docker build -t qmcgaw/private-internet-access https://github.com/qdm12/private-internet-access-docker.git#mullvad
    ```

</p></details>

<details><summary>Mullvad does not work with IPv6?</summary><p>

By default, the Mullvad server tunnels both ipv4 and ipv6, hence openvpn will try to create an
ipv6 route. To allow the container to create such route, you have to specify `net.ipv6.conf.all.disable_ipv6=0`
at runtime, using either:

- For a Docker run command, the flag: `--sysctl net.ipv6.conf.all.disable_ipv6=0`
- In a docker-compose file:

    ```yml
        sysctls:
          - net.ipv6.conf.all.disable_ipv6=0
    ```

</p></details>

<details><summary>What's all this Go code?</summary><p>

The Go code is a big rewrite of the previous shell entrypoint, it allows for:

- better testing
- better maintainability
- ease of implementing new features
- faster boot
- asynchronous/parallel operations

It is mostly made of the [internal directory](https://github.com/qdm12/private-internet-access-docker/tree/master/internal) and the entry Go file [cmd/main.go](https://github.com/qdm12/private-internet-access-docker/blob/master/cmd/main.go).

</p></details>

<details><summary>How to test DNS over TLS?</summary><p>

- You can test DNSSEC using [internet.nl/connection](https://www.internet.nl/connection/)
- Check DNS leak tests with [https://www.dnsleaktest.com](https://www.dnsleaktest.com)
- Some other DNS leaks tests might not work because of [this](https://github.com/qdm12/cloudflare-dns-server#verify-dns-connection) (*TLDR*: Unbound DNS server is a local caching intermediary)

</p></details>

<details><summary>How to fix OpenVPN failing to start?</summary><p>

You can try:

- Installing the tun kernel module on your host with `insmod /lib/modules/tun.ko` or `modprobe tun`
- Adding `--device=/dev/net/tun` to your docker run command (equivalent for docker-compose, kubernetes, etc.)

</p></details>

## Development

1. Setup your environment

    <details><summary>Using VSCode and Docker</summary><p>

    1. Install [Docker](https://docs.docker.com/install/)
       - On Windows, share a drive with Docker Desktop and have the project on that partition
       - On OSX, share your project directory with Docker Desktop
    1. With [Visual Studio Code](https://code.visualstudio.com/download), install the [remote containers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)
    1. In Visual Studio Code, press on `F1` and select `Remote-Containers: Open Folder in Container...`
    1. Your dev environment is ready to go!... and it's running in a container :+1:

    </p></details>

    <details><summary>Locally</summary><p>

    Install [Go](https://golang.org/dl/), [Docker](https://www.docker.com/products/docker-desktop) and [Git](https://git-scm.com/downloads); then:

    ```sh
    go mod download
    go get github.com/golang/mock/gomock
    go get github.com/golang/mock/mockgen
    ```

    And finally install [golangci-lint](https://github.com/golangci/golangci-lint#install)

    </p></details>

1. Commands available:

    ```sh
    # Build the entrypoint binary
    go build cmd/main.go
    # Test the entrypoint code
    go test ./...
    # Lint the code
    golangci-lint run
    # Build the Docker image
    docker build -t qmcgaw/private-internet-access .
    ```

1. The Go code is in the Go file [cmd/main.go](https://github.com/qdm12/private-internet-access-docker/blob/master/cmd/main.go) and the [internal directory](https://github.com/qdm12/private-internet-access-docker/tree/master/internal), you might want to start reading the main.go file.
1. See [Contributing](.github/CONTRIBUTING.md) for more information on how to contribute to this repository.

### Contributors

Thanks for all the contributions, wether small or not so small!

- [@JeordyR](https://github.com/JeordyR) for testing the Mullvad version and opening a [PR with a few fixes](https://github.com/qdm12/private-internet-access-docker/pull/84/files) 👍
- [@rorph](https://github.com/rorph) for a [PR to pick a random region for PIA](https://github.com/qdm12/private-internet-access-docker/pull/70) and a [PR to make the container work with kubernetes](https://github.com/qdm12/private-internet-access-docker/pull/69)
- [@JesterEE](https://github.com/JesterEE) for a [PR to fix silly line endings in block lists back then](https://github.com/qdm12/private-internet-access-docker/pull/55) 📎
- [@elmerfdz](https://github.com/elmerfdz) for a [PR to add timezone information to have correct log timestampts](https://github.com/qdm12/private-internet-access-docker/pull/51) 🕙
- [@Juggels](https://github.com/Juggels) for a [PR to write the PIA forwarded port to a file](https://github.com/qdm12/private-internet-access-docker/pull/43)
- [@gdlx](https://github.com/gdlx) for a [PR to fix and improve PIA port forwarding script](https://github.com/qdm12/private-internet-access-docker/pull/32)
- [@janaz](https://github.com/janaz) for keeping an eye on [updating things in the Dockerfile](https://github.com/qdm12/private-internet-access-docker/pull/8)

## TODOs

<details><summary>Expand me</summary><p>

- Gotify support for notificactions
- Periodic update of malicious block lists with Unbound restart
- Improve healthcheck
    - Check IP address belongs to selected region
    - Check for DNS provider somehow if this is even possible
- Support for other VPN protocols
    - Wireguard (wireguard-go)
- Show new versions/commits available at start
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

</p></details>

## License

This repository is under an [MIT license](https://github.com/qdm12/private-internet-access-docker/master/license)
