# Gluetun VPN client

*Lightweight swiss-knife-like VPN client to tunnel to Cyberghost, FastestVPN,
HideMyAss, IPVanish, IVPN, Mullvad, NordVPN, Privado, Private Internet Access, PrivateVPN,
ProtonVPN, PureVPN, Surfshark, TorGuard, VPNUnlimited, VyprVPN and Windscribe VPN servers
using Go, OpenVPN, iptables, DNS over TLS, ShadowSocks and an HTTP proxy*

**ANNOUNCEMENT**: You can try Wireguard, see #565

![Title image](https://raw.githubusercontent.com/qdm12/gluetun/master/title.svg)

[![Build status](https://github.com/qdm12/gluetun/actions/workflows/ci.yml/badge.svg)](https://github.com/qdm12/gluetun/actions/workflows/ci.yml)

[![Docker pulls qmcgaw/gluetun](https://img.shields.io/docker/pulls/qmcgaw/gluetun.svg)](https://hub.docker.com/r/qmcgaw/gluetun)
[![Docker pulls qmcgaw/private-internet-access](https://img.shields.io/docker/pulls/qmcgaw/private-internet-access.svg)](https://hub.docker.com/r/qmcgaw/gluetun)

[![Docker stars qmcgaw/gluetun](https://img.shields.io/docker/stars/qmcgaw/gluetun.svg)](https://hub.docker.com/r/qmcgaw/gluetun)
[![Docker stars qmcgaw/private-internet-access](https://img.shields.io/docker/stars/qmcgaw/private-internet-access.svg)](https://hub.docker.com/r/qmcgaw/gluetun)

![Last release](https://img.shields.io/github/release/qdm12/gluetun?label=Last%20release)
![Last Docker tag](https://img.shields.io/docker/v/qmcgaw/gluetun?sort=semver&label=Last%20Docker%20tag)
[![Last release size](https://img.shields.io/docker/image-size/qmcgaw/gluetun?sort=semver&label=Last%20released%20image)](https://hub.docker.com/r/qmcgaw/gluetun/tags?page=1&ordering=last_updated)
![GitHub last release date](https://img.shields.io/github/release-date/qdm12/gluetun?label=Last%20release%20date)
![Commits since release](https://img.shields.io/github/commits-since/qdm12/gluetun/latest?sort=semver)

[![Latest size](https://img.shields.io/docker/image-size/qmcgaw/gluetun/latest?label=Latest%20image)](https://hub.docker.com/r/qmcgaw/gluetun/tags)

[![GitHub last commit](https://img.shields.io/github/last-commit/qdm12/gluetun.svg)](https://github.com/qdm12/gluetun/commits/master)
[![GitHub commit activity](https://img.shields.io/github/commit-activity/y/qdm12/gluetun.svg)](https://github.com/qdm12/gluetun/graphs/contributors)
[![GitHub closed PRs](https://img.shields.io/github/issues-pr-closed/qdm12/gluetun.svg)](https://github.com/qdm12/gluetun/pulls?q=is%3Apr+is%3Aclosed)
[![GitHub issues](https://img.shields.io/github/issues/qdm12/gluetun.svg)](https://github.com/qdm12/gluetun/issues)
[![GitHub closed issues](https://img.shields.io/github/issues-closed/qdm12/gluetun.svg)](https://github.com/qdm12/gluetun/issues?q=is%3Aissue+is%3Aclosed)

[![Lines of code](https://img.shields.io/tokei/lines/github/qdm12/gluetun)](https://github.com/qdm12/gluetun)
![Code size](https://img.shields.io/github/languages/code-size/qdm12/gluetun)
![GitHub repo size](https://img.shields.io/github/repo-size/qdm12/gluetun)
![Go version](https://img.shields.io/github/go-mod/go-version/qdm12/gluetun)

![Visitors count](https://visitor-badge.laobi.icu/badge?page_id=gluetun.readme)

## Quick links

- Problem or suggestion?
  - [Start a discussion](https://github.com/qdm12/gluetun/discussions)
  - [Create an issue](https://github.com/qdm12/gluetun/issues)
  - [Check the Wiki](https://github.com/qdm12/gluetun/wiki)
  - [Join the Slack channel](https://join.slack.com/t/qdm12/shared_invite/enQtOTE0NjcxNTM1ODc5LTYyZmVlOTM3MGI4ZWU0YmJkMjUxNmQ4ODQ2OTAwYzMxMTlhY2Q1MWQyOWUyNjc2ODliNjFjMDUxNWNmNzk5MDk)
- Happy?
  - Sponsor me on [github.com/sponsors/qdm12](https://github.com/sponsors/qdm12)
  - Donate to [paypal.me/qmcgaw](https://www.paypal.me/qmcgaw)
  - Drop me [an email](mailto:quentin.mcgaw@gmail.com)
- Video:

  [![Video Gif](https://i.imgur.com/CetWunc.gif)](https://youtu.be/0F6I03LQcI4)

## Features

- Based on Alpine 3.14 for a small Docker image of 31MB
- Supports: **Cyberghost**, **FastestVPN**, **HideMyAss**, **IPVanish**, **IVPN**, **Mullvad**, **NordVPN**, **Privado**, **Private Internet Access**, **PrivateVPN**, **ProtonVPN**, **PureVPN**,  **Surfshark**, **TorGuard**, **VPNUnlimited**, **Vyprvpn**, **Windscribe** servers
- Supports OpenVPN and Wireguard (the latter in progress, see PR #565 and issue #134)
- DNS over TLS baked in with service provider(s) of your choice
- DNS fine blocking of malicious/ads/surveillance hostnames and IP addresses, with live update every 24 hours
- Choose the vpn network protocol, `udp` or `tcp`
- Built in firewall kill switch to allow traffic only with needed the VPN servers and LAN devices
- Built in Shadowsocks proxy (protocol based on SOCKS5 with an encryption layer, tunnels TCP+UDP)
- Built in HTTP proxy (tunnels HTTP and HTTPS through TCP)
- [Connect other containers to it](https://github.com/qdm12/gluetun/wiki/Connect-to-gluetun)
- [Connect LAN devices to it](https://github.com/qdm12/gluetun/wiki/Connect-to-gluetun)
- Compatible with amd64, i686 (32 bit), **ARM** 64 bit, ARM 32 bit v6 and v7, and even ppc64le 🎆
- VPN server side port forwarding for Private Internet Access and Vyprvpn
- Possibility of split horizon DNS by selecting multiple DNS over TLS providers
- Subprograms all drop root privileges once launched
- Subprograms output streams are all merged together
- Can work as a Kubernetes sidecar container, thanks @rorph

## Setup

1. Ensure your `tun` kernel module is setup:

    ```sh
    sudo modprobe tun
    # or, if you don't have modprobe, with
    sudo insmod /lib/modules/tun.ko
    ```

1. Extra steps:
    - [For Synology users](https://github.com/qdm12/gluetun/wiki/Synology-setup)
    - [For 32 bit Operating systems (**Rasberry Pis**)](https://github.com/qdm12/gluetun/wiki/32-bit-setup)
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
    docker-compose up -d
    ```

    You should probably check the many [environment variables](https://github.com/qdm12/gluetun/wiki/Environment-variables) available to adapt the container to your needs.

## Further setup

The following points are all optional but should give you insights on all the possibilities with this container.

- [Test your setup](https://github.com/qdm12/gluetun/wiki/Test-your-setup)
- [How to connect other containers and devices to Gluetun](https://github.com/qdm12/gluetun/wiki/Connect-to-gluetun)
- [VPN server side port forwarding](https://github.com/qdm12/gluetun/wiki/Port-forwarding)
- [HTTP control server](https://github.com/qdm12/gluetun/wiki/HTTP-Control-server) to automate things, restart Openvpn etc.
- Update the image with `docker pull qmcgaw/gluetun:latest`. See this [Wiki document](https://github.com/qdm12/gluetun/wiki/Docker-image-tags) for Docker tags available.
- Use [Docker secrets](https://github.com/qdm12/gluetun/wiki/Docker-secrets) to read your credentials instead of environment variables

## License

[![MIT](https://img.shields.io/github/license/qdm12/gluetun)](https://github.com/qdm12/gluetun/master/LICENSE)
